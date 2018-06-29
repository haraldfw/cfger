// This package provides functions to read configuration files in raw, JSON or YAML format from a file or environment
// variable. And additionally a shortcut to read docker-secrets. It also follows configuration locations in env-vars if
// the value has a valid prefix.
//
// Valid prefixes are:
//  1. env:: reads from environment variables.
//  2. secrets:: reads file from /run/secrets
//  3. file:: reads file with given path
//
//
// Raw Versus Structured Read
//
// ReadCfg(val) wraps ReadStructuredCfg(val, nil), which results in a raw read. So if val here resolves to a JSON-file
// it returns the contents of the JSON-file as a string, without unmarshalling.
//
// For cfger to unmarshal the file, you need to supply a valid interface in the form defined in the documentation for
// the package https://golang.org/pkg/encoding/json/#Unmarshal. Similarly for YAML-unmarshalling you have to supply a valid interface as described in the
// https://godoc.org/gopkg.in/yaml.v2#Unmarshal documentation.
package cfger

import (
	"strings"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"encoding/json"
	"path"
	"os"
	"errors"
	"fmt"
)

var (
	secretPrefLen = len("secret::")
	secretRoot    = "/run/secrets"
	filePrefLen   = len("file::")
	envPrefLen    = len("env::")
)

// Reads an unstructured configuration-value. See ReadStructuredCfg.
func ReadCfg(val string) (string, error) {
	return ReadStructuredCfg(val, nil)
}

// Check if the input has a prefix, e. g. "secret::" or "file::", if it does try to read the file in the given location.
// In the case of neither the input will be returned. If a prefix is given but an error occurs when reading, an empty
// string and the error will be returned. If no prefix is given the given val i returned.
//
// Valid prefixes are:
//  1. env:: reads from environment variables.
//  2. secrets:: reads file from /run/secrets
//  3. file:: reads file with given path
//
// Additionally if a file has .yml or .json as a suffix and the given interface{} is not nil, the interface will be used
// to unmarshal the file at the given path.
//
// If a read environment-variable contains a prefix this function will be called with the environment-variable's value.
func ReadStructuredCfg(val string, structure interface{}) (string, error) {
	switch strings.SplitN(val, "::", 2)[0] {
	case "secret":
		return ReadCfgFile(path.Join(secretRoot, val[secretPrefLen:]), structure)
	case "file":
		return ReadCfgFile(val[filePrefLen:], structure)
	case "env":
		return ReadEnv(val[envPrefLen:], structure)
	default:
		return val, nil
	}
}

// Reads an environment-variable by the given name. If the environment-variable contains a prefix the path will be
// resolved.
func ReadEnv(val string, structure interface{}) (string, error) {
	envVal, ok := os.LookupEnv(val)

	if !ok {
		return "", errors.New(fmt.Sprintf("Environment variable %q not found", val))
	}

	return ReadStructuredCfg(envVal, structure)
}

// Reads the file at the given path and returns the contents as a string. If the suffix of the path is .yml/.yaml/.json
// the contents are unmarshalled before they are returned. Returns an empty string and an error if an error is returned
// while reading or unmarshalling.
func ReadCfgFile(inPath string, structure interface{}) (string, error) {
	content, err := ioutil.ReadFile(inPath)
	if err != nil {
		return "", err
	}

	if structure != nil{
		if strings.HasSuffix(inPath, ".yml") || strings.HasSuffix(inPath, ".yaml") {
			err = yaml.Unmarshal(content, structure)
			return "", err
		} else if strings.HasSuffix(inPath, ".json") {
			err = json.Unmarshal(content, structure)
			return "", err
		}
	}
	return string(content), nil
}
