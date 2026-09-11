package jsonvector

import (
	"encoding/json"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppendEscape(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "no special characters",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "double quote",
			input:    `Hello "World"`,
			expected: `Hello \"World\"`,
		},
		{
			name:     "backslash",
			input:    `C:\Program Files\App`,
			expected: `C:\\Program Files\\App`,
		},
		{
			name:     "newline",
			input:    "Line1\nLine2",
			expected: "Line1\\nLine2",
		},
		{
			name:     "carriage return",
			input:    "Line1\rLine2",
			expected: "Line1\\rLine2",
		},
		{
			name:     "tab",
			input:    "Col1\tCol2",
			expected: "Col1\\tCol2",
		},
		{
			name:     "form feed",
			input:    "Page1\fPage2",
			expected: "Page1\\fPage2",
		},
		{
			name:     "backspace",
			input:    "Hello\bWorld",
			expected: "Hello\\bWorld",
		},
		{
			name:     "DEL character",
			input:    "Hello\x7fWorld",
			expected: "Hello\\u007fWorld",
		},
		{
			name:     "null character",
			input:    "Hello\x00World",
			expected: "Hello\\u0000World",
		},
		{
			name:     "control characters 0x01-0x1F",
			input:    "Hello\x01World\x02Test",
			expected: "Hello\\u0001World\\u0002Test",
		},
		{
			name:     "all control characters",
			input:    string([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F}),
			expected: "\\u0000\\u0001\\u0002\\u0003\\u0004\\u0005\\u0006\\u0007\\b\\t\\n\\u000b\\f\\r\\u000e\\u000f\\u0010\\u0011\\u0012\\u0013\\u0014\\u0015\\u0016\\u0017\\u0018\\u0019\\u001a\\u001b\\u001c\\u001d\\u001e\\u001f",
		},
		{
			name:     "mixed special characters",
			input:    "Hello \"World\"\n\t\\End",
			expected: "Hello \\\"World\\\"\\n\\t\\\\End",
		},
		{
			name:     "unicode characters",
			input:    "Привет 世界",
			expected: "Привет 世界",
		},
		{
			name:     "long string with many escapes",
			input:    func() string { return `{"key":"value"}` + "\n\t\\" }(),
			expected: "{\\\"key\\\":\\\"value\\\"}\\n\\t\\\\",
		},
		{
			name:     "all special chars in one string",
			input:    `" \ \n \r \t \f \b ` + "\x7f\x00",
			expected: `\" \\ \\n \\r \\t \\f \\b \u007f\u0000`,
		},
		{
			name:     "only special chars",
			input:    "\"\n\r\t\f\b\x7f\x00\\",
			expected: "\\\"\\n\\r\\t\\f\\b\\u007f\\u0000\\\\",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := []byte(tt.input)
			expected := []byte(tt.expected)
			result := AppendEscape([]byte(""), input)

			assert.Equal(t, expected, result, "AppendEscape() mismatch\nInput: %q", input)

			jsonStr := `"` + string(result) + `"`
			var decoded string
			err := json.Unmarshal([]byte(jsonStr), &decoded)
			require.NoError(t, err, "Result is not valid JSON")
		})
	}
}

func BenchmarkAppendEscape(b *testing.B) {
	genStage := func(size int, escapeProbability float64) []byte {
		data := make([]byte, size)

		for i := range data {
			if rand.Float64() < escapeProbability {
				specialChars := []byte{'"', '\\', '\n', '\r', '\t', '\f', '\b', 0x7F, 0x00}
				data[i] = specialChars[rand.Intn(len(specialChars))]
			} else {
				data[i] = byte(rand.Intn(0x5F) + 0x20)
				for data[i] == '"' || data[i] == '\\' {
					data[i] = byte(rand.Intn(0x5F) + 0x20)
				}
			}
		}
		return data
	}

	benchmarks := []struct {
		name string
		size int
		prob float64
	}{
		{"Small_NoEscapes", 100, 0.0},
		{"Small_FewEscapes", 100, 0.1},
		{"Small_ManyEscapes", 100, 0.5},
		{"Medium_NoEscapes", 1024, 0.0},
		{"Medium_FewEscapes", 1024, 0.1},
		{"Medium_ManyEscapes", 1024, 0.5},
		{"Large_NoEscapes", 10240, 0.0},
		{"Large_FewEscapes", 10240, 0.1},
		{"Large_ManyEscapes", 10240, 0.5},
		{"Huge_NoEscapes", 102400, 0.0},
		{"Huge_FewEscapes", 102400, 0.1},
		{"Huge_ManyEscapes", 102400, 0.5},
		{"VeryHuge_Mixed", 1048576, 0.3}, // 1MB
	}

	for _, bm := range benchmarks {
		data := genStage(bm.size, bm.prob)

		b.Run(bm.name, func(b *testing.B) {
			var buf []byte
			b.SetBytes(int64(len(data)))
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				buf = AppendEscape(buf[:0], data)
			}
		})
	}
}
