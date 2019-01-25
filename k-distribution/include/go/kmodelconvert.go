
package main

import (
	"fmt"
	koreparser "$INCLUDE_KORE_PARSER$"
	"strconv"
)

func convertParserModelToKModel(pk koreparser.K) K {
	switch v := pk.(type) {
	case koreparser.KApply:
		var convertedList []K
		for _, le := range v.List {
			convertedList = append(convertedList, convertParserModelToKModel(le))
		}
		return KApply{Label: parseKLabel(v.Label), List: convertedList}
	case koreparser.InjectedKLabel:
		return InjectedKLabel{Label: parseKLabel(v.Label)}
	case koreparser.KToken:
		return convertKToken(parseSort(v.Sort), v.Value)
	case koreparser.KVariable:
		return KVariable{Name: v.Name}
	case koreparser.KSequence:
		var convertedKs []K
		for _, ksElem := range v {
			convertedKs = append(convertedKs, convertParserModelToKModel(ksElem))
		}
		return KSequence{ks: convertedKs}
	default:
		panic(fmt.Sprintf("Unknown parser model K type: %#v", v))
	}
}

func convertKToken(sort Sort, value string) K {
	switch sort {
	case sortInt:
		i, err := strconv.Atoi(value)
		if err != nil {
			panic("Could not parse Int token: " + value)
		}
		return Int(i)
	case sortFloat:
		panic("Float token parse not implemented.")
	case sortString:
		unescapedStr := value // TODO: unescape value, see Ocaml impl unescape_k_string
		return String(unescapedStr)
	case sortBool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			panic("Could not parse Int token: " + value)
		}
		return Bool(b)
	default:
		return KToken{Value: value, Sort: sort}
	}
}
