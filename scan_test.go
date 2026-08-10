package picomatch

import (
	"reflect"
	"testing"
)

func TestScan_Both(t *testing.T) {
	tests := []struct {
		pattern string
		base    string
		glob    string
	}{
		{"foo/bar", "foo/bar", ""},
		{"foo/@bar", "foo/@bar", ""},
		{"foo/@bar\\+", "foo/@bar\\+", ""},
		{"foo/bar+", "foo/bar+", ""},
		{"foo/bar*", "foo", "bar*"},

		{"!foo", "foo", ""},
		{"*", "", "*"},
		{"**", "", "**"},
		{"**/*.md", "", "**/*.md"},
		{"**/*.min.js", "", "**/*.min.js"},
		{"**/*foo.js", "", "**/*foo.js"},
		{"**/.*", "", "**/.*"},
		{"**/d", "", "**/d"},
		{"*.*", "", "*.*"},
		{"*.js", "", "*.js"},
		{"*.md", "", "*.md"},
		{"*.min.js", "", "*.min.js"},
		{"*/*", "", "*/*"},
		{"*/*/*/*", "", "*/*/*/*"},
		{"*/*/*/e", "", "*/*/*/e"},
		{"*/b/*/e", "", "*/b/*/e"},
		{"*b", "", "*b"},
		{".*", "", ".*"},
		{"*", "", "*"},
		{"a/**/j/**/z/*.md", "a", "**/j/**/z/*.md"},
		{"a/**/z/*.md", "a", "**/z/*.md"},
		{"node_modules/*-glob/**/*.js", "node_modules", "*-glob/**/*.js"},
		{"{a/b/{c,/foo.js}/e.f.g}", "", "{a/b/{c,/foo.js}/e.f.g}"},
		{".a*", "", ".a*"},
		{".b*", "", ".b*"},
		{"/*", "/", "*"},
		{"a/***", "a", "***"},
		{"a/**/b/*.{foo,bar}", "a", "**/b/*.{foo,bar}"},
		{"a/**/c/*", "a", "**/c/*"},
		{"a/**/c/*.md", "a", "**/c/*.md"},
		{"a/**/e", "a", "**/e"},
		{"a/**/j/**/z/*.md", "a", "**/j/**/z/*.md"},
		{"a/**/z/*.md", "a", "**/z/*.md"},
		{"a/**c*", "a", "**c*"},
		{"a/**c/*", "a", "**c/*"},
		{"a/*/*/e", "a", "*/*/e"},
		{"a/*/c/*.md", "a", "*/c/*.md"},
		{"a/b/**/c{d,e}/**/xyz.md", "a/b", "**/c{d,e}/**/xyz.md"},
		{"a/b/**/e", "a/b", "**/e"},
		{"a/b/*.{foo,bar}", "a/b", "*.{foo,bar}"},
		{"a/b/*/e", "a/b", "*/e"},
		{"a/b/.git/", "a/b/.git/", ""},
		{"a/b/.git/**", "a/b/.git", "**"},
		{"a/b/.{foo,bar}", "a/b", ".{foo,bar}"},
		{"a/b/c/*", "a/b/c", "*"},
		{"a/b/c/**/*.min.js", "a/b/c", "**/*.min.js"},
		{"a/b/c/*.md", "a/b/c", "*.md"},
		{"a/b/c/.*.md", "a/b/c", ".*.md"},
		{"a/b/{c,.gitignore,{a,b}}/{a,b}/abc.foo.js", "a/b", "{c,.gitignore,{a,b}}/{a,b}/abc.foo.js"},
		{"a/b/{c,/.gitignore}", "a/b", "{c,/.gitignore}"},
		{"a/b/{c,d}/", "a/b", "{c,d}/"},
		{"a/b/{c,d}/e/f.g", "a/b", "{c,d}/e/f.g"},
		{"b/*/*/*", "b", "*/*/*"},
		{".md", ".md", ""},

		{"!*.min.js", "", "*.min.js"},
		{"!foo", "foo", ""},
		{"!foo/*.js", "foo", "*.js"},
		{"!foo/(a|b).min.js", "foo", "(a|b).min.js"},
		{"!foo/[a-b].min.js", "foo", "[a-b].min.js"},
		{"!foo/{a,b}.min.js", "foo", "{a,b}.min.js"},
		{"a/b/c/!foo", "a/b/c/!foo", ""},

		{"/a/b/!(a|b)/e.f.g/", "/a/b", "!(a|b)/e.f.g/"},
		{"/a/b/@(a|b)/e.f.g/", "/a/b", "@(a|b)/e.f.g/"},
		{"@(a|b)/e.f.g/", "", "@(a|b)/e.f.g/"},

		{"[a-c]b*", "", "[a-c]b*"},
		{"[a-j]*[^c]", "", "[a-j]*[^c]"},
		{"[a-j]*[^c]b/c", "", "[a-j]*[^c]b/c"},
		{"[a-j]*[^c]bc", "", "[a-j]*[^c]bc"},
		{"[ab][ab]", "", "[ab][ab]"},
		{"foo/[a-b].min.js", "foo", "[a-b].min.js"},

		{"?", "", "?"},
		{"?/?", "", "?/?"},
		{"??", "", "??"},
		{"???", "", "???"},
		{"?a", "", "?a"},
		{"?b", "", "?b"},
		{"a?b", "", "a?b"},
		{"a/?/c.js", "a", "?/c.js"},
		{"a/?/c.md", "a", "?/c.md"},
		{"a/?/c/?/*/f.js", "a", "?/c/?/*/f.js"},
		{"a/?/c/?/*/f.md", "a", "?/c/?/*/f.md"},
		{"a/?/c/?/e.js", "a", "?/c/?/e.js"},
		{"a/?/c/?/e.md", "a", "?/c/?/e.md"},
		{"a/?/c/???/e.js", "a", "?/c/???/e.js"},
		{"a/?/c/???/e.md", "a", "?/c/???/e.md"},
		{"a/??/c.js", "a", "??/c.js"},
		{"a/??/c.md", "a", "??/c.md"},
		{"a/???/c.js", "a", "???/c.js"},
		{"a/???/c.md", "a", "???/c.md"},
		{"a/????/c.js", "a", "????/c.js"},

		{"", "", ""},
		{".", ".", ""},
		{"a", "a", ""},
		{".a", ".a", ""},
		{"/a", "/a", ""},
		{"a/", "a/", ""},
		{"/a/", "/a/", ""},
		{"/a/b/c", "/a/b/c", ""},
		{"/a/b/c/", "/a/b/c/", ""},
		{"a/b/c/", "a/b/c/", ""},
		{"a.min.js", "a.min.js", ""},
		{"a/.x.md", "a/.x.md", ""},
		{"a/b/.gitignore", "a/b/.gitignore", ""},
		{"a/b/c/d.md", "a/b/c/d.md", ""},
		{"a/b/c/d.e.f/g.min.js", "a/b/c/d.e.f/g.min.js", ""},
		{"a/b/.git", "a/b/.git", ""},
		{"a/b/.git/", "a/b/.git/", ""},
		{"a/b/c", "a/b/c", ""},
		{"a/b/c.d/e.md", "a/b/c.d/e.md", ""},
		{"a/b/c.md", "a/b/c.md", ""},
		{"a/b/c.min.js", "a/b/c.min.js", ""},
		{"a/b/git/", "a/b/git/", ""},
		{"aa", "aa", ""},
		{"ab", "ab", ""},
		{"bb", "bb", ""},
		{"c.md", "c.md", ""},
		{"foo", "foo", ""},

		{"/a/b/{c,/foo.js}/e.f.g/", "/a/b", "{c,/foo.js}/e.f.g/"},
		{"{a/b/c.js,/a/b/{c,/foo.js}/e.f.g/}", "", "{a/b/c.js,/a/b/{c,/foo.js}/e.f.g/}"},
		{"/a/b/{c,d}/", "/a/b", "{c,d}/"},
		{"/a/b/{c,d}/*.js", "/a/b", "{c,d}/*.js"},
		{"/a/b/{c,d}/*.min.js", "/a/b", "{c,d}/*.min.js"},
		{"/a/b/{c,d}/e.f.g/", "/a/b", "{c,d}/e.f.g/"},
		{"{.,*}", "", "{.,*}"},
		{"foo/{a,b}.min.js", "foo", "{a,b}.min.js"},
	}

	for _, tc := range tests {
		t.Run(tc.pattern, func(t *testing.T) {
			state := Scan(tc.pattern, nil)
			if state.Base != tc.base {
				t.Errorf("Expected base %q, got %q", tc.base, state.Base)
			}
			if state.Glob != tc.glob {
				t.Errorf("Expected glob %q, got %q", tc.glob, state.Glob)
			}
		})
	}
}

