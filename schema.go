package config

import "fmt"

type Schema struct {
	Cloudflare CloudflareSchema
	TLS        TLSSchema
}
type TLSSchema struct {
	CN string `json:"cn"`
}
type CloudflareSchema struct {
	APIToken     string `json:"api_token"`
	EmailAddress string `json:"email"`
}

func (*TLSSchema) Template() string {
	tpl := `{{ "TLS Configuration" | green | bold }}
  {{ "CN" | faint }}: {{ .CN }}`
	return fmt.Sprintf("%s\n", tpl)
}
func (*CloudflareSchema) Template() string {
	tpl := `{{ "Cloudflare Configuration" | green | bold }}
  {{ "Email" | faint }}: {{ .EmailAddress }}
  {{ "API Token" | faint }}: {{ .APIToken }}`
	return fmt.Sprintf("%s\n", tpl)
}
