package keystone

type SecretString struct {
	Masked   string `json:"masked"`
	Original string `json:"original"`
}
