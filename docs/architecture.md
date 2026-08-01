Original Picomatch Repository Architecture
picomatch/
│
├── bench/
│
├── examples/
│
├── lib/
│
├── test/
│
├── index.js
├── package.json
├── README.md
└── LICENSE

The important part is:

index.js
    |
    v
lib/picomatch.js
    |
    +----------------+
    |                |
    v                v
 lib/parse.js     lib/scan.js
    |
    v
 Regex Compiler
    |
    v
 Matcher
1. Root Files
index.js
Purpose

Public package entry point.

Current role:

module.exports = require('./lib/picomatch');

It exposes the Picomatch API.

Example:

const picomatch = require("picomatch");

const matcher = picomatch("*.js");

Architecture equivalent:

User
 |
 |
index.js
 |
 |
lib/picomatch.js
2. lib/

This is the core engine.

lib/
|
├── picomatch.js
├── parse.js
├── scan.js
├── constants.js
└── utils.js

These files contain the actual implementation.

lib/picomatch.js
Main Engine API

This is the heart of Picomatch.

Responsibilities:

create matcher functions
manage options
call parser
create regex
expose APIs

Flow:

picomatch(pattern)

        |
        v

parse(pattern)

        |
        v

makeRe()

        |
        v

RegExp matcher

Example:

const isMatch = picomatch("*.js");

isMatch("app.js");
// true
lib/scan.js
Glob Scanner

Responsible for:

picomatch.scan()

It performs lexical analysis.

Example:

Input:

!src/**/*.ts

Output:

{
 prefix:"!",
 base:"src",
 glob:"**/*.ts",
 isGlob:true,
 negated:true
}

Architecture:

Raw Pattern

      |
      v

Scanner

      |
      v

Pattern Metadata

It does NOT create regex.

lib/parse.js
Glob Parser

Responsible for converting glob syntax into regex source.

Example:

Input:

*.js

Parser understands:

*
 |
 |
anything except slash

.js
 |
 |
literal extension

Output:

^(?:(?!\.)[^/]*?\.js)$

Architecture:

Glob Tokens

      |
      v

Parser

      |
      v

Regex Source
lib/constants.js

Contains shared parser constants.

Examples:

STAR
QMARK
SLASH
DOT

Instead of:

if(char==="*")

everywhere:

Picomatch uses:

constants.STAR

Purpose:

avoid duplicated values
keep parser readable
maintain consistency
lib/utils.js

Helper functions.

Examples:

string utilities
regex helpers
path utilities
validation

Used by:

picomatch.js
parse.js
scan.js
3. examples/

These are API usage examples.

Structure:

examples/

├── extglob.js
├── extglob-negated.js
├── makeRe.js
├── match.js
├── scan.js
├── windows.js
├── option-ignore.js
├── option-onMatch.js
├── option-onIgnore.js
├── option-onResult.js
└── option-expandRange.js

These demonstrate:

Matching
match.js

Example:

picomatch.isMatch(
"foo.js",
"*.js"
)
Regex generation
makeRe.js

Example:

picomatch.makeRe("*.js")

Output:

/^(?:(?!\.)(?=.)[^/]*?\.js)$/
Scanner
scan.js

Example:

picomatch.scan(
"src/**/*.js"
)
Options

Examples:

option-ignore.js
option-onMatch.js
option-onIgnore.js
option-onResult.js
option-expandRange.js

These should become Go tests/examples.

4. test/

Testing architecture.

The port MUST preserve this philosophy.

Original:

test/

├── matching tests
├── parsing tests
├── options tests
├── edge cases

Port equivalent:

tests/

├── scan_test.go
├── parse_test.go
├── match_test.go
├── options_test.go
└── regression_test.go
5. bench/

Performance testing.

Contains:

first-match-minimatch.js
first-match-picomatch.js
glob-parent.js
load-time.js

Purpose:

Compare:

minimatch
    vs
picomatch

Measures:

startup time
regex compilation
matching speed

Port equivalent:

Go:

bench/

├── match_test.go
├── compile_test.go
└── benchmark_test.go

Using:

go test -bench=.
6. Complete Architecture Mapping

Original:

                 index.js
                    |
                    v
             picomatch.js
                    |
        +-----------+-----------+
        |                       |
        v                       v
     scan.js                parse.js
        |                       |
        |                       v
        |                 Regex Generator
        |
        v
 Pattern Metadata

                    |
                    v

              Matcher Function

                    |
                    v

              User Input Path
Go Port Architecture

Preserve the same boundaries:

picomatch-go/

│
├── cmd/
│
├── internal/
│   |
│   ├── picomatch/
│   │
│   ├── scanner/
│   │     └── scan.go
│   │
│   ├── parser/
│   │     └── parse.go
│   │
│   ├── compiler/
│   │     └── regex.go
│   │
│   ├── matcher/
│   │     └── matcher.go
│   │
│   ├── constants/
│   │     └── constants.go
│   │
│   └── utils/
│         └── utils.go
│
├── examples/
│
├── bench/
│
├── test/
│
├── go.mod
└── README.md
Important Porting Rule

Do not merge:

scan + parse + compile

into:

Match(pattern,input)

because original Picomatch separates them.

The equivalent should remain:

pattern
   |
   v
Scan()
   |
   v
Parse()
   |
   v
CompileRegex()
   |
   v
Match()

This preserves:

API behavior
debugging capability
test mapping
benchmark comparison
future compatibility

For Port Mortem, this is the architecture reviewers expect to see: same engine stages, different language implementation.