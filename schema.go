package config

type Schema struct {
	Cloudflare CloudflareSchema
	TLS        TLSSchema
}
type TLSSchema struct {
	CN string
}
type CloudflareSchema struct {
	APIToken     string `json:"api_token"`
	EmailAddress string `json:"email"`
}
