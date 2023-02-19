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
package generate

import (
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/tomroffe/zone2hcl/pkg/utils"
	"github.com/zclconf/go-cty/cty"
)

func GenerateZone(zone *types.HostedZone) *hclwrite.File {
	// Create File and Root Body
	file, rootBody := utils.CreateFileAndRootBody()
	// Zone FQDN and Zone Resource Name
	fqdn, resourceName := utils.FormatName(*zone.Name)

	// Create Resource
	zoneBlock := rootBody.AppendNewBlock("resource", []string{"aws_route53_zone", resourceName})
	zoneBody := zoneBlock.Body()
	zoneBody.SetAttributeValue("name", cty.StringVal(fqdn))

	// Add New Line To Root Body After Resource Addition
	rootBody.AppendNewline()
	return file
}
