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

type recurseJsonStruct struct {
	Root struct {
		Text string `json:"text"`
	} `json:"root"`
}

func setupJSON() {
	os.Setenv("TESTFILEJSON", "file::./testdata/test.json")
	os.Setenv("RECURSE", "file::./testdata/recurse.json")
}

func TestJSON(t *testing.T) {
	setupJSON()

	a := jsonStruct{}
	err := ReadStructuredCfg("env::TESTFILEJSON", &a)
	if err != nil {
		log.Error(err)
	}

	if a != factualStructured {
		log.Fatal("Read from env::file failed with inequality-error")
	}

	log.Info("env::file:: json file to Go struct passed")

	a = jsonStruct{}
	err = ReadStructuredCfg("file::./testdata/test.json", &a)

	if err != nil {
		log.Error(err)
	}
	log.Info("file:: json file to Go struct passed")

	if a != factualStructured {
		log.Fatal("Read from file failed with inequality-error")
	}


	var b = recurseJsonStruct{}

	err = ReadStructuredCfg("env::RECURSE", &b, true)
	if b.Root.Text != "data" {
		t.Fatalf("Text was expected to be 'data', but was '%s'", b.Root.Text)
	} else if err != nil {
		t.Fatalf("Got err when none was expected: %s", err.Error())
	}

	err = ReadStructuredCfg("env::RECURSE", 0, true)
	if err == nil {
		t.Fatal("Wanted an err, but got none")
	}
}