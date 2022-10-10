// Copyright (c) 2022 Erik Kinnunen.
// license can be found in the LICENSE file.

/*
	Package miniflag implements command-line flag parsing. miniflag offers
	minimal API on top of go flag module found in the standard library.

	The minimal API is created by using the generics functionality added in
	1.18. The main difference between the flag module in standard library and
	miniflag is that all the flags defined in miniflag are setup using just a
	single function. For example creating integer and boolean values in
	standard library requires using two different function calls `Int()` and
	`Bool()`, whereas in miniflag you only use one `Flag`.

	With miniflag shorthands are created as separate flag definitions, but will
	hold the reference to same variable.

	Usage

	Define flags using miniflag.Flag() or miniflag.SetFlag().

	Note: Taken from go flag module in standard library with variation to match
	miniflag implementation.

	This declares an integer flag, -n, stored in the pointer nFlag, with type *int:

		import "github.com/erikjuhani/miniflag"
		var nFlag = miniflag.Flag("n", "", 1234, "help message for flag n")

	Or you can create custom flags that satisfy the Value interface and couple
	them to flag parsing by:

		var namesFlag = miniflag.Flag("names", "n", StringSliceFlag{}, "help message for names flag")

	For such flags, the default value is just the initial value of the
	variable.

	After all flags are defined, call:

		flag.Parse()

	to parse the command line into the defined flags.

	Flags may then be used directly. All flag values are pointers.

		fmt.Println("nFlag has value ", *nFlag)

	After parsing, the arguments following the flags are available as the slice
	flag.Args().
*/
package miniflag

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

// ErrorHandling defines how FlagSet.Parse behaves if the parse fails.
// NOTE: Direct reference to std.
type ErrorHandling = flag.ErrorHandling

// These constants cause FlagSet.Parse to behave as described if the parse
// fails.
// NOTE: Direct reference to std.
const (
	ContinueOnError ErrorHandling = flag.ContinueOnError // Return a descriptive error.
	ExitOnError                   = flag.ExitOnError     // Call os.Exit(2) or for -h/-help Exit(0).
	PanicOnError                  = flag.PanicOnError    // Call panic with a descriptive error.
)

var (
	// CommandLine is the default set of command-line flags, parsed from
	// os.Args.
	CommandLine = NewFlagSet(os.Args[0], ExitOnError)
	// Setup capacity for optimized performance
	flagSets = make(map[string]FlagSet[any], 8)
	// Increase performance by pre-allocating slice capacity
	// flagInfoSlice is used in new FlagSet creation
	flagInfoSlice = make([]flagInfo, 0, 32)
)

// A FlagSet represents a set of defined flags. The zero value of a FlagSet has
// no name and has ContinueOnError error handling.
//
// Flag names must be unique within a FlagSet. An attempt to define a flag
// whose name is already in use will cause a panic.
// NOTE: Direct reference to standard lib.
type FlagSet[T any] struct {
	*flag.FlagSet
	flags []flagInfo
	// TODO: move flagSets into FlagSet
	// Sub flagsets are commands for parent flagset/command
}

func (fs *FlagSet[T]) Usage() {
	usageFn(fs, fs.Name())
}

// SetFlag defines a new flag to a given FlagSet.
func SetFlag[T any](fs *FlagSet[any], name string, shorthand string, value T, usage string) *T {
	return defineFlag(fs, name, shorthand, value, usage)
}

// Flag defines a new flag for CommandLine with the given name, shorthand, usage
// and value. Value type is inferred from the given value. Shorthand for the
// flag is only created if passed shorthand parameter is not an empty string.
func Flag[T any](name string, shorthand string, value T, usage string) *T {
	return defineFlag(CommandLine, name, shorthand, value, usage)
}

// NewFlagSet returns a new, empty flag set with the specified name and error
// handling property.
func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet[any] {
	fs := FlagSet[any]{flag.NewFlagSet(name, errorHandling), flagInfoSlice}
	flagSets[name] = fs
	return &fs
}

// Args returns non-flag arguments.
func Args() []string {
	return args(CommandLine)
}

func Parse() error {
	return parse(CommandLine, os.Args[1:])
}

// flagInfo stores flag information and is used internally.
type flagInfo struct {
	Longhand   string
	Shorthand  string
	UsageValue string
	Usage      string
}

func parse(fs *FlagSet[any], args []string) error {
	l := len(args)
	if l > 1 {
		if f, ok := flagSets[args[0]]; ok {
			return f.Parse(args[1:])
		}
	}
	return fs.Parse(args)
}

func args(fs *FlagSet[any]) []string {
	args := fs.Args()

	pArgs := []string{}
	for i, arg := range args {
		if arg == "" {
			pArgs = append(pArgs, arg)
			continue
		}

		if arg[0] == '-' {
			continue
		}
		if i > 0 && args[i-1][0] == '-' {
			f := fs.Lookup(strings.ReplaceAll(args[i-1], "-", ""))

			if f != nil && reflect.TypeOf(f.Value).Elem().Kind() != reflect.Bool {
				continue
			}
		}
		pArgs = append(pArgs, arg)
	}

	return pArgs
}

