// Copyright (c) 2020 Fandom, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import "testing"

func TestParseTimeDescription(t *testing.T) {
	cases := []struct {
		desc     string
		expected float64
	}{
		{"389 days 9343h 3m 28.000000s", 67244608},
		{"22h 14m 53.898438s", 80093.898438},
		{"0.957994s", 0.957994},
	}

	for _, c := range cases {
		seconds := parseTimeDescription(c.desc)

		if seconds != c.expected {
			t.Errorf("parseTimeDescription(%s) == %f, expected %f", c.desc, seconds, c.expected)
		}
	}
}
