//
// Copyright (c) 2019 Red Matter Ltd. UK
//

package globre

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertMatch(t *testing.T, glob, str string, opts ...Option) bool {
	reStr := RegexFromGlob(glob, opts...)
	t.Logf("regex from glob - glob:%s regex:%s", glob, reStr)
	re, err := regexp.Compile(reStr)
	if assert.Nil(t, err) && assert.NotNil(t, re) {
		return assert.True(t, re.MatchString(str))
	}

	return false
}

func assertNotMatch(t *testing.T, glob, str string, opts ...Option) bool {
	reStr := RegexFromGlob(glob, opts...)
	t.Logf("regex from glob - glob:%s regex:%s", glob, reStr)
	re, err := regexp.Compile(reStr)
	if assert.Nil(t, err) && assert.NotNil(t, re) {
		return assert.False(t, re.MatchString(str))
	}

	return false
}

// legacy from the JS library's perspective; left as it is for copy/paste easiness
func testRegexFromGlob_legacy(t *testing.T, GlobStar bool) {
	// Match everything
	assertMatch(t, "*", "foo")
	assertMatch(t, "*", "foo")

	// Match the end
	assertMatch(t, "f*", "foo")
	assertMatch(t, "f*", "foo")

	// Match the start
	assertMatch(t, "*o", "foo")
	assertMatch(t, "*o", "foo")

	// Match the middle
	assertMatch(t, "f*uck", "firetruck")

	// Don't match
	assertNotMatch(t, "uc", "firetruck")

	// Match zero characters
	assertMatch(t, "f*uck", "fuck")

	// More complex matches
	assertMatch(t, "*.min.js", "http://example.com/jquery.min.js", GlobStarEnabled(false))
	assertMatch(t, "*.min.*", "http://example.com/jquery.min.js", GlobStarEnabled(false))
	assertMatch(t, "*/js/*.js", "http://example.com/js/jquery.min.js", GlobStarEnabled(false))

	// More complex matches with RegExp 'g' flag (complex regression)
	assertMatch(t, "*.min.*", "http://example.com/jquery.min.js")
	assertMatch(t, "*.min.js", "http://example.com/jquery.min.js")
	assertMatch(t, "*/js/*.js", "http://example.com/js/jquery.min.js")

	// Test string  "\\\\/$^+?.()=!|{},[].*"  represents  <glob>\\/$^+?.()=!|{},[].*</glob>
	// The equivalent regex is:  /^\\\/\$\^\+\?\.\(\)\=\!\|\{\}\,\[\]\..*$/
	// Both glob and regex match:  \/$^+?.()=!|{},[].*
	testStr := "\\\\/$^+?.()=!|{},[].*"
	targetStr := "\\/$^+?.()=!|{},[].*"
	assertMatch(t, testStr, targetStr)

	// Equivalent matches without/with using RegExp 'g'
	assertNotMatch(t, ".min.", "http://example.com/jquery.min.js")
	assertMatch(t, "*.min.*", "http://example.com/jquery.min.js")

	assertNotMatch(t, "http:", "http://example.com/jquery.min.js")
	assertMatch(t, "http:*", "http://example.com/jquery.min.js")

	assertNotMatch(t, "min.js", "http://example.com/jquery.min.js")
	assertMatch(t, "*.min.js", "http://example.com/jquery.min.js")

	assertNotMatch(t, "/js*jq*.js", "http://example.com/js/jquery.min.js")

	// Extended mode

	// ?: Match one character, no more and no less
	assertMatch(t, "f?o", "foo", ExtendedSyntaxEnabled(true))
	assertNotMatch(t, "f?o", "fooo", ExtendedSyntaxEnabled(true))
	assertNotMatch(t, "f?oo", "foo", ExtendedSyntaxEnabled(true))

	// ?: Match one character with RegExp 'g'
	assertMatch(t, "f?o", "foo", ExtendedSyntaxEnabled(true), GlobStarEnabled(GlobStar))
	assertMatch(t, "f?o?", "fooo", ExtendedSyntaxEnabled(true), GlobStarEnabled(GlobStar))
	assertNotMatch(t, "?fo", "fooo", ExtendedSyntaxEnabled(true), GlobStarEnabled(GlobStar))
	assertNotMatch(t, "f?oo", "foo", ExtendedSyntaxEnabled(true), GlobStarEnabled(GlobStar))
	assertNotMatch(t, "foo?", "foo", ExtendedSyntaxEnabled(true), GlobStarEnabled(GlobStar))

	// []: Match a character range
	assertMatch(t, "fo[oz]", "foo", ExtendedSyntaxEnabled(true))
	assertMatch(t, "fo[oz]", "foz", ExtendedSyntaxEnabled(true))
	assertNotMatch(t, "fo[oz]", "fog", ExtendedSyntaxEnabled(true))

	// []: Match a character range and RegExp 'g' (regresion)
	assertMatch(t, "fo[oz]", "foo", ExtendedSyntaxEnabled(true), GlobStarEnabled(GlobStar))
	assertMatch(t, "fo[oz]", "foz", ExtendedSyntaxEnabled(true), GlobStarEnabled(GlobStar))
	assertNotMatch(t, "fo[oz]", "fog", ExtendedSyntaxEnabled(true), GlobStarEnabled(GlobStar))

	// Options{}: Match a choice of different substrings
	assertMatch(t, "foo{bar,baaz}", "foobaaz", ExtendedSyntaxEnabled(true))
	assertMatch(t, "foo{bar,baaz}", "foobar", ExtendedSyntaxEnabled(true))
	assertNotMatch(t, "foo{bar,baaz}", "foobuzz", ExtendedSyntaxEnabled(true))
	assertMatch(t, "foo{bar,b*z}", "foobuzz", ExtendedSyntaxEnabled(true))

	// Options{}: Match a choice of different substrings and RegExp 'g' (regression)
	assertMatch(t, "foo{bar,baaz}", "foobaaz", ExtendedSyntaxEnabled(true), GlobStarEnabled(GlobStar))
	assertMatch(t, "foo{bar,baaz}", "foobar", ExtendedSyntaxEnabled(true), GlobStarEnabled(GlobStar))
	assertNotMatch(t, "foo{bar,baaz}", "foobuzz", ExtendedSyntaxEnabled(true), GlobStarEnabled(GlobStar))
	assertMatch(t, "foo{bar,b*z}", "foobuzz", ExtendedSyntaxEnabled(true), GlobStarEnabled(GlobStar))

	// More complex extended matches
	assertMatch(t, "http://?o[oz].b*z.com/{*.js,*.html}",
		"http://foo.baaz.com/jquery.min.js",
		ExtendedSyntaxEnabled(true))
	assertMatch(t, "http://?o[oz].b*z.com/{*.js,*.html}",
		"http://moz.buzz.com/index.html",
		ExtendedSyntaxEnabled(true))
	assertNotMatch(t, "http://?o[oz].b*z.com/{*.js,*.html}",
		"http://moz.buzz.com/index.htm",
		ExtendedSyntaxEnabled(true))
	assertNotMatch(t, "http://?o[oz].b*z.com/{*.js,*.html}",
		"http://moz.bar.com/index.html",
		ExtendedSyntaxEnabled(true))
	assertNotMatch(t, "http://?o[oz].b*z.com/{*.js,*.html}",
		"http://flozz.buzz.com/index.html",
		ExtendedSyntaxEnabled(true))

	// More complex extended matches and RegExp 'g' (regresion)
	assertMatch(t, "http://?o[oz].b*z.com/{*.js,*.html}",
		"http://foo.baaz.com/jquery.min.js",
		ExtendedSyntaxEnabled(true), GlobStarEnabled(GlobStar))
	assertMatch(t, "http://?o[oz].b*z.com/{*.js,*.html}",
		"http://moz.buzz.com/index.html",
		ExtendedSyntaxEnabled(true), GlobStarEnabled(GlobStar))
	assertNotMatch(t, "http://?o[oz].b*z.com/{*.js,*.html}",
		"http://moz.buzz.com/index.htm",
		ExtendedSyntaxEnabled(true), GlobStarEnabled(GlobStar))
	assertNotMatch(t, "http://?o[oz].b*z.com/{*.js,*.html}",
		"http://moz.bar.com/index.html",
		ExtendedSyntaxEnabled(true), GlobStarEnabled(GlobStar))
	assertNotMatch(t, "http://?o[oz].b*z.com/{*.js,*.html}",
		"http://flozz.buzz.com/index.html",
		ExtendedSyntaxEnabled(true), GlobStarEnabled(GlobStar))

	// GlobStar
	assertMatch(t, "http://foo.com/**/{*.js,*.html}",
		"http://foo.com/bar/jquery.min.js",
		ExtendedSyntaxEnabled(true), GlobStarEnabled(GlobStar))
	assertMatch(t, "http://foo.com/**/{*.js,*.html}",
		"http://foo.com/bar/baz/jquery.min.js",
		ExtendedSyntaxEnabled(true), GlobStarEnabled(GlobStar))
	assertMatch(t, "http://foo.com/**",
		"http://foo.com/bar/baz/jquery.min.js",
		ExtendedSyntaxEnabled(true), GlobStarEnabled(GlobStar))

	// Remaining special chars should still match themselves
	// Test string  "\\\\/$^+.()=!|,.*"  represents  <glob>\\/$^+.()=!|,.*</glob>
	// The equivalent regex is:  /^\\\/\$\^\+\.\(\)\=\!\|\,\..*$/
	// Both glob and regex match:  \/$^+.()=!|,.*
	var testExtStr = "\\\\/$^+.()=!|,.*"
	var targetExtStr = "\\/$^+.()=!|,.*"
	assertMatch(t, testExtStr, targetExtStr, ExtendedSyntaxEnabled(true))
	assertMatch(t, testExtStr, targetExtStr, ExtendedSyntaxEnabled(true), GlobStarEnabled(GlobStar))
}

