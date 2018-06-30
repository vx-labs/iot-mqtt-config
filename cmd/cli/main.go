package main

import (
	"log"
	"os"

	vault "github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"
	"github.com/vx-labs/iot-mqtt-config"
)

func CloudFlare() *cobra.Command {
	return &cobra.Command{
		Use:   "cloudflare",
		Short: "cloudflare connection credentials",
		Run: func(_ *cobra.Command, _ []string) {
			vaultConfig := vault.DefaultConfig()
			api, err := vault.NewClient(vaultConfig)
			if err != nil {
				log.Fatal(err)
			}
			api.SetToken(os.Getenv("VAULT_TOKEN"))
			cf, err := config.Cloudflare(api)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(cf.EmailAddress)
		},
	}
}
func main() {
	cmd := &cobra.Command{
		Use:   "config CLI",
		Short: "common configuration keys for platform operators",
	}
	cmd.AddCommand(CloudFlare())
	cmd.Execute()
}
