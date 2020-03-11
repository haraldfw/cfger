package cfger

import (
	"os"
	"testing"
)

var factualFile = `{
  "asdf": "lkjh",
  "a": [
    1,
    2,
    3
  ]
}`

func TestUnstrFile(t *testing.T) {
	os.Setenv("TESTFILEUNSTR", "file::./testdata/unstructured.json")
	val, err := ReadCfg("env::TESTFILEUNSTR")
	if err != nil {
		t.Fatal(err)
	}

	if val != factualFile {
		t.Fatal("Read from file failed with inequality-error")
	}
}
