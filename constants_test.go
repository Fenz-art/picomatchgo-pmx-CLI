package picomatch

import "testing"

func TestGetGlobChars_Posix(t *testing.T) {
	chars := GetGlobChars(false)
	if chars.Sep != "/" {
		t.Errorf("Expected Sep '/', got %q", chars.Sep)
	}
	if chars.Qmark != "[^/]" {
		t.Errorf("Expected Qmark '[^/]', got %q", chars.Qmark)
	}
	if chars.Star != "[^/]*?" {
		t.Errorf("Expected Star '[^/]*?', got %q", chars.Star)
	}
}

func TestGetGlobChars_Windows(t *testing.T) {
	chars := GetGlobChars(true)
	if chars.Sep != "\\" {
		t.Errorf("Expected Sep '\\\\', got %q", chars.Sep)
	}
	if chars.Qmark != `[^\\/]` {
		t.Errorf("Expected Qmark '[^\\\\/]', got %q", chars.Qmark)
	}
}

func TestPosixRegexSource(t *testing.T) {
	expected := map[string]string{
		"alnum":  `a-zA-Z0-9`,
		"alpha":  `a-zA-Z`,
		"digit":  `0-9`,
		"lower":  `a-z`,
		"upper":  `A-Z`,
		"word":   `A-Za-z0-9_`,
		"xdigit": `A-Fa-f0-9`,
	}

	for name, want := range expected {
		got, ok := PosixRegexSource[name]
		if !ok {
			t.Errorf("Missing POSIX class %q", name)
			continue
		}
		if got != want {
			t.Errorf("PosixRegexSource[%q] = %q, want %q", name, got, want)
		}
	}
}

func TestExtglobChars(t *testing.T) {
	chars := ExtglobChars(PosixChars)

	tests := []struct {
		key      byte
		wantType string
	}{
		{'!', "negate"},
		{'?', "qmark"},
		{'+', "plus"},
		{'*', "star"},
		{'@', "at"},
	}

	for _, tc := range tests {
		eg, ok := chars[tc.key]
		if !ok {
			t.Errorf("Missing extglob char %q", tc.key)
			continue
		}
		if eg.Type != tc.wantType {
			t.Errorf("ExtglobChars[%q].Type = %q, want %q", tc.key, eg.Type, tc.wantType)
		}
	}
}
