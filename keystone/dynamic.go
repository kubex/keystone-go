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

type Amount struct {
	Currency string `json:"currency"`
	Units    int64  `json:"units"`
}

func NewAmount(currency string, units int64) Amount {
	return Amount{
		Currency: currency,
		Units:    units,
	}
}
