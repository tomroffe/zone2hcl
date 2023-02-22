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
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/stretchr/testify/assert"
)

type mockListHostedZonesByNameAPI func(ctx context.Context, params *route53.ListHostedZonesByNameInput, optFns ...func(*route53.Options)) (*route53.ListHostedZonesByNameOutput, error)

func (m mockListHostedZonesByNameAPI) ListHostedZonesByName(ctx context.Context, params *route53.ListHostedZonesByNameInput, optFns ...func(*route53.Options)) (*route53.ListHostedZonesByNameOutput, error) {
	return m(ctx, params, optFns...)
}

func TestGetHostedZoneID(t *testing.T) {
	domains := []struct {
		id      string
		ref     string
		name    string
		count   int64
		comment string
		private bool
	}{
		{id: "testing", ref: "testing", name: "testing.com.", count: 5, comment: "testing!", private: false},
	}
	cases := []struct {
		client func(t *testing.T) Route53API
		expect []byte
		domain string
		error  bool
	}{
		{
			client: func(t *testing.T) Route53API {
				return mockListHostedZonesByNameAPI(func(ctx context.Context, params *route53.ListHostedZonesByNameInput, optFns ...func(*route53.Options)) (*route53.ListHostedZonesByNameOutput, error) {
					t.Helper()
					return &route53.ListHostedZonesByNameOutput{
						HostedZones: []types.HostedZone{
							{
								CallerReference:        &domains[0].ref,
								Id:                     &domains[0].id,
								Name:                   &domains[0].name,
								ResourceRecordSetCount: &domains[0].count,
								Config: &types.HostedZoneConfig{
									Comment:     &domains[0].comment,
									PrivateZone: domains[0].private,
								},
							},
						},
					}, nil
				})
			},
			expect: []byte(""),
			domain: "testing.com.",
			error:  false,
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			content := GetHostedZoneID(ctx, tt.client(t), tt.domain)

			assert.Equal(t, tt.domain, *content.Name)
		})
	}
}
