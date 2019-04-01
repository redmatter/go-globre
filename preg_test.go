//
// Copyright (c) 2019 Red Matter Ltd. UK. All rights reserved.
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
//
// Created by Dino Korah on 01/04/2019.
//

package globtoregexp

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertMatch(t *testing.T, glob, str string, opts GlobOptions) bool {
	reStr := RegexFromGlobWithOptions(glob, opts)
	t.Logf("regex from glob - glob:%s regex:%s", glob, reStr)
	re, err := regexp.Compile(reStr)
	if assert.Nil(t, err) && assert.NotNil(t, re) {
		return assert.True(t, re.MatchString(str))
	}

	return false
}

func assertNotMatch(t *testing.T, glob, str string, opts GlobOptions) bool {
	reStr := RegexFromGlobWithOptions(glob, opts)
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
	assertMatch(t, "*", "foo", GlobOptions{})
	assertMatch(t, "*", "foo", GlobOptions{})

	// Match the end
	assertMatch(t, "f*", "foo", GlobOptions{})
	assertMatch(t, "f*", "foo", GlobOptions{})

	// Match the start
	assertMatch(t, "*o", "foo", GlobOptions{})
	assertMatch(t, "*o", "foo", GlobOptions{})

	// Match the middle
	assertMatch(t, "f*uck", "firetruck", GlobOptions{})

	// Don't match
	assertNotMatch(t, "uc", "firetruck", GlobOptions{})

	// Match zero characters
	assertMatch(t, "f*uck", "fuck", GlobOptions{})

	// More complex matches
	assertMatch(t, "*.min.js", "http://example.com/jquery.min.js", GlobOptions{GlobStar: false})
	assertMatch(t, "*.min.*", "http://example.com/jquery.min.js", GlobOptions{GlobStar: false})
	assertMatch(t, "*/js/*.js", "http://example.com/js/jquery.min.js", GlobOptions{GlobStar: false})

	// More complex matches with RegExp 'g' flag (complex regression)
	assertMatch(t, "*.min.*", "http://example.com/jquery.min.js", GlobOptions{})
	assertMatch(t, "*.min.js", "http://example.com/jquery.min.js", GlobOptions{})
	assertMatch(t, "*/js/*.js", "http://example.com/js/jquery.min.js", GlobOptions{})

	// Test string  "\\\\/$^+?.()=!|{},[].*"  represents  <glob>\\/$^+?.()=!|{},[].*</glob>
	// The equivalent regex is:  /^\\\/\$\^\+\?\.\(\)\=\!\|\{\}\,\[\]\..*$/
	// Both glob and regex match:  \/$^+?.()=!|{},[].*
	testStr := "\\\\/$^+?.()=!|{},[].*"
	targetStr := "\\/$^+?.()=!|{},[].*"
	assertMatch(t, testStr, targetStr, GlobOptions{})

	// Equivalent matches without/with using RegExp 'g'
	assertNotMatch(t, ".min.", "http://example.com/jquery.min.js", GlobOptions{})
	assertMatch(t, "*.min.*", "http://example.com/jquery.min.js", GlobOptions{})

	assertNotMatch(t, "http:", "http://example.com/jquery.min.js", GlobOptions{})
	assertMatch(t, "http:*", "http://example.com/jquery.min.js", GlobOptions{})

	assertNotMatch(t, "min.js", "http://example.com/jquery.min.js", GlobOptions{})
	assertMatch(t, "*.min.js", "http://example.com/jquery.min.js", GlobOptions{})

	assertNotMatch(t, "/js*jq*.js", "http://example.com/js/jquery.min.js", GlobOptions{})

	// Extended mode

	// ?: Match one character, no more and no less
	assertMatch(t, "f?o", "foo", GlobOptions{Extended: true})
	assertNotMatch(t, "f?o", "fooo", GlobOptions{Extended: true})
	assertNotMatch(t, "f?oo", "foo", GlobOptions{Extended: true})

	// ?: Match one character with RegExp 'g'
	assertMatch(t, "f?o", "foo", GlobOptions{Extended: true, GlobStar: GlobStar})
	assertMatch(t, "f?o?", "fooo", GlobOptions{Extended: true, GlobStar: GlobStar})
	assertNotMatch(t, "?fo", "fooo", GlobOptions{Extended: true, GlobStar: GlobStar})
	assertNotMatch(t, "f?oo", "foo", GlobOptions{Extended: true, GlobStar: GlobStar})
	assertNotMatch(t, "foo?", "foo", GlobOptions{Extended: true, GlobStar: GlobStar})

	// []: Match a character range
	assertMatch(t, "fo[oz]", "foo", GlobOptions{Extended: true})
	assertMatch(t, "fo[oz]", "foz", GlobOptions{Extended: true})
	assertNotMatch(t, "fo[oz]", "fog", GlobOptions{Extended: true})

	// []: Match a character range and RegExp 'g' (regresion)
	assertMatch(t, "fo[oz]", "foo", GlobOptions{Extended: true, GlobStar: GlobStar})
	assertMatch(t, "fo[oz]", "foz", GlobOptions{Extended: true, GlobStar: GlobStar})
	assertNotMatch(t, "fo[oz]", "fog", GlobOptions{Extended: true, GlobStar: GlobStar})

	// GlobOptions{}: Match a choice of different substrings
	assertMatch(t, "foo{bar,baaz}", "foobaaz", GlobOptions{Extended: true})
	assertMatch(t, "foo{bar,baaz}", "foobar", GlobOptions{Extended: true})
	assertNotMatch(t, "foo{bar,baaz}", "foobuzz", GlobOptions{Extended: true})
	assertMatch(t, "foo{bar,b*z}", "foobuzz", GlobOptions{Extended: true})

	// GlobOptions{}: Match a choice of different substrings and RegExp 'g' (regression)
	assertMatch(t, "foo{bar,baaz}", "foobaaz", GlobOptions{Extended: true, GlobStar: GlobStar})
	assertMatch(t, "foo{bar,baaz}", "foobar", GlobOptions{Extended: true, GlobStar: GlobStar})
	assertNotMatch(t, "foo{bar,baaz}", "foobuzz", GlobOptions{Extended: true, GlobStar: GlobStar})
	assertMatch(t, "foo{bar,b*z}", "foobuzz", GlobOptions{Extended: true, GlobStar: GlobStar})

	// More complex extended matches
	assertMatch(t, "http://?o[oz].b*z.com/{*.js,*.html}",
		"http://foo.baaz.com/jquery.min.js",
		GlobOptions{Extended: true})
	assertMatch(t, "http://?o[oz].b*z.com/{*.js,*.html}",
		"http://moz.buzz.com/index.html",
		GlobOptions{Extended: true})
	assertNotMatch(t, "http://?o[oz].b*z.com/{*.js,*.html}",
		"http://moz.buzz.com/index.htm",
		GlobOptions{Extended: true})
	assertNotMatch(t, "http://?o[oz].b*z.com/{*.js,*.html}",
		"http://moz.bar.com/index.html",
		GlobOptions{Extended: true})
	assertNotMatch(t, "http://?o[oz].b*z.com/{*.js,*.html}",
		"http://flozz.buzz.com/index.html",
		GlobOptions{Extended: true})

	// More complex extended matches and RegExp 'g' (regresion)
	assertMatch(t, "http://?o[oz].b*z.com/{*.js,*.html}",
		"http://foo.baaz.com/jquery.min.js",
		GlobOptions{Extended: true, GlobStar: GlobStar})
	assertMatch(t, "http://?o[oz].b*z.com/{*.js,*.html}",
		"http://moz.buzz.com/index.html",
		GlobOptions{Extended: true, GlobStar: GlobStar})
	assertNotMatch(t, "http://?o[oz].b*z.com/{*.js,*.html}",
		"http://moz.buzz.com/index.htm",
		GlobOptions{Extended: true, GlobStar: GlobStar})
	assertNotMatch(t, "http://?o[oz].b*z.com/{*.js,*.html}",
		"http://moz.bar.com/index.html",
		GlobOptions{Extended: true, GlobStar: GlobStar})
	assertNotMatch(t, "http://?o[oz].b*z.com/{*.js,*.html}",
		"http://flozz.buzz.com/index.html",
		GlobOptions{Extended: true, GlobStar: GlobStar})

	// GlobStar
	assertMatch(t, "http://foo.com/**/{*.js,*.html}",
		"http://foo.com/bar/jquery.min.js",
		GlobOptions{Extended: true, GlobStar: GlobStar})
	assertMatch(t, "http://foo.com/**/{*.js,*.html}",
		"http://foo.com/bar/baz/jquery.min.js",
		GlobOptions{Extended: true, GlobStar: GlobStar})
	assertMatch(t, "http://foo.com/**",
		"http://foo.com/bar/baz/jquery.min.js",
		GlobOptions{Extended: true, GlobStar: GlobStar})

	// Remaining special chars should still match themselves
	// Test string  "\\\\/$^+.()=!|,.*"  represents  <glob>\\/$^+.()=!|,.*</glob>
	// The equivalent regex is:  /^\\\/\$\^\+\.\(\)\=\!\|\,\..*$/
	// Both glob and regex match:  \/$^+.()=!|,.*
	var testExtStr = "\\\\/$^+.()=!|,.*"
	var targetExtStr = "\\/$^+.()=!|,.*"
	assertMatch(t, testExtStr, targetExtStr, GlobOptions{Extended: true})
	assertMatch(t, testExtStr, targetExtStr, GlobOptions{Extended: true, GlobStar: GlobStar})
}

