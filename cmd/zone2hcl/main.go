package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("Unable to load config")
	}

	svc := route53.NewFromConfig(cfg)
	listZones(svc)
}

func listZones(svc *route53.Client) {
	resp, err := svc.ListHostedZonesByName(context.TODO(), &route53.ListHostedZonesByNameInput{})

	if err != nil {
		log.Fatalf("Unable to fetch zone list")
	}

	for _, zone := range resp.HostedZones {
		zoneHCL := generateZone(&zone)
		fmt.Printf("%s", zoneHCL.Bytes())
	}
}

func generateZone(zone *types.HostedZone) *hclwrite.File {
	// Create File and Root Body
	f, rootBody := createFileAndRootBody()
	// Remove trailing .(dot), Replace remaining with _(underscore)
	domainName := strings.TrimRight(*zone.Name, ".")
	resourceName := strings.ReplaceAll(domainName, ".", "_")
	// Create Resource
	zoneBlock := createResourceBlock(rootBody, "aws_route53_zone", resourceName)
	zoneBody := zoneBlock.Body()
	zoneBody.SetAttributeValue("name", cty.StringVal(domainName))
	// Add New Line To Root Body After Resource Addition
	rootBody.AppendNewline()
	return f
}

func generateRecord(zone *types.HostedZone, resource *types.ResourceRecordSet) *hclwrite.File {
	// Create File and Root Body
	f, rootBody := createFileAndRootBody()
	// Remove trailing .(dot), Replace remaining with _(underscore)
	domainName := strings.TrimRight(*resource.Name, ".")
	recordName := strings.ReplaceAll(domainName, ".", "_")
	resourceRecordBlock := createResourceBlock(rootBody, "aws_route53_record", recordName)
	resourceRecordBody := resourceRecordBlock.Body()

	parentDomain := strings.TrimRight(*zone.Name, ".")
	resourceName := strings.ReplaceAll(parentDomain, ".", "_")
	zoneId := fmt.Sprintf("aws_route53_zone.%s.zone", resourceName)
	resourceToken := hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte(zoneId),
		},
	}
	resourceRecordBody.SetAttributeRaw("zone_id", resourceToken)
	resourceRecordBody.SetAttributeValue("name", cty.StringVal(domainName))
	resourceRecordBody.SetAttributeValue("type", cty.StringVal(string(resource.Type)))

	if len(resource.ResourceRecords) > 0 {
		resourceRecordBody.SetAttributeValue("ttl", cty.NumberIntVal(*resource.TTL))
		// records := []string{}
		// for _, record := range resource.ResourceRecords {
		// 	attr := cty.ListVal([]cty.Value{cty.StringVal(""), cty.StringVal("")}) 
		// 	records = append(records, *record.Value)
		// }
		resourceRecordBody.SetAttributeValue("records",
			cty.ListVal([]cty.Value{cty.StringVal("123"), cty.StringVal("abc")}))
	} else {
		fmt.Println("Target Alias: ", *resource.AliasTarget.DNSName)
		fmt.Println("Zone ID: ", *resource.AliasTarget.HostedZoneId)
	}

	return f
}

func createFileAndRootBody() (*hclwrite.File, *hclwrite.Body) {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	return f, rootBody
}

func createResourceBlock(rootBody *hclwrite.Body, resourceType string, resourceName string) *hclwrite.Block {
	// Remove trailing .(dot), Replace remaining with _(underscore)
	zoneBlock := rootBody.AppendNewBlock("resource", []string{resourceType, resourceName})
	return zoneBlock
}
