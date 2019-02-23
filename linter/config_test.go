package linter

import (
	"reflect"
	"strings"
	"testing"

	"github.com/skeema/mybase"
	"github.com/skeema/skeema/fs"
)

func TestOptionsForDir(t *testing.T) {
	dir := getDir(t, "../testdata/linter/valid")
	if opts, err := OptionsForDir(dir); err != nil {
		t.Errorf("Unexpected error from OptionsForDir: %s", err)
	} else {
		expected := Options{
			ProblemSeverity: map[string]Severity{
				"no-pk":       SeverityError,
				"bad-charset": SeverityWarning,
				"bad-engine":  SeverityWarning,
			},
			AllowedCharSets: []string{"utf8mb4", "utf8"},
			AllowedEngines:  []string{"myisam"},
		}
		if !reflect.DeepEqual(opts, expected) {
			t.Errorf("OptionsForDir returned %+v, did not match expectation %+v", opts, expected)
		}
	}

	dir = getDir(t, "../testdata/linter/badcfgerrors")
	if _, err := OptionsForDir(dir); err == nil {
		t.Errorf("Expected an error from OptionsForDir, but it was nil")
	} else {
		if _, ok := err.(ConfigError); !ok {
			t.Errorf("Expected error to be a ConfigError, but instead type is %T", err)
		}
		if !strings.HasPrefix(err.Error(), "Option errors ") {
			t.Errorf("Error message does not contain expected prefix. Message: %s", err.Error())
		}
	}

	dir = getDir(t, "../testdata/linter/badcfgwarnings")
	if _, err := OptionsForDir(dir); err == nil {
		t.Errorf("Expected an error from OptionsForDir, but it was nil")
	} else {
		if _, ok := err.(ConfigError); !ok {
			t.Errorf("Expected error to be a ConfigError, but instead type is %T", err)
		}
		if !strings.HasPrefix(err.Error(), "Option warnings ") {
			t.Errorf("Error message does not contain expected prefix. Message: %s", err.Error())
		}
	}
}

func getValidConfig(t *testing.T) *mybase.Config {
	cmd := mybase.NewCommand("lintertest", "", "", nil)
	AddCommandOptions(cmd)
	cmd.AddArg("environment", "production", false)
	return mybase.ParseFakeCLI(t, cmd, "lintertest")
}

func getDir(t *testing.T, dirPath string) *fs.Dir {
	t.Helper()
	dir, err := fs.ParseDir(dirPath, getValidConfig(t))
	if err != nil {
		t.Fatalf("Unexpected error parsing dir %s: %s", dirPath, err)
	}
	return dir
}