func TestScan_Scan(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		opts     *Options
		expected ScanState
	}{
		{
			name:    "should handle leading './'",
			pattern: "./foo/bar/*.js",
			expected: ScanState{
				Input: "./foo/bar/*.js", Prefix: "./", Start: 2, Base: "foo/bar", Glob: "*.js",
				IsBrace: false, IsBracket: false, IsGlob: true, IsGlobstar: false,
				IsExtglob: false, Negated: false, NegatedExtglob: false,
			},
		},
		{
			name:    "should detect braces",
			pattern: "foo/{a,b,c}/*.js",
			opts:    &Options{ScanToEnd: true},
			expected: ScanState{
				Input: "foo/{a,b,c}/*.js", Prefix: "", Start: 0, Base: "foo", Glob: "{a,b,c}/*.js",
				IsBrace: true, IsBracket: false, IsGlob: true, IsGlobstar: false,
				IsExtglob: false, Negated: false, NegatedExtglob: false,
			},
		},
		{
			name:    "should detect globstars",
			pattern: "./foo/**/*.js",
			opts:    &Options{ScanToEnd: true},
			expected: ScanState{
				Input: "./foo/**/*.js", Prefix: "./", Start: 2, Base: "foo", Glob: "**/*.js",
				IsBrace: false, IsBracket: false, IsGlob: true, IsGlobstar: true,
				IsExtglob: false, Negated: false, NegatedExtglob: false,
			},
		},
		{
			name:    "should detect extglobs",
			pattern: "./foo/@(foo)/*.js",
			expected: ScanState{
				Input: "./foo/@(foo)/*.js", Prefix: "./", Start: 2, Base: "foo", Glob: "@(foo)/*.js",
				IsBrace: false, IsBracket: false, IsGlob: true, IsGlobstar: false,
				IsExtglob: true, Negated: false, NegatedExtglob: false,
			},
		},
		{
			name:    "should detect extglobs and globstars",
			pattern: "./foo/@(bar)/**/*.js",
			opts:    &Options{Parts: true},
			expected: ScanState{
				Input: "./foo/@(bar)/**/*.js", Prefix: "./", Start: 2, Base: "foo", Glob: "@(bar)/**/*.js",
				IsBrace: false, IsBracket: false, IsGlob: true, IsGlobstar: true,
				IsExtglob: true, Negated: false, NegatedExtglob: false,
				Slashes: []int{1, 5, 12, 15}, Parts: []string{"foo", "@(bar)", "**", "*.js"},
			},
		},
		{
			name:    "should handle leading '!'",
			pattern: "!foo/bar/*.js",
			expected: ScanState{
				Input: "!foo/bar/*.js", Prefix: "!", Start: 1, Base: "foo/bar", Glob: "*.js",
				IsBrace: false, IsBracket: false, IsGlob: true, IsGlobstar: false,
				IsExtglob: false, Negated: true, NegatedExtglob: false,
			},
		},
		{
			name:    "should detect negated extglobs at the beginning 1",
			pattern: "!(foo)*",
			expected: ScanState{
				Input: "!(foo)*", Prefix: "", Start: 0, Base: "", Glob: "!(foo)*",
				IsBrace: false, IsBracket: false, IsGlob: true, IsGlobstar: false,
				IsExtglob: true, Negated: false, NegatedExtglob: true,
			},
		},
		{
			name:    "should detect negated extglobs at the beginning 2",
			pattern: "!(foo)",
			expected: ScanState{
				Input: "!(foo)", Prefix: "", Start: 0, Base: "", Glob: "!(foo)",
				IsBrace: false, IsBracket: false, IsGlob: true, IsGlobstar: false,
				IsExtglob: true, Negated: false, NegatedExtglob: true,
			},
		},
		{
			name:    "should not detect negated extglobs in the middle",
			pattern: "test/!(foo)/*",
			expected: ScanState{
				Input: "test/!(foo)/*", Prefix: "", Start: 0, Base: "test", Glob: "!(foo)/*",
				IsBrace: false, IsBracket: false, IsGlob: true, IsGlobstar: false,
				IsExtglob: true, Negated: false, NegatedExtglob: false,
			},
		},
		{
			name:    "should handle leading './' when negated 1",
			pattern: "./!foo/bar/*.js",
			expected: ScanState{
				Input: "./!foo/bar/*.js", Prefix: "./!", Start: 3, Base: "foo/bar", Glob: "*.js",
				IsBrace: false, IsBracket: false, IsGlob: true, IsGlobstar: false,
				IsExtglob: false, Negated: true, NegatedExtglob: false,
			},
		},
		{
			name:    "should handle leading './' when negated 2",
			pattern: "!./foo/bar/*.js",
			expected: ScanState{
				Input: "!./foo/bar/*.js", Prefix: "!./", Start: 3, Base: "foo/bar", Glob: "*.js",
				IsBrace: false, IsBracket: false, IsGlob: true, IsGlobstar: false,
				IsExtglob: false, Negated: true, NegatedExtglob: false,
			},
		},
		{
			name:    "noext option",
			pattern: "./foo/bar/*.js",
			opts:    &Options{Noext: true},
			expected: ScanState{
				Input: "./foo/bar/*.js", Prefix: "./", Start: 2, Base: "foo/bar/*.js", Glob: "",
				IsBrace: false, IsBracket: false, IsGlob: false, IsGlobstar: false,
				IsExtglob: false, Negated: false, NegatedExtglob: false,
			},
		},
		{
			name:    "nonegate option",
			pattern: "!foo/bar/*.js",
			opts:    &Options{Nonegate: true},
			expected: ScanState{
				Input: "!foo/bar/*.js", Prefix: "", Start: 0, Base: "!foo/bar", Glob: "*.js",
				IsBrace: false, IsBracket: false, IsGlob: true, IsGlobstar: false,
				IsExtglob: false, Negated: false, NegatedExtglob: false,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			state := Scan(tc.pattern, tc.opts)
			if state.Input != tc.expected.Input {
				t.Errorf("Input: want %q, got %q", tc.expected.Input, state.Input)
			}
			if state.Prefix != tc.expected.Prefix {
				t.Errorf("Prefix: want %q, got %q", tc.expected.Prefix, state.Prefix)
			}
			if state.Start != tc.expected.Start {
				t.Errorf("Start: want %d, got %d", tc.expected.Start, state.Start)
			}
			if state.Base != tc.expected.Base {
				t.Errorf("Base: want %q, got %q", tc.expected.Base, state.Base)
			}
			if state.Glob != tc.expected.Glob {
				t.Errorf("Glob: want %q, got %q", tc.expected.Glob, state.Glob)
			}
			if state.IsBrace != tc.expected.IsBrace {
				t.Errorf("IsBrace: want %v, got %v", tc.expected.IsBrace, state.IsBrace)
			}
			if state.IsBracket != tc.expected.IsBracket {
				t.Errorf("IsBracket: want %v, got %v", tc.expected.IsBracket, state.IsBracket)
			}
			if state.IsGlob != tc.expected.IsGlob {
				t.Errorf("IsGlob: want %v, got %v", tc.expected.IsGlob, state.IsGlob)
			}
			if state.IsGlobstar != tc.expected.IsGlobstar {
				t.Errorf("IsGlobstar: want %v, got %v", tc.expected.IsGlobstar, state.IsGlobstar)
			}
			if state.IsExtglob != tc.expected.IsExtglob {
				t.Errorf("IsExtglob: want %v, got %v", tc.expected.IsExtglob, state.IsExtglob)
			}
			if state.Negated != tc.expected.Negated {
				t.Errorf("Negated: want %v, got %v", tc.expected.Negated, state.Negated)
			}
			if state.NegatedExtglob != tc.expected.NegatedExtglob {
				t.Errorf("NegatedExtglob: want %v, got %v", tc.expected.NegatedExtglob, state.NegatedExtglob)
			}
			if tc.expected.Slashes != nil && !reflect.DeepEqual(state.Slashes, tc.expected.Slashes) {
				t.Errorf("Slashes: want %v, got %v", tc.expected.Slashes, state.Slashes)
			}
			if tc.expected.Parts != nil && !reflect.DeepEqual(state.Parts, tc.expected.Parts) {
				t.Errorf("Parts: want %v, got %v", tc.expected.Parts, state.Parts)
			}
		})
	}
}

