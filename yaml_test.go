package cfger

import (
	"os"
	"reflect"
	"testing"
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
	Key1    struct {
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
		t.Fatal(err)
	}

	if !reflect.DeepEqual(a, factualYAMLStructured) {
		t.Fatal("Read from environment failed with inequality-error")
	}

	a = yamlStruct{}
	_, err = ReadStructuredCfg("file::./testdata/test.yml", &a)

	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, factualYAMLStructured) {
		t.Fatal("Read from file failed with inequality-error")
	}
}