func TestRegexFromGlob(t *testing.T) {
	// regression
	// GlobStar false
	testRegexFromGlob_legacy(t, false)
	// GlobStar true
	testRegexFromGlob_legacy(t, true)

	// GlobStar specific tests
	assertMatch(t, "/foo/*", "/foo/bar.txt", GlobStarEnabled(true))
	assertMatch(t, "/foo/**", "/foo/baz.txt", GlobStarEnabled(true))
	assertMatch(t, "/foo/**", "/foo/bar/baz.txt", GlobStarEnabled(true))
	assertMatch(t, "/foo/*/*.txt", "/foo/bar/baz.txt", GlobStarEnabled(true))
	assertMatch(t, "/foo/**/*.txt", "/foo/bar/baz.txt", GlobStarEnabled(true))
	assertMatch(t, "/foo/**/*.txt", "/foo/bar/baz/qux.txt", GlobStarEnabled(true))
	assertMatch(t, "/foo/**/bar.txt", "/foo/bar.txt", GlobStarEnabled(true))
	assertMatch(t, "/foo/**/**/bar.txt", "/foo/bar.txt", GlobStarEnabled(true))
	assertMatch(t, "/foo/**/*/baz.txt", "/foo/bar/baz.txt", GlobStarEnabled(true))
	assertMatch(t, "/foo/**/*.txt", "/foo/bar.txt", GlobStarEnabled(true))
	assertMatch(t, "/foo/**/**/*.txt", "/foo/bar.txt", GlobStarEnabled(true))
	assertMatch(t, "/foo/**/*/*.txt", "/foo/bar/baz.txt", GlobStarEnabled(true))
	assertMatch(t, "**/*.txt", "/foo/bar/baz/qux.txt", GlobStarEnabled(true))
	assertMatch(t, "**/foo.txt", "foo.txt", GlobStarEnabled(true))
	assertMatch(t, "**/*.txt", "foo.txt", GlobStarEnabled(true))

	assertNotMatch(t, "/foo/*", "/foo/bar/baz.txt", GlobStarEnabled(true))
	assertNotMatch(t, "/foo/*.txt", "/foo/bar/baz.txt", GlobStarEnabled(true))
	assertNotMatch(t, "/foo/*/*.txt", "/foo/bar/baz/qux.txt", GlobStarEnabled(true))
	assertNotMatch(t, "/foo/*/bar.txt", "/foo/bar.txt", GlobStarEnabled(true))
	assertNotMatch(t, "/foo/*/*/baz.txt", "/foo/bar/baz.txt", GlobStarEnabled(true))
	assertNotMatch(t, "/foo/**.txt", "/foo/bar/baz/qux.txt", GlobStarEnabled(true))
	assertNotMatch(t, "/foo/bar**/*.txt", "/foo/bar/baz/qux.txt", GlobStarEnabled(true))
	assertNotMatch(t, "/foo/bar**", "/foo/bar/baz.txt", GlobStarEnabled(true))
	assertNotMatch(t, "**/.txt", "/foo/bar/baz/qux.txt", GlobStarEnabled(true))
	assertNotMatch(t, "*/*.txt", "/foo/bar/baz/qux.txt", GlobStarEnabled(true))
	assertNotMatch(t, "*/*.txt", "foo.txt", GlobStarEnabled(true))

	assertNotMatch(t, "http://foo.com/*",
		"http://foo.com/bar/baz/jquery.min.js",
		ExtendedSyntaxEnabled(true), GlobStarEnabled(true))
	assertNotMatch(t, "http://foo.com/*",
		"http://foo.com/bar/baz/jquery.min.js",
		GlobStarEnabled(true))

	assertMatch(t, "http://foo.com/*",
		"http://foo.com/bar/baz/jquery.min.js",
		GlobStarEnabled(false))
	assertMatch(t, "http://foo.com/**",
		"http://foo.com/bar/baz/jquery.min.js",
		GlobStarEnabled(true))

	assertMatch(t, "http://foo.com/*/*/jquery.min.js",
		"http://foo.com/bar/baz/jquery.min.js",
		GlobStarEnabled(true))
	assertMatch(t, "http://foo.com/**/jquery.min.js",
		"http://foo.com/bar/baz/jquery.min.js",
		GlobStarEnabled(true))
	assertMatch(t, "http://foo.com/*/*/jquery.min.js",
		"http://foo.com/bar/baz/jquery.min.js",
		GlobStarEnabled(false))
	assertMatch(t, "http://foo.com/*/jquery.min.js",
		"http://foo.com/bar/baz/jquery.min.js",
		GlobStarEnabled(false))
	assertNotMatch(t, "http://foo.com/*/jquery.min.js",
		"http://foo.com/bar/baz/jquery.min.js",
		GlobStarEnabled(true))
}

func TestRegexFromGlobWithDelimiter(t *testing.T) {
	assertMatch(t, "http*?**",
		"http://foo.com/bar/baz/jquery.min.js?yahhooo",
		GlobStarEnabled(true), WithDelimiter('?'))
}
