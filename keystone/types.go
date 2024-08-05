package keystone

// Amount represents money
type Amount struct {
	Currency string `json:"currency"`
	Units    int64  `json:"units"`
}

// NewAmount creates a new Amount
func NewAmount(currency string, units int64) Amount {
	return Amount{
		Currency: currency,
		Units:    units,
	}
}

// SecretString is a string that represents sensitive Data
type SecretString struct {
	Masked   string `json:"masked,omitempty"`
	Original string `json:"original,omitempty"`
}

// String returns the original string if it exists, otherwise the masked string
func (e SecretString) String() string {
	if e.Original != "" {
		return e.Original
	}
	return e.Masked
}

// NewSecretString creates a new SecretString
func NewSecretString(original, masked string) SecretString {
	return SecretString{
		Masked:   masked,
		Original: original,
	}
}
