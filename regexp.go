//
// Copyright (c) 2019 Red Matter Ltd. UK
//
// The core of the logic is translated from https://github.com/fitzgen/glob-to-regexp
//

package globre

type options struct {
	Extended  bool
	GlobStar  bool
	Delimiter rune
}

// based on https://stackoverflow.com/a/400316/83005
func escapeDelimiter(delimiter rune) (outsideClass, insideClass string) {
	switch delimiter {
	case '^', '\\':
		return "\\" + string(delimiter), "\\" + string(delimiter)
	case '.', '$', '*', '+', '?', '(', ')', '[', '{', '|':
		return "\\" + string(delimiter), string(delimiter)
	case '-', ']':
		return string(delimiter), "\\" + string(delimiter)
	default:
		return string(delimiter), string(delimiter)
	}
}

// ExtendedSyntaxEnabled enable or disable extended glob syntax (default is disabled).
func ExtendedSyntaxEnabled(enabled bool) Option {
	return func(options *options) {
		options.Extended = enabled
	}
}

// GlobStarEnabled enable or disable globstar syntax (default is disabled).
func GlobStarEnabled(enabled bool) Option {
	return func(options *options) {
		options.GlobStar = enabled
	}
}

// WithDelimiter sets the delimiter to use (default is '/').
func WithDelimiter(delimiter rune) Option {
	return func(options *options) {
		options.Delimiter = delimiter
	}
}

type Option func(*options)

// RegexFromGlob converts the glob string into a regexp string with options "extended" & "glob-star" enabled, and with
// the default delimiter '/'
func RegexFromGlob(glob string, opts ...Option) string {
	o := options{
		Extended:  false,
		GlobStar:  false,
		Delimiter: '/',
	}

	// apply options
	for _, opt := range opts {
		opt(&o)
	}

	reStr := ""

	// escape the delimiter
	delimiterOutsideClass, delimiterInsideClass := escapeDelimiter(o.Delimiter)

	inGroup := false

	globLen := len(glob)

	for i := 0; i < globLen; i++ {
		c := glob[i]

		switch c {
		case '/', '$', '^', '+', '.', '(', ')', '=', '!', '|':
			reStr += "\\" + string(c)

		case '?':
			if o.Extended {
				reStr += "."
				break
			}

			reStr += "\\" + string(c)

		case '[', ']':
			if o.Extended {
				reStr += string(c)
				break
			}

			reStr += "\\" + string(c)

		case '{':
			if o.Extended {
				inGroup = true
				reStr += "("
				break
			}

			reStr += "\\" + string(c)

		case '}':
			if o.Extended {
				inGroup = false
				reStr += ")"
				break
			}

			reStr += "\\" + string(c)

		case ',':
			if inGroup {
				reStr += "|"
				break
			}

			reStr += "\\" + string(c)

		case '*':
			// Move over all consecutive "*"'s.
			// Also store the previous and next characters
			var nextChar, prevChar rune
			if i > 0 {
				prevChar = rune(glob[i-1])
			}
			starCount := 1
			for i < globLen-1 && glob[i+1] == '*' {
				starCount++
				i++
			}

			if i < globLen-1 {
				nextChar = rune(glob[i+1])
			}

			if !o.GlobStar {
				// globstar is disabled, so treat any number of "*" as one
				reStr += ".*"
			} else {
				// globstar is enabled, so determine if this is a globstar segment
				isGlobstar := starCount > 1 && // multiple "*"'s
					(prevChar == o.Delimiter || prevChar == 0) && // from the start of the segment
					(nextChar == o.Delimiter || nextChar == 0) // to the end of the segment

				if isGlobstar {
					// it's a globstar, so match zero or more path segments
					reStr += "(?:(?:[^" + delimiterInsideClass + "]*(?:" + delimiterOutsideClass + "|$))*)"
					i++ // move over the delimiter
				} else {
					// it's not a globstar, so only match one path segment
					reStr += "(?:[^" + delimiterInsideClass + "]*)"
				}
			}

		default:
			reStr += string(c)
		}
	}

	return "^" + reStr + "$"
}
