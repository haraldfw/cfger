package cfger

import (
	log "github.com/sirupsen/logrus"
	"testing"
	"os"
)

var factualFileAsString = `hello


wut
`

func TestUnstrBytesFile(t *testing.T) {
	os.Setenv("TESTFILEBYTESUNSTR", "file::./testdata/bytes")
	var val []byte
	err := ReadStructuredCfg("env::TESTFILEBYTESUNSTR", &val)
	if err != nil {
		log.Error(err)
	}

	if string(val) != factualFileAsString {
		log.Fatal("Read from file failed with inequality-error")
	}

	log.Info("Unstructured file test passed")
}
