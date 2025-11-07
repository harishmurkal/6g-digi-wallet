package models

import "time"

// VerifiablePresentation follows W3C VP Data Model v1.1
// Unique Key or ID (i.e. Holder) would be vp+Nonce
type VerifiablePresentation struct {
	Context              []string                `json:"@context"`
	Type                 []string                `json:"type"`
	VerifiableCredential []*VerifiableCredential `json:"verifiableCredential"`
	Holder               string                  `json:"holder,omitempty"`
	Proof                *Proof                  `json:"proof,omitempty"`
	Nonce                string                  `json:"nonce,omitempty"`
	Created              time.Time               `json:"created"`
}
