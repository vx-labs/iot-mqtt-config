package config

import "fmt"

type Schema struct {
	Cloudflare CloudflareSchema
	TLS        TLSSchema
	HTTP       HTTPSchema
}
type HTTPSchema struct {
	Proxy string `json:"proxy"`
}
type TLSSchema struct {
	CN                      string `json:"cn"`
	LetsEncryptAccountEmail string `json:"le_email"`
}
type CloudflareSchema struct {
	APIToken     string `json:"api_token"`
	EmailAddress string `json:"email"`
}

type AuthenticationSchema struct {
	StaticTokens []string `json:"static_tokens"`
}

func (*AuthenticationSchema) Template() string {
	tpl := `{{ "Authentication Configuration" | green | bold }}
  {{ "Static tokens" | faint }}: {{ .StaticTokens }}`
	return fmt.Sprintf("%s\n", tpl)
}
func (*TLSSchema) Template() string {
	tpl := `{{ "TLS Configuration" | green | bold }}
  {{ "Letsencrypt account email" | faint }}: {{ .LetsEncryptAccountEmail }}
  {{ "CN" | faint }}: {{ .CN }}`
	return fmt.Sprintf("%s\n", tpl)
}
func (*CloudflareSchema) Template() string {
	tpl := `{{ "Cloudflare Configuration" | green | bold }}
  {{ "Email" | faint }}: {{ .EmailAddress }}
  {{ "API Token" | faint }}: {{ .APIToken }}`
	return fmt.Sprintf("%s\n", tpl)
}
func (*HTTPSchema) Template() string {
	tpl := `{{ "HTTP Configuration" | green | bold }}
  {{ "Proxy" | faint }}: {{ .Proxy }}`
	return fmt.Sprintf("%s\n", tpl)
}
