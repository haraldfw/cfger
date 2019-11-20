// This package provides functions to read configuration files in raw, JSON or YAML format from a file or environment
// variable. And additionally a shortcut to read docker-secrets. It also follows configuration locations in env-vars if
// the value has a valid prefix.
//
// Valid prefixes are:
//  1. env:: reads from environment variables.
//  2. secret:: reads file from /run/secrets
//  3. file:: reads file with given path
//
//
// Raw Versus Structured Read
//
// ReadCfg(val) wraps ReadStructuredCfg(val, nil), which results in a raw read. So if val here resolves to a JSON-file
// it returns the contents of the JSON-file as a string, without unmarshalling
//
// For cfger to unmarshal the file, you need to supply a valid interface in the form defined in the documentation for
// the package https://golang.org/pkg/encoding/json/#Unmarshal. Similarly for YAML-unmarshalling you have to supply a valid interface as described in the
// https://godoc.org/gopkg.in/yaml.v2#Unmarshal documentation.
package cfger

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	secretPrefLen = len("secret::")
	secretRoot    = "/run/secrets"
	filePrefLen   = len("file::")
	envPrefLen    = len("env::")
)

// Reads an unstructured configuration-value. See ReadStructuredCfg.
func ReadCfg(val string) (string, error) {
	var output string
	return output, ReadStructuredCfg(val, &output)
}

// Check if the input has a prefix, e. g. "secret::" or "file::", if it does try to read the file in the given location.
// In the case of neither the input will be returned. If a prefix is given but an error occurs when reading, the error
// is returned,
//
// Valid prefixes are:
//  1. env:: reads from environment variables.
//  2. secret:: reads file from /run/secrets
//  3. file:: reads file with given path
//
// If not prefix is provided, the input structure type is checked. If a string pointer, the value of the pointer
// is set to the input.
//
// Additionally if a file has .yml or .json as a suffix and the given interface{} is not nil, the interface will be used
// to unmarshal the file at the given path.
//
// If a read environment-variable contains a prefix this function will be called with the environment-variable's value.
func ReadStructuredCfg(val string, structure interface{}, populateRecursively ...bool) (error) {
	if structure == nil {
		return errors.New("no output structure provided")
	} else if len(populateRecursively) > 1 {
		return errors.New("more than one populateRecursively value passed")
	}

	switch strings.SplitN(val, "::", 2)[0] {
	case "secret":
		return ReadCfgFile(path.Join(secretRoot, val[secretPrefLen:]), structure, populateRecursively...)
	case "file":
		return ReadCfgFile(val[filePrefLen:], structure, populateRecursively...)
	case "env":
		return ReadEnv(val[envPrefLen:], structure, populateRecursively...)
	default:
		switch v := structure.(type) {
		case *string:
			var t *string
			t = structure.(*string)
			*t = val
			return nil
		default:
			return fmt.Errorf("unsupported structure type '%s'", v)
		}
	}
}

// Reads an environment-variable by the given name. If the environment-variable contains a prefix the path will be
// resolved.
func ReadEnv(val string, structure interface{}, populateRecursively ...bool) (error) {
	envVal, ok := os.LookupEnv(val)

	if !ok {
		return errors.New(fmt.Sprintf("Environment variable %q not found", val))
	}

	return ReadStructuredCfg(envVal, structure, populateRecursively...)
}

// Reads the file at the given path and returns the contents as a string. If the suffix of the path is .yml/.yaml/.json
// the contents are unmarshalled before they are returned. Returns error if an error is returned
// while reading or unmarshalling.
func ReadCfgFile(inPath string, structure interface{}, populateRecursively ...bool) (err error) {
	if structure == nil {
		return errors.New("no output structure provided")
	}

	var content []byte
	if content, err = ioutil.ReadFile(inPath); err != nil {
		return
	}

	switch structure.(type) {
	case *string:
		var t *string
		t = structure.(*string)
		*t = string(content)
	case *[]byte:
		var b *[]byte
		b = structure.(*[]byte)
		*b = content
	default:
		if strings.HasSuffix(inPath, ".yml") || strings.HasSuffix(inPath, ".yaml") {
			err = yaml.Unmarshal(content, structure)
		} else if strings.HasSuffix(inPath, ".json") {
			err = json.Unmarshal(content, structure)
		} else {
			err = errors.New("unsupported file type - supported types are 'yml' and 'json' ")
		}

		if err == nil && doRecurse(populateRecursively...) {
			return findVal(structure, -1, []int{})
		}
	}
	return
}

func findVal(structure interface{}, numFields int, fieldIndices []int) error {
	var elem reflect.Value
	{
		vos := reflect.ValueOf(structure)
		if vos.Kind() == reflect.Ptr {
			elem = vos.Elem()
		} else {
			elem = vos
		}
	}

	if elem.Kind() == reflect.Struct {

		if numFields == -1 {
			numFields = elem.NumField()
		}

		for i := 0; i < numFields; i++ {
			_fieldIndices := append(fieldIndices, i)

			field := elem.FieldByIndex(_fieldIndices)
			if !field.IsValid() || !field.CanSet() {
				continue
			}

			if field.Kind() == reflect.Struct {
				err := findVal(structure, field.NumField(), _fieldIndices)
				if err != nil {
					return err
				}
			} else if field.Kind() == reflect.String {
				//TODO: Generate test to check that this works.
				var output string
				err := ReadStructuredCfg(field.String(), &output)
				if err != nil {
					return err
				}
				field.SetString(output)
			}
		}
	}
	return nil
}

func doRecurse(populateRecursively ...bool) bool {
	return len(populateRecursively) > 0 && populateRecursively[0]
}
