package %PACKAGE_INTERPRETER%

type bytesHooksType int

const bytesHooks bytesHooksType = 0

func (bytesHooksType) empty(lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (bytesHooksType) bytes2int(c1 K, c2 K, c3 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (bytesHooksType) int2bytes(c1 K, c2 K, c3 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (bytesHooksType) bytes2string(c K,lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (bytesHooksType) string2bytes(c K,lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (bytesHooksType) substr(c1 K, c2 K, c3 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (bytesHooksType) replaceAt(c1 K, c2 K, c3 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (bytesHooksType) length(c K,lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (bytesHooksType) padRight(c1 K, c2 K, c3 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (bytesHooksType) padLeft(c1 K, c2 K, c3 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (bytesHooksType) reverse(c K,lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (bytesHooksType) concat(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

