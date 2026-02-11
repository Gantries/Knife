package lang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatL(t *testing.T) {
	assert.True(t, FormatL(10, "hello %s", "world") == "hello worl")
	assert.True(t, FormatL(15, "hello %s", "world") == "hello world")

	s := "hello world"
	assert.True(t, Substring(&s, 0, 5) == "hello")
	assert.True(t, Substring(&s, 0, 50) == "hello world")
	assert.True(t, Substring(&s, -1, 50) == "")
	assert.True(t, Substring(&s, -1, -2) == "")
	assert.True(t, Substring(&s, 3, 1) == "")
	assert.True(t, Substring(&s, 3, 3) == "")
}

func TestIsBlank(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", true},     // 空字符串
		{" ", true},    // 只包含空格的字符串
		{"\t", true},   // 只包含制表符的字符串
		{"\n", true},   // 只包含换行符的字符串
		{"\r", true},   // 只包含回车符的字符串
		{"   ", true},  // 包含多个空白字符的字符串
		{"a", false},   // 非空白字符的字符串
		{" b ", false}, // 包含前后空白字符的非空白字符串
		{"abc", false}, // 普通字符串
	}

	for _, test := range tests {
		result := IsBlank(test.input)
		if result != test.expected {
			t.Errorf("IsBlank(%q) expected %v, got %v", test.input, test.expected, result)
		}
	}
}

// TestJoinWith tests the JoinWith function with various test cases.
func TestJoinWith(t *testing.T) {
	tests := []struct {
		Input    []*string // Input is a slice of string pointers
		Expected string    // Expected is the expected result of the JoinWith function
	}{
		{Input: []*string{stringToPtr("Hello"), stringToPtr("World")}, Expected: "Hello, World"},
		{Input: []*string{stringToPtr("Hello"), nil, stringToPtr("World")}, Expected: "Hello, World"},
		{Input: []*string{nil, stringToPtr("Hello"), stringToPtr("World")}, Expected: "Hello, World"},
		{Input: []*string{}, Expected: ""},
		{Input: []*string{nil}, Expected: ""},
	}

	for _, test := range tests {
		result := JoinWith(", ", test.Input...)
		if result != test.Expected {
			t.Errorf("TestJoinWith failed: expected '%s', got '%s'", test.Expected, result)
		}
	}
}

// stringToPtr is a helper function to convert a string to a *string.
func stringToPtr(s string) *string {
	return &s
}

func TestTrimJoin(t *testing.T) {
	assert.Equal(t, "abc", TrimJoin("", Dup("   a "), Dup(" b"), Dup("c  ")))
	assert.Equal(t, "abc///def/", TrimJoin("//", Dup(" //abc"), Dup("/def///")))
	assert.Equal(t, "abc/def", TrimJoin("/", Dup(" //abc///"), Dup("///def///")))
	assert.Equal(t, "//a/b/c////a////d/e/f///", TrimJoin("/a/", Dup(" //a/b/c///"), Dup("///d/e/f///")))
	assert.Equal(t, "/b/c////a/d/e/f///", TrimJoin("/a", Dup(" /a/b/c///"), Dup("/d/e/f///")))
	assert.Equal(t, "a/b/c/d/e/f", TrimJoin("/", Dup(" /a/b"), Dup("/c/"), Dup("d/e/f")))
	assert.Equal(t, "a/b/c/d/e/f", TrimJoin("/", Dup(" /a/"), Dup(" b/c/  "), Dup("d/e/f")))
	assert.Equal(t, "我/爱/China/！", TrimJoin("/", Dup(" /我/"), Dup(" 爱/China/  "), Dup("！")))
	assert.Equal(t, "我/前缀/爱/China/前缀/！", TrimJoin("前缀/", Dup(" 前缀/我/前缀/"), Dup(" 爱/China/  "), Dup("！")))
}
