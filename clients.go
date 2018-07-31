package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	consul "github.com/hashicorp/consul/api"
	vault "github.com/hashicorp/vault/api"
)

func loadVaultToken(api *vault.Client) {
	fallback := func() {
		api.SetToken(os.Getenv("VAULT_TOKEN"))
	}
	_, err := os.Stat("secrets/vault_token")
	if err != nil {
		fallback()
		return
	}
	token, err := ioutil.ReadFile("secrets/vault_token")
	if err != nil {
		fallback()
		return
	}
	api.SetToken(string(token))
	go func() {
		sigUsr1 := make(chan os.Signal, 1)
		signal.Notify(sigUsr1, syscall.SIGUSR1)
		<-sigUsr1
		log.Println("INFO: received SIGUSR1, reloading vault token")
		loadVaultToken(api)
	}()
}

func wait(name string, retries int, test func() bool) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for retries > 0 {
		if test() {
			return nil
		}
		log.Printf("%s is not ready, waiting for another 5s (%d retries left)", name, retries)
		retries--
		<-ticker.C
	}
	return errors.New("retries expired")
}

func consulWaiter(api *consul.Client) func() bool {
	return func() bool {
		resp, err := api.Status().Leader()
		return err == nil && resp != ""
	}
}

func discoverVaultAddr(client *consul.Client) string {
	if addr := os.Getenv("VAULT_ADDR"); addr != "" {
		return addr
	}
	opt := &consul.QueryOptions{}
	services, _, err := client.Health().Service("vault", "active", true, opt)
	if err != nil {
		panic(err)
	}
	for _, service := range services {
		return fmt.Sprintf("http://%s:%d", service.Service.Address, service.Service.Port)
	}
	return ""
}
func DefaultClients() (*consul.Client, *vault.Client, error) {
	consulConfig := consul.DefaultConfig()
	consulAPI, err := consul.NewClient(consulConfig)
	if err != nil {
		return nil, nil, err
	}
	if wait("consul", 5, consulWaiter(consulAPI)) != nil {
		return nil, nil, errors.New("unable to connect to consul")
	}
	config := vault.DefaultConfig()
	config.Address = discoverVaultAddr(consulAPI)

	vaultAPI, err := vault.NewClient(config)
	if err != nil {
		return nil, nil, err
	}
	loadVaultToken(vaultAPI)
	return consulAPI, vaultAPI, nil
}