func defineFlag[T any](fs *FlagSet[any], name string, shorthand string, value T, usage string) *T {
	if name == shorthand {
		shorthand = ""
	}

	defineUsage(&fs.flags, name, shorthand, usage)

	switch v := any(value).(type) {
	case bool:
		return any((boolVar(fs, name, shorthand, v, usage))).(*T)
	case string:
		return any(stringVar(fs, name, shorthand, v, usage)).(*T)
	case int:
		return any(intVar(fs, name, shorthand, v, usage)).(*T)
	case int64:
		return any(int64Var(fs, name, shorthand, v, usage)).(*T)
	case uint:
		return any(uintVar(fs, name, shorthand, v, usage)).(*T)
	case uint64:
		return any(uint64Var(fs, name, shorthand, v, usage)).(*T)
	case float64:
		return any(float64Var(fs, name, shorthand, v, usage)).(*T)
	case time.Duration:
		return any(durationVar(fs, name, shorthand, v, usage)).(*T)
	case T:
		return valueVar(fs, name, shorthand, v, usage)
	}
	return nil
}

func usageFn[T any](fs *FlagSet[T], name string) {
	var s, u strings.Builder

	s.WriteString("usage: " + name)

	p := s.Len()

	for i, f := range fs.flags {
		var c strings.Builder

		if f.Shorthand == "" && f.Longhand == "" {
			continue
		}

		if f.Shorthand != "" {
			fmt.Fprintf(&c, "-%s", f.Shorthand)
		}

		if f.Longhand != "" {
			if f.Shorthand != "" {
				c.WriteRune(' ')
			}
			fmt.Fprintf(&c, "--%s", f.Longhand)
		}

		compound := c.String()

		if f.UsageValue != "" {
			fmt.Fprintf(&s, " [%s=%s]", compound, f.UsageValue)
		} else {
			fmt.Fprintf(&s, " [%s]", compound)
		}

		if (i+1)%4 == 0 {
			fmt.Fprintf(&s, "\n%*s", p, "")
		}

		fmt.Fprintf(
			&u,
			"%*s%*s\n",
			len(compound)+4,
			compound,
			len(f.Usage)-len(compound)+16,
			f.Usage,
		)
	}

	fmt.Fprint(fs.Output(), s.String(), "\n", u.String())
}

// TODO: Maybe this can be optimized
// Extract usage value from the usage text. Looks for the first occurrance
// between backtick characters
func extractUsageValue(s string) string {
	var pos int

	for i := 0; i < len(s); i++ {
		if s[i] == '`' {
			if pos == 0 {
				pos = i + 1
				continue
			}
			return s[pos:i]
		}
	}

	return ""
}

func defineUsage(flags *[]flagInfo, name string, shorthand string, usage string) {
	if len(name) == 1 {
		shorthand = name
		name = ""
	}

	*flags = append(
		*flags,
		flagInfo{
			Longhand:   name,
			Shorthand:  shorthand,
			Usage:      usage,
			UsageValue: extractUsageValue(usage),
		},
	)
}

func boolVar(fs *FlagSet[any], name string, shorthand string, value bool, usage string) *bool {
	fs.BoolVar(&value, name, value, usage)
	if shorthand != "" {
		fs.BoolVar(&value, shorthand, value, usage)
	}
	return &value
}

func stringVar(fs *FlagSet[any], name string, shorthand string, value string, usage string) *string {
	fs.StringVar(&value, name, value, usage)
	if shorthand != "" {
		fs.StringVar(&value, shorthand, value, usage)
	}
	return &value
}

func intVar(fs *FlagSet[any], name string, shorthand string, value int, usage string) *int {
	fs.IntVar(&value, name, value, usage)
	if shorthand != "" {
		fs.IntVar(&value, shorthand, value, usage)
	}
	return &value
}

func int64Var(fs *FlagSet[any], name string, shorthand string, value int64, usage string) *int64 {
	fs.Int64Var(&value, name, value, usage)
	if shorthand != "" {
		fs.Int64Var(&value, shorthand, value, usage)
	}
	return &value
}

func uintVar(fs *FlagSet[any], name string, shorthand string, value uint, usage string) *uint {
	fs.UintVar(&value, name, value, usage)
	if shorthand != "" {
		fs.UintVar(&value, shorthand, value, usage)
	}
	return &value
}

func uint64Var(fs *FlagSet[any], name string, shorthand string, value uint64, usage string) *uint64 {
	fs.Uint64Var(&value, name, value, usage)
	if shorthand != "" {
		fs.Uint64Var(&value, shorthand, value, usage)
	}
	return &value
}

func float64Var(fs *FlagSet[any], name string, shorthand string, value float64, usage string) *float64 {
	fs.Float64Var(&value, name, value, usage)
	if shorthand != "" {
		fs.Float64Var(&value, shorthand, value, usage)
	}
	return &value
}

func durationVar(fs *FlagSet[any], name string, shorthand string, value time.Duration, usage string) *time.Duration {
	fs.DurationVar(&value, name, value, usage)
	if shorthand != "" {
		fs.DurationVar(&value, shorthand, value, usage)
	}
	return &value
}

func valueVar[T any](fs *FlagSet[any], name string, shorthand string, value T, usage string) *T {
	fs.Var(any(&value).(flag.Value), name, usage)
	if shorthand != "" {
		fs.Var(any(&value).(flag.Value), shorthand, usage)
	}
	return &value
}
