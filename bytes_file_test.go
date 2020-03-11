package cfger

import (
	"os"
	"testing"
)

var factualFileAsString = `hello


wut
`

func TestUnstrBytesFile(t *testing.T) {
	os.Setenv("TESTFILEBYTESUNSTR", "file::./testdata/bytes")
	var val []byte
	_, err := ReadStructuredCfg("env::TESTFILEBYTESUNSTR", &val)
	if err != nil {
		t.Fatal(err)
	}

	if string(val) != factualFileAsString {
		t.Fatal("Read from file failed with inequality-error")
	}
}
