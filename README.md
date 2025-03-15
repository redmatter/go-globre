# Glob to RegExp

[![GoDoc](https://godoc.org/github.com/redmatter/go-globre?status.svg)](https://godoc.org/github.com/redmatter/go-globre)
[![Go Report Card](https://goreportcard.com/badge/github.com/redmatter/go-globre)](https://goreportcard.com/report/github.com/redmatter/go-globre)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

A library for converting glob strings to regular expressions in golang.

# Usage

```go
import (
    "github.com/redmatter/go-globre/v2"
    "regexp"
)

urlGlob := "{http,https}://example.com/**"
// converts to ^(http|https):\/\/example\.com\/(?:(?:[^/]*(?:/|$))*)$
urlRe := globre.RegexFromGlob(
	urlGlob,
    globre.ExtendedSyntaxEnabled(true),
    globre.GlobStarEnabled(true),
    globre.WithDelimiter('/'),
)

// Compile the regexp
urlRegexp := regexp.MustCompile(urlRe)

// Match the url
if urlRegexp.MatchString("https://example.com/index.htm") {
	// do something
}
```

## Options

The conversion to regexp can be modified using the following Option functions:

| Option                | Default Value | Description                                        |
|:----------------------|---------------|:---------------------------------------------------|
| WithDelimiter         | '/'           | Delimiter used for tokenising the compared string  |
| ExtendedSyntaxEnabled | false         | Enable extended globs, supporting classes, etc.    |
| GlobStarEnabled       | false         | Enables "double-star" match for one or more tokens |

# License

Please see [LICENSE](LICENSE) and [LICENSE-THIRD-PARTY](LICENSE-THIRD-PARTY.md).
