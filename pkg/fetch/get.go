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

func GetHostedZoneID(ctx context.Context, svc *route53.Client, zoneFQDN string) (zone *types.HostedZone) {
	zoneInput := route53.ListHostedZonesByNameInput{
		DNSName: &zoneFQDN,
	}
	resp, err := svc.ListHostedZonesByName(ctx, &zoneInput)

	if err != nil {
		log.Fatalf("Unable to look-up zone ID: %s\n%s", zoneFQDN, err)
	}

	// return the first types.HostedZone ordered by the requested zones fqdn.
	return &resp.HostedZones[0]
}
