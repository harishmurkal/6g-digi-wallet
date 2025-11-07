// internal/models/did.go
package models

type DIDDocument struct {
	ID        string           `json:"id"`
	Context   []string         `json:"@context,omitempty"`
	PublicKey []PublicKeyEntry `json:"publicKey,omitempty"`
	// ... other standard fields
}

type PublicKeyEntry struct {
	ID           string         `json:"id"`
	Type         string         `json:"type"`
	Controller   string         `json:"controller"`
	PublicKeyJWK map[string]any `json:"publicKeyJwk,omitempty"`
	// etc
}
