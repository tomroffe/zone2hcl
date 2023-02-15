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
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/hashicorp/hcl/v2/hclsyntax"
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

func GenerateRecord(zone *types.HostedZone, resourceSet *types.ResourceRecordSet) *hclwrite.File {
	// Create File and Root Body
	file, rootBody := utils.CreateFileAndRootBody()
	// ResourceRecord FQDN and ResourceRecord Resource Name
	fqdn, resourceName := utils.FormatName(*resourceSet.Name)
	// ResourceRecord Parent/Root Zone FQDN and ResourceRecord Parent/Root Zone Resource Name
	_, zoneResourceName := utils.FormatName(*zone.Name)

	resourceName = fmt.Sprintf("%s_%s", strings.ToLower(string(resourceSet.Type)), resourceName)
	resourceRecordBlock := rootBody.AppendNewBlock("resource", []string{"aws_route53_record", resourceName})
	resourceRecordBody := resourceRecordBlock.Body()

	zoneId := fmt.Sprintf("aws_route53_zone.%s.zone", zoneResourceName)
	resourceRecordBody.SetAttributeRaw("zone_id", hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte(zoneId),
		},
	})
	resourceRecordBody.SetAttributeValue("name", cty.StringVal(fqdn))
	resourceRecordBody.SetAttributeValue("type", cty.StringVal(string(resourceSet.Type)))

	if len(resourceSet.ResourceRecords) > 0 {
		resourceRecordBody.SetAttributeValue("ttl", cty.NumberIntVal(*resourceSet.TTL))
		records := []cty.Value{}
		for _, record := range resourceSet.ResourceRecords {
			records = append(records, cty.StringVal(*record.Value))
		}
		resourceRecordBody.SetAttributeValue("records",
			cty.ListVal(records))
	} else {
		resourceRecordBody.AppendNewline()

		aliasBlock := resourceRecordBody.AppendNewBlock("alias", nil)
		aliasBody := aliasBlock.Body()
		aliasZoneVar := fmt.Sprintf("var.%s_alias_zone", resourceName)
		aliasBody.SetAttributeRaw("name", hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenIdent,
				Bytes: []byte(aliasZoneVar),
			},
		})
		aliasZoneIdVar := fmt.Sprintf("var.%s_alias_zone_id", resourceName)
		aliasBody.SetAttributeRaw("zone_id", hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenIdent,
				Bytes: []byte(aliasZoneIdVar),
			},
		})
		aliasZoneEvalTargetHealthVar := fmt.Sprintf("var.%s_alias_healthcheck", resourceName)
		aliasBody.SetAttributeRaw("evaluate_target_health", hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenIdent,
				Bytes: []byte(aliasZoneEvalTargetHealthVar),
			},
		})
		// fmt.Println("Target Alias: ", *resource.AliasTarget.DNSName)
		// fmt.Println("Zone ID: ", *resource.AliasTarget.HostedZoneId)
	}

	return file
}
