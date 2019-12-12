package cfger

import (
	log "github.com/sirupsen/logrus"
	"os"
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

type recurseYamlStruct struct {
	Root struct {
		Text string `yaml:"text"`
	} `yaml:"root"`
}

type recurseMixedStruct struct {
	Root struct {
		Data recurseJsonStruct `yaml:"data"`
	} `yaml:"root"`
}

func setupYAML() {
	os.Setenv("TESTFILE", "file::./testdata/test.yml")
	os.Setenv("RECURSE", "file::testdata/recurse.yml")
	os.Setenv("MIXED", "file::testdata/mixed.yml")
}

func TestYAML(t *testing.T) {
	setupYAML()

	a := yamlStruct{}
	err := ReadStructuredCfg("env::TESTFILE", &a)
	if err != nil {
		log.Error(err)
	}

	if a != factualYAMLStructured {
		log.Fatal("Read from environment failed with inequality-error")
	}

	log.Info("env::file:: yaml file to Go struct passed")

	a = yamlStruct{}
	err = ReadStructuredCfg("file::./testdata/test.yml", &a)

	if err != nil {
		log.Error(err)
	}
	log.Info("file:: yaml file to Go struct passed")

	if a != factualYAMLStructured {
		log.Fatal("Read from file failed with inequality-error")
	}

	var b = recurseYamlStruct{}

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
