package main

import (
	"reflect"
	"testing"
)

func TestGetURLsFromHTML(t *testing.T) {

	tests := []struct {
		name      string
		inputHTML string
		expected  []string
	}{
		{
			name: "external",
			inputHTML: `
<html>
    <body>
        <a href="https://blog.boot.dev"><span>Go to Boot.dev, you React Andy</span></a>
    </body>
</html>
			`,
			expected: []string{"https://blog.boot.dev"},
		},
		{
			name: "internal",
			inputHTML: `
<html>
    <body>
        <a href="/test"><span>Go to Boot.dev, you React Andy</span></a>
    </body>
</html>
			`,
			expected: []string{"boot.dev/test"},
		},
		{
			name: "absolute and relative URLs",
			inputHTML: `
<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected: []string{"boot.dev/path/one", "https://other.com/path/one"},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.inputHTML, "boot.dev")
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
