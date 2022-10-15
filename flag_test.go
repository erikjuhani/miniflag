package miniflag

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	tests := []struct {
		args     []string
		expected bool
	}{
		{},
		{
			args:     []string{"-b"},
			expected: true,
		},
		{
			args:     []string{"foo", "-b"},
			expected: true,
		},
	}

	for _, tt := range tests {
		fs := NewFlagSet("foo", ContinueOnError)
		t.Run("", func(t *testing.T) {
			actual := SetFlag(fs, "b", "", false, "")

			if err := parse(fs, tt.args); err != nil {
				t.Fatal(err)
			}

			if tt.expected != *actual {
				t.Fatalf("flag value did not match expected %t, got %t", tt.expected, *actual)
			}
		})
	}
}

func TestNameIsShorthand(t *testing.T) {
	tests := []struct {
		args     []string
		expected bool
	}{
		{
			args:     []string{"-b"},
			expected: true,
		},
	}

	for _, tt := range tests {
		fs := NewFlagSet("", ContinueOnError)
		t.Run("", func(t *testing.T) {
			actual := SetFlag(fs, "b", "b", false, "Test name is shorthand flag")

			if err := fs.Parse(tt.args); err != nil {
				t.Fatal(err)
			}

			if tt.expected != *actual {
				t.Fatalf("flag value did not match expected %t, got %t", tt.expected, *actual)
			}
		})
	}
}

