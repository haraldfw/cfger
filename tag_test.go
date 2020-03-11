package cfger

import (
	"os"
	"reflect"
	"testing"
)

type ValueBoi struct {
	Value1 string `json:"val1"`
	Value2 string `json:"val2"`
	Struct VBS    `json:"struct1"`
}

type VBS struct {
	Value        string `json:"structVal1"`
	NestedStruct VBNS   `json:"struct1.1"`
}

type VBNS struct {
	Value string `json:"structVal1.1"`
}

func TestTags(t *testing.T) {
	os.Setenv("VAL", "file::./tag-test.json")
	os.Setenv("ELLO_VALUE", "eyy")

	var ello ValueBoi

	_, err := ReadStructuredCfgRecursive("env::VAL", &ello)
	if err != nil {
		t.Fatal(err)
	}

	factualCfg := ValueBoi{
		Value1: "hello",
		Value2: "hoihoi",
		Struct: VBS{
			Value: "hoihoi",
			NestedStruct: VBNS{
				Value: "eyy",
			},
		},
	}

	if !reflect.DeepEqual(ello, factualCfg) {
		t.Fatal("mismatch between read file and factual file")
	}
}
