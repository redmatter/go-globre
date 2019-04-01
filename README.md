# Glob to RegExp

Glob string to regexp string conversion in golang 

# Usage

```go
urlGlob := "{http,https}://example.com/**"
// converts to ^(http|https):\/\/example\.com\/(?:(?:[^/]*(?:/|$))*)$
urlRe := globtoregexp.RegexFromGlobWithOptions(
	urlGlob, globtoregexp.GlobOptions{
		Extended:  true,
		GlobStar:  true,
		Delimiter: '/',
	})
urlRegexp := regexp.MustCompile(urlRe)
urlRegexp.MatchString("https://example.com/index.htm")
```

## Options

The conversion to regexp can be modified using `globtoregexp.GlobOptions`.

|   Attr    | Type |                        Desc                        |
| :-------- | ---- | :------------------------------------------------- |
| Delimiter | rune | Delimiter used for tonenisation                    |
| Extended  | bool | Enable extended globs, supporting classes, etc     |
| GlobStar  | bool | Enables "gouble-star" match for one or more tokens |

# License

Please see [LICENSE](LICENSE) and [LICENSE-THIRD-PARTY](LICENSE-THIRD-PARTY.md).
