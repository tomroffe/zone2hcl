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
	"github.com/spf13/viper"
	"github.com/tomroffe/zone2hcl/pkg/fetch"
)

// zonesCmd represents the zones command
var zonesCmd = &cobra.Command{
	Use:   "zones",
	Short: "Generate hosted zones Terraform resource",
	Long:  ``,
	Run:   ZonesCmd,
}

func ZonesCmd(cmd *cobra.Command, args []string) {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		log.Fatalf("Unable to load config")
	}

	svc := route53.NewFromConfig(cfg)
	// listResources, _ := cmd.Flags().GetBool("records")
	fetch.ListZones(ctx, svc)
}

func init() {
	rootCmd.AddCommand(zonesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// zonesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	var flagName string = "records"
	var shortHand = strings.ToLower(string([]rune(flagName)[0]))
	zonesCmd.Flags().BoolP(flagName, shortHand, false, "return zone record resources as well as zone resources")
	viper.BindPFlag(flagName, zonesCmd.Flags().Lookup(flagName))
}
