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
	"bytes"
	"testing"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestFormatName(t *testing.T) {
	var tests = []struct {
		a              string
		want_a, want_b string
	}{
		{"example.com.", "example.com", "example_com"},
		{"www.example.co.uk.", "www.example.co.uk", "www_example_co_uk"},
		{"an.another.example.io.", "an.another.example.io", "an_another_example_io"},
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

func TestCreateFileAndRootBody(t *testing.T) {
	hclFile, hclBody := CreateFileAndRootBody()
	assert.IsType(t, &hclwrite.File{}, hclFile)
	assert.IsType(t, &hclwrite.Body{}, hclBody)

}

func TestIsDomain(t *testing.T) {
	viper.Set("DomainNameValidationFilterRegEx", `(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]`)
	var tests = []struct {
		a      string
		want_a bool
	}{
		{"example.com.", true},
		{"www.example.co.uk.", true},
		{"an.another.example.io.", true},
		{"this is some text", false},
		{"0-example.com.br", true},
		{"xn--example.com", true},
		{"xn--diseolatinoamericano-76b.com", true},
		{"xn--d1ai6ai.xn--p1ai", true},
		{"xn--rksmrgs-5wao1o.josefsson.org", true},
		{"xn--99zt52a.w3.example.ac.jp", true},
		{"this.is.a.test.com.au", true},
		{"xn--fiqs8s.asia", true},
		{"Another Test", false},
	}

	for _, tt := range tests {
		t.Run(tt.a, func(t *testing.T) {
			ansA, err := IsDomain(tt.a)
			if ansA != tt.want_a {
				t.Errorf("Got %t, wanted %t", ansA, tt.want_a)
			}
			assert.Nil(t, err)
		})
	}
}

func TestIsDomainRegexNotComplining(t *testing.T) {
	// pass a non-compilable regex to viper for IsDomain function to pickup
	viper.Set("DomainNameValidationFilterRegEx", `*$)(`)

	t.Run("BadRegex", func(t *testing.T) {
		_, err := IsDomain("test")
		assert.Error(t, err)
	})
}

func emptyRun(*cobra.Command, []string) {}

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	_, err = root.ExecuteC()

	return buf.String(), err
}

func TestValidateDomain(t *testing.T) {
	var good_regex string = `(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]`
	var bad_regex string = `*$)(`
	var tests = []struct {
		name    string
		command *cobra.Command
		args    []string
		want    string
		pass    bool
		regex   string
	}{
		{"valid_domain", &cobra.Command{
			Use:  "testing",
			Args: cobra.MatchAll(cobra.ExactArgs(1), VaildateDomain),
			Run:  emptyRun,
		}, []string{"testing.com"}, "", true, good_regex},
		{"invalid_domain", &cobra.Command{
			Use:  "testing",
			Args: cobra.MatchAll(cobra.ExactArgs(1), VaildateDomain),
			Run:  emptyRun,
		}, []string{"testing"}, "Error: domain name is invalid testing", false, good_regex},
		{"too_many_args", &cobra.Command{
			Use:  "testing",
			Args: cobra.MatchAll(cobra.ExactArgs(1), VaildateDomain),
			Run:  emptyRun,
		}, []string{"testing", "testing2"}, "Error: accepts 1 arg(s), received 2", false, good_regex},
		{"no_args", &cobra.Command{
			Use:  "testing",
			Args: cobra.MatchAll(cobra.ExactArgs(1), VaildateDomain),
			Run:  emptyRun,
		}, []string{""}, "Error: domain name is invalid", false, good_regex},
		{"bad_domain_filter_validation_regex", &cobra.Command{
			Use:  "testing",
			Args: cobra.MatchAll(cobra.ExactArgs(1), VaildateDomain),
			Run:  emptyRun,
		}, []string{""}, "Error: bad regex. compilation error: *$)(", false, bad_regex},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("DomainNameValidationFilterRegEx", tt.regex)
			output, err := executeCommand(tt.command, tt.args...)

			if tt.pass {
				assert.EqualValues(t, tt.want, output)
			} else {
				assert.Error(t, err)
				assert.Contains(t, output, tt.want)
			}
		})
	}
}
