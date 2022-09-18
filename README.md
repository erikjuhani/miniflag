## Description

`miniflag` is minimal implementation done on top of go's standard flag library.
It leverages go generics and provides simpler more minimal API for creating
flags.

The minimal API is created by using the generics functionality added in
1.18. The main difference between the flag module in standard library and
miniflag is that all the flags defined in miniflag are setup using just a
single function. For example creating integer and boolean values in
standard library requires using two different function calls `Int()` and
`Bool()`, whereas in miniflag you only use one `Flag()`.

With miniflag shorthands are created as separate flag definitions, but will
hold the pointer reference to the same variable.

Example of main difference in API:

```diff
import (
-	"flag"
+	flag "github.com/erikjuhani/miniflag"
)

-var stringSliceFlag StringSliceFlag

var (
-	enable   = flag.Bool("enable", false, "description for enable flag")
-	name     = flag.String("name", "", "description for name flag")
-	custom   = flag.Var(&stringSliceFlag, "custom", "description for custom flag")
+	enable   = flag.Flag("enable", "e", false, "description for enable flag")
+	name     = flag.Flag("name", "n", "", "description for name flag")
+	custom   = flag.Flag("custom", "c", StringSliceFlag{}, "description for custom flag")
)

type StringSliceFlag []string

func (s *StringSliceFlag) String() string {
	return fmt.Sprintf("%s", *s)
}

func (s *StringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}
```

## Requirements

The minimum go version required is 1.18.

## Installation

Install miniflag by running:

```
go get github.com/erikjuhani/miniflag
```

Run tests by running:

```
go test github.com/erikjuhani/miniflag
```

## Usage

Define flags using `miniflag.Flag()` or `miniflag.SetFlag()`. The API is
minimalistic in purpose so the amount of noise is as little as possible. Types
can be inferred from the given default value, which is always required.

This declares an integer flag, -n, stored in the pointer nFlag, with type *int:

```go
import "github.com/erikjuhani/miniflag"

var nFlag = miniflag.Flag("n", "", 1234, "help message for flag n")
// Inferred as *int type
```

Or you can create custom flags that satisfy the Value interface and couple
them to flag parsing by:

```go
var namesFlag = miniflag.Flag("names", "n", StringSliceFlag{}, "help message for names flag")
// Inferred as *StringSliceFlag
```

For such flags, the default value is just the initial value of the
variable.

After all flags are defined, call:

```go
miniflag.Parse()
```

to parse the command line into the defined flags.

### Setting flags to flagsets

By default flags are defined to CommandLine ie using `miniflag.Flag` function.
FlagSet or command specific flags can be defined by using separate FlagSets.
New flag sets can be created by using `miniflag.NewFlagSet()` function.

To define flags to specific FlagSet `miniflag.SetFlag` can be used. The first
argument is the pointer to the flag set. 

```go
subCmd := miniflag.NewFlagSet("subCmd", flag.ExitOnError)
var (
    subCmdEnable = miniflag.SetFlag(subCmd, "enable", "e", false, "help message for subCmd enable flag")
    subCmdName   = miniflag.SetFlag(subCmd, "name", "n", "", "help message for subCmd name flag")
)
```

### Help usage

`miniflag` has a default custom help usage messsage, which takes inspiration
fromhelp usage printed by `git` command-line tool. Value text can be changed by
using backticks `<value>` in the flags help message.

Example of custom default help message:

```
usage: cmd [-e --enable=<bool>] [-n --name]
    -e --enable     help message for cmd enable `<bool>` flag
    -n --name       help message for cmd name flag
```
