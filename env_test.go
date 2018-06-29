package cfger

import (
	log "github.com/sirupsen/logrus"
	"testing"
	"os"
)

var (
	factualEnvValue = "this is a test value"
)

func setupEnv() {
	os.Setenv("TESTENV", factualEnvValue)
}

func TestEnv(t *testing.T) {
	setupEnv()

	val, err := ReadCfg("env::TESTENV")

	if err != nil {
		log.Fatal(err)
	} else if val != factualEnvValue {
		log.Fatal("Env read failed due to inequality checks")
	}

	log.Info("Env read passed")
}