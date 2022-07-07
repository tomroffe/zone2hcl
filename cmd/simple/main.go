package main

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func main() {
	type Zone struct {
		Type         string `hcl:"type,label"`
		ResourceName string `hcl:"name,label"`

		Name string `hcl:"name"`
		Desc string `hcl:"description"`
	}

	type Resource struct {
		Type         string `hcl:"type,label"`
		ResourceName string `hcl:"name,label"`

		Name string `hcl:"name"`
	}

	type App struct {
		Zone     *Zone     `hcl:"resourcea,block"`
		Resource *Resource `hcl:"resource,block"`
	}

	app := App{
		Zone: &Zone{
			Type:         "aws_route53_zone",
			ResourceName: "test_com",
			Name:         "test.com.",
			Desc:         "Test Domain",
		},
		Resource: &Resource{
			Type:         "aws_route53_record",
			ResourceName: "dev_test_com",
			Name:         "dev.test.com.",
		},
	}

	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(&app, f.Body())
	fmt.Printf("%s", f.Bytes())

}
