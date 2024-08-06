package keystone

// VerifyString is a string that can be verified
type VerifyString struct {
	Original string `json:"original,omitempty"`
}

// String returns the original string if it exists, otherwise the masked string
func (e VerifyString) String() string {
	return e.Original
}

// NewVerifyString creates a new VerifyString
func NewVerifyString(original string) VerifyString {
	return VerifyString{Original: original}
}
