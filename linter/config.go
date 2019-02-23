package linter

import (
	"fmt"
	"strings"

	"github.com/skeema/mybase"
	"github.com/skeema/skeema/fs"
)

// Severity represents different annotation severity levels.
type Severity string

// Constants enumerating valid severity levels
const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

// AddCommandOptions adds linting-related mybase options to the supplied
// mybase.Command.
func AddCommandOptions(cmd *mybase.Command) {
	cmd.AddOption(mybase.StringOption("warnings", 0, "no-pk,bad-charset,bad-engine", "Linter problems to display as warnings (non-fatal); see manual for usage"))
	cmd.AddOption(mybase.StringOption("errors", 0, "", "Linter problems to treat as fatal errors; see manual for usage"))
	cmd.AddOption(mybase.StringOption("allow-charset", 0, "latin1,utf8mb4", "Whitelist of acceptable character sets"))
	cmd.AddOption(mybase.StringOption("allow-engine", 0, "innodb", "Whitelist of acceptable storage engines"))
}

// Options contains parsed settings controlling linter behavior.
type Options struct {
	ProblemSeverity map[string]Severity
	AllowedCharSets []string
	AllowedEngines  []string
}

// OptionsForDir returns Options based on the configuration in an fs.Dir,
// effectively converting between mybase options and linter options.
func OptionsForDir(dir *fs.Dir) (Options, error) {
	opts := Options{
		ProblemSeverity: make(map[string]Severity),
		AllowedCharSets: dir.Config.GetSlice("allow-charset", ',', true),
		AllowedEngines:  dir.Config.GetSlice("allow-engine", ',', true),
	}

	allAllowed := strings.Join(allProblemNames(), ", ")
	for _, val := range dir.Config.GetSlice("warnings", ',', true) {
		val = strings.ToLower(val)
		if !problemExists(val) {
			return opts, ConfigError(fmt.Sprintf("Option warnings must be a comma-separated list including these values: %s", allAllowed))
		}
		opts.ProblemSeverity[val] = SeverityWarning
	}
	for _, val := range dir.Config.GetSlice("errors", ',', true) {
		val = strings.ToLower(val)
		if !problemExists(val) {
			return opts, ConfigError(fmt.Sprintf("Option errors must be a comma-separated list including these values: %s", allAllowed))
		}
		opts.ProblemSeverity[val] = SeverityError
	}

	return opts, nil
}

// ConfigError represents a configuration problem encountered at runtime.
type ConfigError string

// Error satisfies the builtin error interface.
func (ce ConfigError) Error() string {
	return string(ce)
}
