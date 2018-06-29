package cfger

import (
	log "github.com/sirupsen/logrus"
	"testing"
	"os"
)

var factualYAMLStructured = yamlStruct{
	Version: "3.3",
	Key1: struct {
		Valkey1 struct {
			Version int
		}
		Valkey2 struct {
			Valkeykey1 string `yaml:"valkeykey_1"`
			Valkeykey2 int    `yaml:"valkeykey_2"`
		}
	}{
		Valkey1: struct{ Version int }{
			Version: 22,
		},
		Valkey2: struct {
			Valkeykey1 string `yaml:"valkeykey_1"`
			Valkeykey2 int    `yaml:"valkeykey_2"`
		}{
			Valkeykey1: "stringval",
			Valkeykey2: 3,
		},
	},
}

type yamlStruct struct {
	Version string
	Key1 struct {
		Valkey1 struct {
			Version int
		}
		Valkey2 struct {
			Valkeykey1 string `yaml:"valkeykey_1"`
			Valkeykey2 int    `yaml:"valkeykey_2"`
		}
	}
}


func setupYAML() {
	os.Setenv("TESTFILE", "file::./testdata/test.yml")
}

func TestYAML(t *testing.T) {
	setupYAML()

	a := yamlStruct{}
	_, err := ReadStructuredCfg("env::TESTFILE", &a)
	if err != nil {
		log.Error(err)
	}

	if a != factualYAMLStructured {
		log.Fatal("Read from environment failed with inequality-error")
	}

	log.Info("env::file:: yaml file to Go struct passed")

	a = yamlStruct{}
	_, err = ReadStructuredCfg("file::./testdata/test.yml", &a)

	if err != nil {
		log.Error(err)
	}
	log.Info("file:: yaml file to Go struct passed")


	if a != factualYAMLStructured {
		log.Fatal("Read from file failed with inequality-error")
	}
}
