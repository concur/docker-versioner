package main

import "testing"

func TestProcess(t *testing.T) {
	tests := []struct {
		name       string
		base       string
		pattern    string
		ignorePrev bool
		prerelease string
		metadata   string
		prefix     string
		expected   string
	}{
		{name: "Inc major from base 3.0.0", base: "3.0.0", ignorePrev: true, pattern: "^.*.*", prerelease: "test", expected: "4.0.0-test"},
		{name: "Inc major-minor-patch from base 3.0.0", base: "3.0.0", ignorePrev: true, pattern: "^.^.^", prerelease: "test", expected: "4.1.1-test"},
		{name: "Inc minor from base 1.4.0", base: "1.4.0", ignorePrev: true, pattern: "*.^.*", prerelease: "test", expected: "1.5.0-test"},
		{name: "Inc patch from base 1.4.0", base: "1.4.0", ignorePrev: true, pattern: "*.*.^", prerelease: "test", expected: "1.4.1-test"},
		{name: "Inc major from base 1.4.0", base: "1.4.0", ignorePrev: true, pattern: "^.*.*", prerelease: "test", expected: "2.0.0-test"},
		{name: "Include metadata with base 1.6.0", base: "1.6.0", ignorePrev: true, pattern: "^.*.*", prerelease: "test", expected: "2.0.0-test+patch", metadata: "patch"},
		{name: "Include metadata with base 1.6.0 plus prefix", base: "1.6.0", prefix: "v", ignorePrev: true, pattern: "^.*.*", prerelease: "test", expected: "v2.0.0-test+patch", metadata: "patch"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if res, e := Process(test.base, test.pattern, test.prerelease, test.metadata, test.prefix, test.ignorePrev); e != nil {
				t.Errorf("Failed to get version: %v", e)
			} else {
				if res != test.expected {
					t.Errorf("Retrieved version does not match: Expected: %s - Got: %s", test.expected, res)
				}
			}
		})
	}
}