func TestRegexFromGlob(t *testing.T) {
	// regression
	// GlobStar false
	testRegexFromGlob_legacy(t, false)
	// GlobStar true
	testRegexFromGlob_legacy(t, true)

	// GlobStar specific tests
	assertMatch(t, "/foo/*", "/foo/bar.txt", GlobOptions{GlobStar: true})
	assertMatch(t, "/foo/**", "/foo/baz.txt", GlobOptions{GlobStar: true})
	assertMatch(t, "/foo/**", "/foo/bar/baz.txt", GlobOptions{GlobStar: true})
	assertMatch(t, "/foo/*/*.txt", "/foo/bar/baz.txt", GlobOptions{GlobStar: true})
	assertMatch(t, "/foo/**/*.txt", "/foo/bar/baz.txt", GlobOptions{GlobStar: true})
	assertMatch(t, "/foo/**/*.txt", "/foo/bar/baz/qux.txt", GlobOptions{GlobStar: true})
	assertMatch(t, "/foo/**/bar.txt", "/foo/bar.txt", GlobOptions{GlobStar: true})
	assertMatch(t, "/foo/**/**/bar.txt", "/foo/bar.txt", GlobOptions{GlobStar: true})
	assertMatch(t, "/foo/**/*/baz.txt", "/foo/bar/baz.txt", GlobOptions{GlobStar: true})
	assertMatch(t, "/foo/**/*.txt", "/foo/bar.txt", GlobOptions{GlobStar: true})
	assertMatch(t, "/foo/**/**/*.txt", "/foo/bar.txt", GlobOptions{GlobStar: true})
	assertMatch(t, "/foo/**/*/*.txt", "/foo/bar/baz.txt", GlobOptions{GlobStar: true})
	assertMatch(t, "**/*.txt", "/foo/bar/baz/qux.txt", GlobOptions{GlobStar: true})
	assertMatch(t, "**/foo.txt", "foo.txt", GlobOptions{GlobStar: true})
	assertMatch(t, "**/*.txt", "foo.txt", GlobOptions{GlobStar: true})

	assertNotMatch(t, "/foo/*", "/foo/bar/baz.txt", GlobOptions{GlobStar: true})
	assertNotMatch(t, "/foo/*.txt", "/foo/bar/baz.txt", GlobOptions{GlobStar: true})
	assertNotMatch(t, "/foo/*/*.txt", "/foo/bar/baz/qux.txt", GlobOptions{GlobStar: true})
	assertNotMatch(t, "/foo/*/bar.txt", "/foo/bar.txt", GlobOptions{GlobStar: true})
	assertNotMatch(t, "/foo/*/*/baz.txt", "/foo/bar/baz.txt", GlobOptions{GlobStar: true})
	assertNotMatch(t, "/foo/**.txt", "/foo/bar/baz/qux.txt", GlobOptions{GlobStar: true})
	assertNotMatch(t, "/foo/bar**/*.txt", "/foo/bar/baz/qux.txt", GlobOptions{GlobStar: true})
	assertNotMatch(t, "/foo/bar**", "/foo/bar/baz.txt", GlobOptions{GlobStar: true})
	assertNotMatch(t, "**/.txt", "/foo/bar/baz/qux.txt", GlobOptions{GlobStar: true})
	assertNotMatch(t, "*/*.txt", "/foo/bar/baz/qux.txt", GlobOptions{GlobStar: true})
	assertNotMatch(t, "*/*.txt", "foo.txt", GlobOptions{GlobStar: true})

	assertNotMatch(t, "http://foo.com/*",
		"http://foo.com/bar/baz/jquery.min.js",
		GlobOptions{Extended: true, GlobStar: true})
	assertNotMatch(t, "http://foo.com/*",
		"http://foo.com/bar/baz/jquery.min.js",
		GlobOptions{GlobStar: true})

	assertMatch(t, "http://foo.com/*",
		"http://foo.com/bar/baz/jquery.min.js",
		GlobOptions{GlobStar: false})
	assertMatch(t, "http://foo.com/**",
		"http://foo.com/bar/baz/jquery.min.js",
		GlobOptions{GlobStar: true})

	assertMatch(t, "http://foo.com/*/*/jquery.min.js",
		"http://foo.com/bar/baz/jquery.min.js",
		GlobOptions{GlobStar: true})
	assertMatch(t, "http://foo.com/**/jquery.min.js",
		"http://foo.com/bar/baz/jquery.min.js",
		GlobOptions{GlobStar: true})
	assertMatch(t, "http://foo.com/*/*/jquery.min.js",
		"http://foo.com/bar/baz/jquery.min.js",
		GlobOptions{GlobStar: false})
	assertMatch(t, "http://foo.com/*/jquery.min.js",
		"http://foo.com/bar/baz/jquery.min.js",
		GlobOptions{GlobStar: false})
	assertNotMatch(t, "http://foo.com/*/jquery.min.js",
		"http://foo.com/bar/baz/jquery.min.js",
		GlobOptions{GlobStar: true})
}

func TestRegexFromGlobWithDelimiter(t *testing.T) {
	assertMatch(t, "http*?**",
		"http://foo.com/bar/baz/jquery.min.js?yahhooo",
		GlobOptions{GlobStar: true, Delimiter: '?'})
}
