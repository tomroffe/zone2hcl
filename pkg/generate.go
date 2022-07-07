package generate

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func generateZone(zone *types.HostedZone) *hclwrite.File {
	// Create File and Root Body
	file, rootBody := createFileAndRootBody()
	// Zone FQDN and Zone Resource Name
	fqdn, resourceName := formatName(*zone.Name)

	// Create Resource
	zoneBlock := rootBody.AppendNewBlock("resource", []string{"aws_route53_zone", resourceName})
	zoneBody := zoneBlock.Body()
	zoneBody.SetAttributeValue("name", cty.StringVal(fqdn))

	// Add New Line To Root Body After Resource Addition
	rootBody.AppendNewline()
	return file
}

func generateRecord(zone *types.HostedZone, resourceSet *types.ResourceRecordSet) *hclwrite.File {
	// Create File and Root Body
	file, rootBody := createFileAndRootBody()
	// ResourceRecord FQDN and ResourceRecord Resource Name
	fqdn, resourceName := formatName(*resourceSet.Name)
	// ResourceRecord Parent/Root Zone FQDN and ResourceRecord Parent/Root Zone Resource Name
	_, zoneResourceName := formatName(*zone.Name)

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

func formatName(name string) (string, string) {
	// Strip the trailing dot
	fqdn := strings.TrimRight(name, ".")
	// Replace '.' with '_'. Format for TF HCL Unique Resource Name
	resourceName := strings.ReplaceAll(fqdn, ".", "_")
	return fqdn, resourceName
}

func createFileAndRootBody() (*hclwrite.File, *hclwrite.Body) {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	return f, rootBody
}
