# TODO before version 1.0.0

This document is a place to organize my thoughts. Everything is subject to
change.

## Extensions api

Allow users to register their own extensions for commands that cannot be achieved through the basic macro system. This
should look something like:
```go
type CmdArg struct{
    toks []Token
    kind CmdArgKind
}

type ExtFunc func(n *MMLNode, ctx parseContext, args []CmdArg) error

// RegisterExtension associates all commands provided in the slice with the function f.
// the 'consumes' argument behaves like a list of regular expressions that
// define what types of arguments the command may accept.
func RegisterExtension(commands []string, f ExtFunc, consumes []CmdArgExpr)
```

Take for example the `\dv, \pdv, \fdv...` family of commands for derivatives.
They may accept an `[optional argument]` of comma-separated values, followed by
  1. a single argument enclosed in {curly braces} (let's call this a *grouped
     arg*); can be comma-separated
  2. two grouped args, the second can be comma-separated
  3. a grouped arg, a literal '/', a grouped arg (comma-separated)
  4. a grouped arg, a literal '!', a grouped arg (comma-separated)
  5. a grouped arg, a literal '!', a literal '/', a grouped arg
     (comma-separated)
It would be convienient to lift ideas from regular expressions to express these
requirements.

My initial idea for the `CmdArg` and `CmdArgExpr` types is simple:
```go
type CmdArgExpr

type CmdFlavor uint64

const (
    CF_LITERAL CmdFlavor = iota << 1
    CF_GROUPED
    CF_TOKEN
    CF_OPTION
    CF_COMMA_SEP
)

type CmdArg struct{
    Flavor CmdFlavor
    Value []Token
}
```

then those 5 patterns could be expressed as something like

```go
dvConsumes := []CmdArgExpr{
    "[,]{,}",
    "[,]{}{,}",
    "[,]{}'/'{,}",
    "[,]{}'!'{,}",
    "[,]{}'/''!'{,}",
}
```

and if we continue to borrow from regexes,

```go
dvConsumes := []CmdArgExpr{
    "[,]?{}?'/'?'!'?{,}"
}
```

which would correspond to the sequence
`(0 or 1 comma-separated option), (0 or 1 grouped expression), (0 or 1 '/' literal), (0 or 1 '!' literal), (required comma-separated grouped expression)`
