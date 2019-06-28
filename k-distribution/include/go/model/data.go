%COMMENT%

package %PACKAGE_MODEL%

// ModelState holds the state of the executor at a certain moment
type ModelState struct {
	initialized bool

	// allKs keeps all KSequences into one large structure
	// all KSequences point to this structure
	allKs *ksequenceSliceContainer

	// keeps object types mixed together (for now)
	allObjects []KObject

	// memoTables is a structure containing all memoization maps.
	// Memoization tables are implemented as maps of maps of maps of ...
	memoTables map[MemoTable]interface{}
}

// constantsModel is another instance of the model, but which only contains a few constants.
var constantsModel = NewModel()

func (ms *ModelState) getObject(ref KReference) KObject {
	index := ref.value1
	if ref.constantObject {
		return constantsModel.allObjects[index]
	}
	if index >= len(ms.allObjects) {
		panic("trying to reference object beyond allocated objects")
	}
	return ms.allObjects[index]
}

func (ms *ModelState) addObject(obj KObject) KReference {
	newIndex := len(ms.allObjects)
	ms.allObjects = append(ms.allObjects, obj)
	return KReference{refType: obj.referenceType(), constantObject: false, value1: newIndex, value2: 0}
}

func addConstantObject(obj KObject) KReference {
	newIndex := len(constantsModel.allObjects)
	constantsModel.allObjects = append(constantsModel.allObjects, obj)
	return KReference{refType: obj.referenceType(), constantObject: true, value1: newIndex, value2: 0}
}

// NewModel creates a new blank model.
func NewModel() *ModelState {
	ms := &ModelState{}
	ms.Init()
	return ms
}

// Init prepares model for execution
func (ms *ModelState) Init() {
	if ms.initialized {
		return
	}
	ms.initialized = true
	ms.allKs = &ksequenceSliceContainer{}
	ms.allObjects = nil
	ms.memoTables = nil
}

// ClearModel ... clean up any data left from previous executions, to save memory
func (ms *ModelState) ClearModel() {
	ms.initialized = false
	ms.Init()
}

// MemoTable is a reference to a memoization table
type MemoTable int

// GetMemoizedValue searches for a value in the memo tables structure of the model.
func (ms *ModelState) GetMemoizedValue(memoTable MemoTable, keys ...KMapKey) (KReference, bool) {
	if ms.memoTables == nil {
		return NullReference, false
	}
	currentObj, tableFound := ms.memoTables[memoTable]
	if !tableFound {
		return NullReference, false
	}
	for _, key := range keys {
		currentMap, isMap := currentObj.(map[KMapKey]interface{})
		if !isMap {
			panic("wrong object found: memo tables need a level of map[KMapKey]interface{} for each key")
		}
		objectForKey, isKeyPresent := currentMap[key]
		if !isKeyPresent {
			return NullReference, false
		}
		currentObj = objectForKey
	}
	kref, isKref := currentObj.(KReference)
	if !isKref {
		panic("wrong object found: memo tables need to have a KReference on the last level")
	}
	return kref, true
}

// SetMemoizedValue inserts a value into the memo table structure, for a variable number of keys.
// It extends the tree up to where it is required.
func (ms *ModelState) SetMemoizedValue(memoized KReference, memoTable MemoTable, keys ...KMapKey) {
	if ms.memoTables == nil {
		ms.memoTables = make(map[MemoTable]interface{})
	}
	if len(keys) == 0 {
		// no keys, memo table is not really a table, it just contains one value
		ms.memoTables[memoTable] = memoized
		return
	}

	currentMapObj, tableFound := ms.memoTables[memoTable]
	if !tableFound {
		currentMapObj = make(map[KMapKey]interface{})
		ms.memoTables[memoTable] = currentMapObj
	}
	for i, key := range keys {
		currentMap, isMap := currentMapObj.(map[KMapKey]interface{})
		if !isMap {
			panic("wrong object found: memo tables need a level of map[KMapKey]interface{} for each key")
		}
		if i < len(keys)-1 {
			nextMap, nextMapExists := currentMap[key]
			if !nextMapExists {
				nextMap = make(map[KMapKey]interface{})
				currentMap[key] = nextMap
				currentMapObj = nextMap
			}
		} else {
			// last key
			currentMap[key] = memoized
		}
	}
}