func TestBool(t *testing.T) {
	tests := []struct {
		args     []string
		expected bool
	}{
		{
			expected: false,
		},
		{
			args:     []string{"-b"},
			expected: true,
		},
		{
			args:     []string{"--bool"},
			expected: true,
		},
		{
			args:     []string{"--bool=false"},
			expected: false,
		},
	}

	for _, tt := range tests {
		fs := NewFlagSet("", ContinueOnError)
		t.Run("", func(t *testing.T) {
			actual := SetFlag(fs, "bool", "b", false, "Test bool flag")

			if err := fs.Parse(tt.args); err != nil {
				t.Fatal(err)
			}

			if tt.expected != *actual {
				t.Fatalf("flag value did not match expected %t, got %t", tt.expected, *actual)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		args     []string
		expected string
	}{
		{
			expected: "",
		},
		{
			args:     []string{"-s", "string"},
			expected: "string",
		},
		{
			args:     []string{"--string", "long", "string"},
			expected: "long",
		},
	}

	for _, tt := range tests {
		fs := NewFlagSet("", ContinueOnError)
		t.Run("", func(t *testing.T) {
			actual := SetFlag(fs, "string", "s", "", "Test string flag")

			if err := fs.Parse(tt.args); err != nil {
				t.Fatal(err)
			}

			if tt.expected != *actual {
				t.Fatalf("flag value did not match expected %s, got %s", tt.expected, *actual)
			}
		})
	}
}

func TestInt(t *testing.T) {
	tests := []struct {
		args     []string
		expected int
	}{
		{
			expected: 0,
		},
		{
			args:     []string{"-i", "1"},
			expected: 1,
		},
		{
			args:     []string{"--int", "-1"},
			expected: -1,
		},
	}

	for _, tt := range tests {
		fs := NewFlagSet("", ContinueOnError)
		t.Run("", func(t *testing.T) {
			actual := SetFlag(fs, "int", "i", 0, "Test int flag")

			if err := fs.Parse(tt.args); err != nil {
				t.Fatal(err)
			}

			if tt.expected != *actual {
				t.Fatalf("flag value did not match expected %d, got %d", tt.expected, *actual)
			}
		})
	}
}

func TestInt64(t *testing.T) {
	tests := []struct {
		args     []string
		expected int64
	}{
		{
			expected: 0,
		},
		{
			args:     []string{"-i", "1"},
			expected: 1,
		},
		{
			args:     []string{"--int64", "-1"},
			expected: -1,
		},
	}

	for _, tt := range tests {
		fs := NewFlagSet("", ContinueOnError)
		t.Run("", func(t *testing.T) {
			actual := SetFlag(fs, "int64", "i", int64(0), "Test int64 flag")

			if err := fs.Parse(tt.args); err != nil {
				t.Fatal(err)
			}

			if tt.expected != *actual {
				t.Fatalf("flag value did not match expected %d, got %d", tt.expected, *actual)
			}
		})
	}
}

func TestUint(t *testing.T) {
	tests := []struct {
		args     []string
		expected uint
	}{
		{
			expected: 0,
		},
		{
			args:     []string{"-u", "1"},
			expected: 1,
		},
		{
			args:     []string{"--uint", "10"},
			expected: 10,
		},
	}

	for _, tt := range tests {
		fs := NewFlagSet("", ContinueOnError)
		t.Run("", func(t *testing.T) {
			actual := SetFlag(fs, "uint", "u", uint(0), "Test uint flag")

			if err := fs.Parse(tt.args); err != nil {
				t.Fatal(err)
			}

			if tt.expected != *actual {
				t.Fatalf("flag value did not match expected %d, got %d", tt.expected, *actual)
			}
		})
	}
}

func TestNil(t *testing.T) {
	tests := []struct {
		expected any
	}{
		{
			expected: nil,
		},
	}

	for _, tt := range tests {
		fs := NewFlagSet("", ContinueOnError)
		t.Run("", func(t *testing.T) {
			actual := SetFlag[any](fs, "nil", "n", nil, "Test nil flag")

			if actual != nil {
				t.Fatalf("flag value did not match expected %v, got %v", tt.expected, actual)
			}
		})
	}
}

func TestUint64(t *testing.T) {
	tests := []struct {
		args     []string
		expected uint64
	}{
		{
			expected: 0,
		},
		{
			args:     []string{"-u", "1"},
			expected: 1,
		},
		{
			args:     []string{"--uint64", "10"},
			expected: 10,
		},
	}

	for _, tt := range tests {
		fs := NewFlagSet("", ContinueOnError)
		t.Run("", func(t *testing.T) {
			actual := SetFlag(fs, "uint64", "u", uint64(0), "Test uint64 flag")

			if err := fs.Parse(tt.args); err != nil {
				t.Fatal(err)
			}

			if tt.expected != *actual {
				t.Fatalf("flag value did not match expected %d, got %d", tt.expected, *actual)
			}
		})
	}
}

func TestFloat64(t *testing.T) {
	tests := []struct {
		args     []string
		expected float64
	}{
		{
			args:     []string{},
			expected: 0.000000,
		},
		{
			args:     []string{"-f", "1.0"},
			expected: 1.000000,
		},
		{
			args:     []string{"--float64", "10.000001"},
			expected: 10.000001,
		},
	}

	for _, tt := range tests {
		fs := NewFlagSet("", ContinueOnError)
		t.Run("", func(t *testing.T) {
			actual := SetFlag(fs, "float64", "f", float64(0), "Test float64 flag")

			if err := fs.Parse(tt.args); err != nil {
				t.Fatal(err)
			}

			if tt.expected != *actual {
				t.Fatalf("flag value did not match expected %f, got %f", tt.expected, *actual)
			}
		})
	}
}

func TestDuration(t *testing.T) {
	tests := []struct {
		args     []string
		expected time.Duration
	}{
		{
			args:     []string{},
			expected: 0,
		},
		{
			args:     []string{"-d", "1ns"},
			expected: 1,
		},
		{
			args:     []string{"--duration", "1ms"},
			expected: 1000000,
		},
	}

	for _, tt := range tests {
		fs := NewFlagSet("", ContinueOnError)
		t.Run("", func(t *testing.T) {
			actual := SetFlag(fs, "duration", "d", time.Duration(0), "Test duration flag")

			if err := fs.Parse(tt.args); err != nil {
				t.Fatal(err)
			}

			if tt.expected != *actual {
				t.Fatalf("flag value did not match expected %d, got %d", tt.expected, *actual)
			}
		})
	}
}

type customSliceFlag []string

func (s *customSliceFlag) String() string {
	return fmt.Sprintf("%s", *s)
}

func (s *customSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func TestCustom(t *testing.T) {
	tests := []struct {
		args     []string
		expected customSliceFlag
	}{
		{
			args:     []string{},
			expected: customSliceFlag{},
		},
		{
			args:     []string{"-s", "A"},
			expected: customSliceFlag{"A"},
		},
		{
			args:     []string{"--slice", "A", "--slice", "B"},
			expected: customSliceFlag{"A", "B"},
		},
	}

	for _, tt := range tests {
		fs := NewFlagSet("", ContinueOnError)
		t.Run("", func(t *testing.T) {
			actual := SetFlag(fs, "slice", "s", customSliceFlag{}, "Test customSliceFlag flag")

			if err := fs.Parse(tt.args); err != nil {
				t.Fatal(err)
			}

			if len(tt.expected) != len(*actual) {
				t.Fatalf("flag value did not match expected %q, got %q", tt.expected, *actual)
			}
		})
	}
}

func TestArgs(t *testing.T) {
	tests := []struct {
		args     []string
		expected []string
	}{
		{
			args:     []string{},
			expected: []string{},
		},
		{
			args:     []string{""},
			expected: []string{""},
		},
		{
			args:     []string{"-s", "string"},
			expected: []string{},
		},
		{
			args:     []string{"--bool", "arg0"},
			expected: []string{"arg0"},
		},
		{
			args:     []string{"arg0", "--bool", "arg1", "-s", "string", "arg2"},
			expected: []string{"arg0", "arg1", "arg2"},
		},
	}

	for _, tt := range tests {
		fs := NewFlagSet("", ContinueOnError)
		t.Run("", func(t *testing.T) {
			// setup flags
			// boolean is a special case so that needs to be tested
			SetFlag(fs, "bool", "b", false, "bool flag")
			SetFlag(fs, "string", "s", "", "string flag")

			if err := fs.Parse(tt.args); err != nil {
				t.Fatal(err)
			}

			actual := args(fs)

			if len(tt.expected) != len(actual) {
				t.Fatalf("flag value did not match expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestUsage(t *testing.T) {
	tests := []struct {
		flags    []flagInfo
		expected string
	}{
		{
			flags:    []flagInfo{},
			expected: "usage: test\n",
		},
		{
			flags: []flagInfo{
				{},
			},
			expected: "usage: test\n",
		},
		{
			flags: []flagInfo{
				{},
				{
					Longhand:   "test",
					Shorthand:  "t",
					Usage:      "Usage for `<test>`",
					UsageValue: "<test>",
				},
			},
			expected: "usage: test [-t --test=<test>]\n    -t --test       Usage for `<test>`\n",
		},
		{
			flags: []flagInfo{
				{
					Longhand:  "aa",
					Shorthand: "a",
					Usage:     "Usage for a",
				},
				{
					Longhand:  "bb",
					Shorthand: "b",
					Usage:     "Usage for b",
				},
				{
					Longhand:  "cc",
					Shorthand: "c",
					Usage:     "Usage for c",
				},
				{
					Longhand:  "dd",
					Shorthand: "d",
					Usage:     "Usage for d",
				},
				{
					Longhand:   "ee",
					Shorthand:  "e",
					Usage:      "Usage for e",
					UsageValue: "bool",
				},
			},
			expected: `usage: test [-a --aa] [-b --bb] [-c --cc] [-d --dd]
            [-e --ee=bool]
    -a --aa         Usage for a
    -b --bb         Usage for b
    -c --cc         Usage for c
    -d --dd         Usage for d
    -e --ee         Usage for e
`,
		},
	}

	for _, tt := range tests {
		var b bytes.Buffer
		fs := NewFlagSet("test", ContinueOnError)
		fs.SetOutput(&b)
		fs.flags = tt.flags
		fs.Usage()

		t.Run("", func(t *testing.T) {
			actual := b.String()

			if tt.expected != actual {
				t.Fatalf("Help string did not match expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestDefineUsage(t *testing.T) {
	tests := []struct {
		name      string
		shorthand string
		usage     string
		actual    []flagInfo
		expected  flagInfo
	}{
		{},
		{
			name:      "test",
			shorthand: "t",
			usage:     "Adjust with `bool` value",
			expected: flagInfo{
				Longhand:   "test",
				Shorthand:  "t",
				Usage:      "Adjust with `bool` value",
				UsageValue: "bool",
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			defineUsage(&tt.actual, tt.name, tt.shorthand, tt.usage)

			actual := tt.actual[0]

			if tt.expected != actual {
				t.Fatalf("flag usage did not match expected %q, got %q", tt.expected, actual)
			}
		})
	}
}
