package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	consul "github.com/hashicorp/consul/api"
	vault "github.com/hashicorp/vault/api"
)

const (
	ConfigPrefix = "mqtt/conf"
)

type opts struct {
	Notify chan struct{}
}
type getOpt func(opts) opts

func configKey(key string) string {
	return fmt.Sprintf("%s/%s", ConfigPrefix, key)
}

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
func watchKey(api *consul.Client, index uint64, key string) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		for {
			opts := &consul.QueryOptions{
				WaitIndex: index,
				WaitTime:  30 * time.Second,
			}
			_, meta, err := api.KV().Get(key, opts)
			if err != nil {
				log.Printf("WARN: exiting watch on %s: %v", key, err)
				return
			}
			if meta.LastIndex > index {
				return
			}
		}
	}()
	return ch
}
func TLS(api *consul.Client) (TLSSchema, <-chan struct{}, error) {
	opts := &consul.QueryOptions{}
	out := TLSSchema{}
	key := configKey("tls")
	pair, meta, err := api.KV().Get(key, opts)
	if err != nil {
		return out, nil, err
	}
	if pair == nil {
		return out, nil, errors.New("key not found")
	}
	err = json.Unmarshal(pair.Value, &out)
	if err != nil {
		return out, nil, err
	}
	ch := watchKey(api, meta.LastIndex, key)
	return out, ch, nil
}
