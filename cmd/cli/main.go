package main

import (
	"log"
	"os"

	"text/template"

	consul "github.com/hashicorp/consul/api"
	vault "github.com/hashicorp/vault/api"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/vx-labs/iot-mqtt-config"
)

func Template(tpl string) (*template.Template, error) {
	return template.New("").Funcs(promptui.FuncMap).Parse(tpl)
}

func TLS(api *consul.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "tls",
		Short: "display TLS configuration keys",
		Run: func(cmd *cobra.Command, _ []string) {
			tls, _, err := config.TLS(api)
			if err != nil {
				log.Fatalf("failed to get TLS config key: %v", err)
			}
			tpl, err := Template(tls.Template())
			if err != nil {
				log.Fatal(err)
			}
			tpl.Execute(cmd.OutOrStdout(), tls)
		},
	}
}
func CloudFlare(api *vault.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "cloudflare",
		Short: "display cloudflare connection credentials",
		Run: func(cmd *cobra.Command, _ []string) {

			cf, err := config.Cloudflare(api)
			if err != nil {
				log.Fatal(err)
			}
			tpl, err := Template(cf.Template())
			if err != nil {
				log.Fatal(err)
			}
			tpl.Execute(cmd.OutOrStdout(), cf)
		},
	}
}
func main() {
	vaultAPI, err := vault.NewClient(vault.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}
	vaultAPI.SetToken(os.Getenv("VAULT_TOKEN"))

	consulAPI, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}
	cmd := &cobra.Command{
		Use:   "config CLI",
		Short: "common configuration keys for platform operators",
	}
	cmd.AddCommand(CloudFlare(vaultAPI))
	cmd.AddCommand(TLS(consulAPI))
	cmd.Execute()
}
