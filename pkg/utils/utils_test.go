/*
zone2hcl
Copyright © 2023 Tom Roffe tom.roffe@gmail.com

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package utils

import (
	"testing"
)

func TestFormatName(t *testing.T) {
	var tests = []struct {
		a              string
		want_a, want_b string
	}{
		{"example.com.", "example.com", "example_com"},
		{"www.example.co.uk.", "www.example.co.uk", "www_example_co_uk"},
	}

	for _, tt := range tests {
		t.Run(tt.a, func(t *testing.T) {
			ansA, ansB := FormatName(tt.a)
			if ansA != tt.want_a {
				t.Errorf("Got %s, wanted %s", ansA, tt.want_a)
			}
			if ansB != tt.want_b {
				t.Errorf("Got %s, wanted %s", ansB, tt.want_b)
			}
		})
	}
}
