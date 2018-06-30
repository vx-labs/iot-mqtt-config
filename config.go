package config

import (
	vault "github.com/hashicorp/vault/api"
)

type opts struct {
	Notify chan struct{}
}
type getOpt func(opts) opts

func Cloudflare(api *vault.Client) (CloudflareSchema, error) {
	out := struct {
		Data CloudflareSchema `json:"data"`
	}{}
	req := api.NewRequest("GET", "/v1/secret/data/vx/cloudflare")
	resp, err := api.RawRequest(req)
	if err != nil {
		return out.Data, err
	}
	defer resp.Body.Close()
	return out.Data, resp.DecodeJSON(&out)
}
