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
package list

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/tomroffe/zone2hcl/pkg/generate"
)

func ListZone(svc *route53.Client, zoneInput *route53.ListHostedZonesByNameInput, listResources bool) {
	// Fetching the first zone returned by ListHostedZonesByName by using to ListHostedZonesByNameInput
	// to ensure first result was at the top of the results[0].
	resp, err := svc.ListHostedZonesByName(context.TODO(), zoneInput)

	if err != nil {
		log.Fatalf("Unable to fetch zone list")
	}

	zoneHCL := generate.GenerateZone(&resp.HostedZones[0])
	fmt.Printf("%s", zoneHCL.Bytes())

	if listResources {
		ListResourceRecords(svc, &resp.HostedZones[0])
	}

}

func ListZones(svc *route53.Client, listResources bool) {
	resp, err := svc.ListHostedZonesByName(context.TODO(), &route53.ListHostedZonesByNameInput{})

	if err != nil {
		log.Fatalf("Unable to fetch zone list")
	}

	for _, zone := range resp.HostedZones {
		zoneHCL := generate.GenerateZone(&zone)
		fmt.Printf("%s", zoneHCL.Bytes())

		if listResources {
			ListResourceRecords(svc, &zone)
		}
	}

}

func ListResourceRecords(svc *route53.Client, zone *types.HostedZone) {
	resp, err := svc.ListResourceRecordSets(context.TODO(), &route53.ListResourceRecordSetsInput{
		HostedZoneId: zone.Id,
	})

	if err != nil {
		log.Fatalf("Unable to fetch zone resource record list")
	}

	for _, recordSet := range resp.ResourceRecordSets {
		recordHCL := generate.GenerateRecord(zone, &recordSet)
		fmt.Printf("%s", recordHCL.Bytes())
	}
}

func GetHostedZoneID(svc *route53.Client, zoneFQDN string) (zone *types.HostedZone) {
	zoneInput := route53.ListHostedZonesByNameInput{
		DNSName: &zoneFQDN,
	}
	resp, err := svc.ListHostedZonesByName(context.TODO(), &zoneInput)

	if err != nil {
		log.Fatalf("Unable to look-up zone ID.\n%s", err)
	}

	// return the first types.HostedZone ordered by the requested zones fqdn.
	return &resp.HostedZones[0]
}
