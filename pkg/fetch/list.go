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
package fetch

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func ListZone(ctx context.Context, svc *route53.Client, zoneInput *route53.ListHostedZonesByNameInput) *types.HostedZone {
	// Fetching the first zone returned by ListHostedZonesByName by using to ListHostedZonesByNameInput
	// to ensure first result was at the top of the results[0].
	resp, err := svc.ListHostedZonesByName(ctx, zoneInput)

	if err != nil {
		log.Fatalf("Unable to fetch zone %s.", *zoneInput.DNSName)
	}

	return &resp.HostedZones[0]
}

func ListZones(ctx context.Context, svc *route53.Client) *[]types.HostedZone {
	resp, err := svc.ListHostedZonesByName(ctx, &route53.ListHostedZonesByNameInput{})

	if err != nil {
		log.Fatalf("Unable to fetch zone list.")
	}

	return &resp.HostedZones

}

func ListResourceRecords(ctx context.Context, svc *route53.Client, zone *types.HostedZone) *[]types.ResourceRecordSet {
	resp, err := svc.ListResourceRecordSets(ctx, &route53.ListResourceRecordSetsInput{
		HostedZoneId: zone.Id,
	})

	if err != nil {
		log.Fatalf("Unable to fetch zones' resource record's")
	}

	return &resp.ResourceRecordSets
}
