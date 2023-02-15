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

import "fmt"

type Version struct {
	Major      int
	Minor      int
	PatchLevel int
	Suffix     string
}

var CurrentVersion = Version{
	Major:      0,
	Minor:      1,
	PatchLevel: 0,
	Suffix:     "",
}

func GetVersion(v Version) string {
	return fmt.Sprintf("v%d.%d.%d%s", v.Major, v.Minor, v.PatchLevel, v.Suffix)
}
