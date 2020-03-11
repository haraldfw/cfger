package cfger

import (
	"os"
	"testing"
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
		t.Fatal(err)
	} else if val != factualEnvValue {
		t.Fatal("Env read failed due to inequality checks")
	}
}
