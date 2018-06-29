package cfger

import (
	log "github.com/sirupsen/logrus"
	"testing"
	"os"
)

type jsonStruct struct {
	Version string
	Key1 struct {
		Valkey1 struct {
			Version int
		}
		Valkey2 struct {
			Valkeykey1 string `json:"valkeykey_1"`
			Valkeykey2 int    `json:"valkeykey_2"`
		}
	}
}

var factualStructured = jsonStruct{
	Version: "3.3",
	Key1: struct {
		Valkey1 struct {
			Version int
		}
		Valkey2 struct {
			Valkeykey1 string `json:"valkeykey_1"`
			Valkeykey2 int    `json:"valkeykey_2"`
		}
	}{
		Valkey1: struct{ Version int }{
			Version: 22,
		},
		Valkey2: struct {
			Valkeykey1 string `json:"valkeykey_1"`
			Valkeykey2 int    `json:"valkeykey_2"`
		}{
			Valkeykey1: "stringval",
			Valkeykey2: 3,
		},
	},
}

func setupJSON() {
	os.Setenv("TESTFILEJSON", "file::./testdata/test.json")
}

func TestJSON(t *testing.T) {
	setupJSON()

	a := jsonStruct{}
	_, err := ReadStructuredCfg("env::TESTFILEJSON", &a)
	if err != nil {
		log.Error(err)
	}

	if a != factualStructured {
		log.Fatal("Read from env::file failed with inequality-error")
	}

	log.Info("env::file:: json file to Go struct passed")

	a = jsonStruct{}
	_, err = ReadStructuredCfg("file::./testdata/test.json", &a)

	if err != nil {
		log.Error(err)
	}
	log.Info("file:: json file to Go struct passed")

	if a != factualStructured {
		log.Fatal("Read from file failed with inequality-error")
	}

}
