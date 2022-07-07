func listZoneCount(svc *route53.Client) {
	resp, err := svc.GetHostedZoneCount(context.TODO(), &route53.GetHostedZoneCountInput{})
	if err != nil {
		log.Fatalf("Unable to list zone count.")
	}

	count := resp.HostedZoneCount
	fmt.Println("Zone Count: ", *count)
}

func printZoneImportStatement(zone *types.HostedZone) {
	_, zoneId, _ := strings.Cut(*zone.Id, "/hostedzone/")
	domainName := strings.TrimRight(*zone.Name, ".")
	resourceName := strings.ReplaceAll(domainName, ".", "_")
	fmt.Printf("# terraform import aws_route53_zone.%s %s\n\n", resourceName, zoneId)
}

func listResourceRecord(svc *route53.Client, zone *types.HostedZone) {
	resp, err := svc.ListResourceRecordSets(context.TODO(), &route53.ListResourceRecordSetsInput{
		HostedZoneId: zone.Id,
	})
	if err != nil {
		log.Fatalf("Unable to get the zone resource record.")
	}

	for _, resource := range resp.ResourceRecordSets {
		printResource(zone, &resource)
	}
}

func printResource(zone *types.HostedZone, resource *types.ResourceRecordSet) {
	fmt.Println("Name: ", *resource.Name)
	fmt.Println("Type: ", resource.Type)

	f := hclwrite.NewEmptyFile()

	// Remove trailing .(dot), Replace remaining with _(underscore)
	domainName := strings.TrimRight(*zone.Name, ".")
	recordName := strings.ReplaceAll(domainName, ".", "_")
	zoneBlock := createResourceBlock(f, "aws_route53_record", recordName)
	zoneBody := zoneBlock.Body()
	zoneBody.SetAttributeValue("name", cty.StringVal(domainName))

	if len(resource.ResourceRecords) > 0 {
		fmt.Println("TLL: ")
		zoneBody.SetAttributeValue("ttl", cty.NumberIntVal(*resource.TTL))
		for _, record := range resource.ResourceRecords {
			fmt.Println("Value = ", *record.Value)
		}
	} else {
		fmt.Println("Target Alias: ", *resource.AliasTarget.DNSName)
		fmt.Println("Zone ID: ", *resource.AliasTarget.HostedZoneId)
	}

	fmt.Printf("%s", f.Bytes())
}