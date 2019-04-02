# Glob to RegExp

Glob string to regexp string conversion in golang 

# Usage

```go
package example

urlGlob := "{http,https}://example.com/**"
// converts to ^(http|https):\/\/example\.com\/(?:(?:[^/]*(?:/|$))*)$
urlRe := globre.RegexFromGlobWithOptions(
	urlGlob, globre.Options{
		Extended:  true,
		GlobStar:  true,
		Delimiter: '/',
	})
urlRegexp := regexp.MustCompile(urlRe)
urlRegexp.MatchString("https://example.com/index.htm")
```

## Options

The conversion to regexp can be modified using `globre.Options`.

|   Attr    | Type |                        Desc                        |
| :-------- | ---- | :------------------------------------------------- |
| Delimiter | rune | Delimiter used for tokenising the compared string  |
| Extended  | bool | Enable extended globs, supporting classes, etc.    |
| GlobStar  | bool | Enables "double-star" match for one or more tokens |

# License

Please see [LICENSE](LICENSE) and [LICENSE-THIRD-PARTY](LICENSE-THIRD-PARTY.md).