func TestScan_Base(t *testing.T) {
	tests := []struct {
		pattern string
		opts    *Options
		base    string
	}{
		{"./(a|b)", nil, ""},
		{".", nil, "."},
		{".*", nil, ""},
		{"/.*", nil, "/"},
		{"/.*/", nil, "/"},
		{"a/.*/b", nil, "a"},
		{"a*/.*/b", nil, ""},
		{"*/a/b/c", nil, ""},
		{"*", nil, ""},
		{"*/", nil, ""},
		{"*/*", nil, ""},
		{"*/*/", nil, ""},
		{"**", nil, ""},
		{"**/", nil, ""},
		{"**/*", nil, ""},
		{"**/*/", nil, ""},
		{"/*.js", nil, "/"},
		{"*.js", nil, ""},
		{"**/*.js", nil, ""},
		{"/root/path/to/*.js", nil, "/root/path/to"},
		{"[a-z]", nil, ""},
		{"chapter/foo [bar]/", nil, "chapter"},
		{"path/!/foo", nil, "path/!/foo"},
		{"path/!/foo/", nil, "path/!/foo/"},
		{"path/!subdir/foo.js", nil, "path/!subdir/foo.js"},
		{"path/**/*", nil, "path"},
		{"path/**/subdir/foo.*", nil, "path"},
		{"path/*/foo", nil, "path"},
		{"path/*/foo/", nil, "path"},
		{"path/+/foo", nil, "path/+/foo"},
		{"path/+/foo/", nil, "path/+/foo/"},
		{"path/?/foo", nil, "path"},
		{"path/?/foo/", nil, "path"},
		{"path/@/foo", nil, "path/@/foo"},
		{"path/@/foo/", nil, "path/@/foo/"},
		{"path/[a-z]", nil, "path"},
		{"path/subdir/**/foo.js", nil, "path/subdir"},
		{"path/to/*.js", nil, "path/to"},

		{"path/\\*\\*/subdir/foo.*", nil, "path/\\*\\*/subdir"},
		{"path/\\[\\*\\]/subdir/foo.*", nil, "path/\\[\\*\\]/subdir"},
		{"path/\\[foo bar\\]/subdir/foo.*", nil, "path/\\[foo bar\\]/subdir"},
		{"path/\\[bar]/", nil, "path/\\[bar]/"},
		{"path/\\[bar]", nil, "path/\\[bar]"},
		{"[bar]", nil, ""},
		{"[bar]/", nil, ""},
		{"./\\[bar]", nil, "\\[bar]"},
		{"\\[bar]/", nil, "\\[bar]/"},
		{"\\[bar\\]/", nil, "\\[bar\\]/"},
		{"[bar\\]/", nil, "[bar\\]/"},
		{"path/foo \\[bar]/", nil, "path/foo \\[bar]/"},
		{"\\[bar]", nil, "\\[bar]"},
		{"[bar\\]", nil, "[bar\\]"},

		{"path", nil, "path"},
		{"path/foo", nil, "path/foo"},
		{"path/foo/", nil, "path/foo/"},
		{"path/foo/bar.js", nil, "path/foo/bar.js"},

		{"js/*.js", nil, "js"},
		{"js/**/test/*.js", nil, "js"},
		{"js/test/wow.js", nil, "js/test/wow.js"},
		{"js/t[a-z]st}/*.js", nil, "js"},
		{"js/t+(wo|est)/*.js", nil, "js"},

		{"(a|b)", nil, ""},
		{"foo/(a|b)", nil, "foo"},
		{"/(a|b)", nil, "/"},
		{"a/(b c)", nil, "a"},
		{"foo/(b c)/baz", nil, "foo"},
		{"a/(b c)/", nil, "a"},
		{"a/(b c)/d", nil, "a"},
		{"a/(b c)", &Options{Noparen: true}, "a/(b c)"},
		{"a/(b c)/", &Options{Noparen: true}, "a/(b c)/"},
		{"a/(b c)/d", &Options{Noparen: true}, "a/(b c)/d"},
		{"foo/(b c)/baz", &Options{Noparen: true}, "foo/(b c)/baz"},
		{"path/(foo bar)/subdir/foo.*", &Options{Noparen: true}, "path/(foo bar)/subdir"},
		{"a/\\(b c)", nil, "a/\\(b c)"},
		{"a/\\+\\(b c)/foo", nil, "a/\\+\\(b c)/foo"},
		{"js/t(wo|est)/*.js", nil, "js"},
		{"js/t/(wo|est)/*.js", nil, "js/t"},
		{"path/(foo bar)/subdir/foo.*", nil, "path"},
		{"path/(foo/bar|baz)", nil, "path"},
		{"path/(foo/bar|baz)/", nil, "path"},
		{"path/(to|from)", nil, "path"},
		{"path/\\(foo/bar|baz)/", nil, "path/\\(foo/bar|baz)/"},
		{"path/\\*(a|b)", nil, "path"},
		{"path/\\*(a|b)/subdir/foo.*", nil, "path"},
		{"path/\\*/(a|b)/subdir/foo.*", nil, "path/\\*"},
		{"path/\\*\\(a\\|b\\)/subdir/foo.*", nil, "path/\\*\\(a\\|b\\)/subdir"},

		{"path/!(to|from)", nil, "path"},
		{"path/*(to|from)", nil, "path"},
		{"path/+(to|from)", nil, "path"},
		{"path/?(to|from)", nil, "path"},
		{"path/@(to|from)", nil, "path"},

		{"path/{to,from}", nil, "path"},
		{"path/{foo,bar}/", nil, "path"},
		{"js/{src,test}/*.js", nil, "js"},
		{"{a,b}", nil, ""},
		{"/{a,b}", nil, "/"},
		{"/{a,b}/", nil, "/"},
		{"js/test{0..9}/*.js", nil, "js"},

		{"path/{,/,bar/baz,qux}/", &Options{Unescape: true}, "path"},
		{"path/\\{,/,bar/baz,qux}/", &Options{Unescape: true}, "path/{,/,bar/baz,qux}/"},
		{"path/\\{,/,bar/baz,qux\\}/", &Options{Unescape: true}, "path/{,/,bar/baz,qux}/"},
		{"/{,/,bar/baz,qux}/", &Options{Unescape: true}, "/"},
		{"/\\{,/,bar/baz,qux}/", &Options{Unescape: true}, "/{,/,bar/baz,qux}/"},
		{"{,/,bar/baz,qux}", &Options{Unescape: true}, ""},
		{"\\{,/,bar/baz,qux\\}", &Options{Unescape: true}, "{,/,bar/baz,qux}"},
		{"\\{,/,bar/baz,qux}/", &Options{Unescape: true}, "{,/,bar/baz,qux}/"},

		{"\\{../,./,\\{bar,/baz},qux}", &Options{Unescape: true}, "{../,./,{bar,/baz},qux}"},
		{"\\{../,./,\\{bar,/baz},qux}/", &Options{Unescape: true}, "{../,./,{bar,/baz},qux}/"},
		{"path/\\{,/,bar/{baz,qux}}/", &Options{Unescape: true}, "path/{,/,bar/{baz,qux}}/"},
		{"path/\\{../,./,\\{bar,/baz},qux}/", &Options{Unescape: true}, "path/{../,./,{bar,/baz},qux}/"},
		{"path/\\{../,./,{bar,/baz},qux}/", &Options{Unescape: true}, "path/{../,./,{bar,/baz},qux}/"},
		{"path/{,/,bar/\\{baz,qux}}/", &Options{Unescape: true}, "path"},

		{"\\{foo,bar\\}", &Options{Unescape: true}, "{foo,bar}"},
		{"\\{foo,bar\\}/", &Options{Unescape: true}, "{foo,bar}/"},
		{"\\{foo,bar}/", &Options{Unescape: true}, "{foo,bar}/"},
		{"path/\\{foo,bar}/", &Options{Unescape: true}, "path/{foo,bar}/"},

		{"one/{foo,bar}/**/{baz,qux}/*.txt", nil, "one"},
		{"two/baz/**/{abc,xyz}/*.js", nil, "two/baz"},
		{"foo/{bar,baz}/**/aaa/{bbb,ccc}", nil, "foo"},
	}

	for _, tc := range tests {
		t.Run(tc.pattern, func(t *testing.T) {
			state := Scan(tc.pattern, tc.opts)
			if state.Base != tc.base {
				t.Errorf("Expected Base %q, got %q", tc.base, state.Base)
			}
		})
	}
}

