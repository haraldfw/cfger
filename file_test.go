package cfger

import (
	log "github.com/sirupsen/logrus"
	"testing"
	"os"
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
		log.Error(err)
	}

	if val != factualFile {
		log.Fatal("Read from file failed with inequality-error")
	}

	log.Info("Unstructured file test passed")
}
