package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Shopify/themekit/cmd/wizard"
	"github.com/Shopify/themekit/kit"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		wiz := wizard.New(fmt.Sprintf("Shopify Themekit v%v", kit.ThemeKitVersion))
		domain, _ := wiz.Text("shop domain", validateDomain)
		action, _ := wiz.List("Action", []string{
			"work on an existing theme",
			"duplicate an existing theme",
			"create new from a template",
		})
		fmt.Println(domain, action)
		return nil
	},
}

func validateDomain(value string) []string {
	errors := []string{}
	if len(value) == 0 {
		errors = append(errors, "Must provide a value.")
	} else if !strings.HasSuffix(value, "myshopify.com") {
		errors = append(errors, "Store domain must end in '.myshopify.com'")
	}
	return errors
}