func TestScan_Parts(t *testing.T) {
	tests := []struct {
		pattern string
		parts   []string
	}{
		{"./foo", []string{"foo"}},
		{"../foo", []string{"..", "foo"}},
		{"foo/bar", []string{"foo", "bar"}},
		{"foo/*", []string{"foo", "*"}},
		{"foo/**", []string{"foo", "**"}},
		{"foo/**/*", []string{"foo", "**", "*"}},
		{"フォルダ/**/*", []string{"フォルダ", "**", "*"}},
		{"foo/!(abc)", []string{"foo", "!(abc)"}},
		{"c/!(z)/v", []string{"c", "!(z)", "v"}},
		{"c/@(z)/v", []string{"c", "@(z)", "v"}},
		{"foo/(bar|baz)", []string{"foo", "(bar|baz)"}},
		{"foo/(bar|baz)*", []string{"foo", "(bar|baz)*"}},
		{"**/*(W*, *)*", []string{"**", "*(W*, *)*"}},
		{"a/**@(/x|/z)/*.md", []string{"a", "**@(/x|/z)", "*.md"}},
		{"foo/(bar|baz)/*.js", []string{"foo", "(bar|baz)", "*.js"}},
		{"XXX/*/*/12/*/*/m/*/*", []string{"XXX", "*", "*", "12", "*", "*", "m", "*", "*"}},
		{"foo/\\\"**\\\"/bar", []string{"foo", "\\\"**\\\"", "bar"}},
		{"[0-9]/[0-9]", []string{"[0-9]", "[0-9]"}},
		{"foo/[0-9]/[0-9]", []string{"foo", "[0-9]", "[0-9]"}},
		{"foo[0-9]/bar[0-9]", []string{"foo[0-9]", "bar[0-9]"}},
	}

	for _, tc := range tests {
		t.Run(tc.pattern, func(t *testing.T) {
			state := Scan(tc.pattern, &Options{Parts: true})
			if !reflect.DeepEqual(state.Parts, tc.parts) {
				t.Errorf("Expected Parts %v, got %v", tc.parts, state.Parts)
			}
		})
	}
}
