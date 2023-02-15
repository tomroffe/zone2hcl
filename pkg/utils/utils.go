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
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/spf13/cobra"
)

func FormatName(name string) (string, string) {
	// Strip the trailing dot
	fqdn := strings.TrimRight(name, ".")
	// Replace '.' with '_'. Format for TF HCL Unique Resource Name
	resourceName := strings.ReplaceAll(fqdn, ".", "_")
	return fqdn, resourceName
}

func CreateFileAndRootBody() (*hclwrite.File, *hclwrite.Body) {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	return f, rootBody
}

func IsDomain(domainname string) bool {
	re, err := regexp.Compile(`(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]`)
	if err != nil {
		return false
	}

	if re.MatchString(domainname) {
		return true
	} else {
		return false
	}
}

func VaildateDomain(cmd *cobra.Command, args []string) error {
	if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
		return err
	}
	if IsDomain(args[0]) {
		return nil
	}
	return fmt.Errorf("domain name is invalid: %s", args[0])
}
