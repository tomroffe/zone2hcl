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
package cmd

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/spf13/cobra"
	"github.com/tomroffe/zone2hcl/pkg/list"
	"github.com/tomroffe/zone2hcl/pkg/utils"
)

// zoneCmd represents the zone command
var zoneCmd = &cobra.Command{
	Use:   "zone [zone name]",
	Short: "Generate hosted zone Terraform resource",
	Args:  utils.VaildateDomain,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadDefaultConfig(context.TODO())

		if err != nil {
			log.Fatalf("Unable to load config")
		}

		svc := route53.NewFromConfig(cfg)
		zoneInput := route53.ListHostedZonesByNameInput{
			DNSName: &args[0],
		}
		listResources, _ := cmd.Flags().GetBool("records")
		list.ListZone(svc, &zoneInput, listResources)
	},
}

func init() {
	rootCmd.AddCommand(zoneCmd)
	var flagName string = "records"
	var shortHand = strings.ToLower(string([]rune(flagName)[0]))
	zoneCmd.Flags().BoolP(flagName, shortHand, false, "return zone record resources as well as zone resources")
}
