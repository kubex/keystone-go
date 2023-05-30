package keystone

type SecretString struct {
	Masked   string `json:"masked"`
	Original string `json:"original"`
}

func (e SecretString) String() string {
	if e.Original != "" {
		return e.Original
	}
	return e.Masked
}

func NewSecretString(original, masked string) SecretString {
	return SecretString{
		Masked:   masked,
		Original: original,
	}
}
