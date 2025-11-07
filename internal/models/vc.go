package models

import "time"

// VCRequest represents an input payload to create a Verifiable Credential.
type VCRequest struct {
	IssuerDID      string         `json:"issuerDID"`
	SubjectDID     string         `json:"subjectDID"`
	CredentialType []string       `json:"credentialType"`
	Claims         map[string]any `json:"claims"`
	ValidityDays   int            `json:"validityDays"`
	Options        map[string]any `json:"options"`
}

// VerifiableCredential follows W3C VC Data Model v1.1
type VerifiableCredential struct {
	Context           []string          `json:"@context"`
	ID                string            `json:"id"`
	Type              []string          `json:"type"`
	Issuer            string            `json:"issuer"`
	IssuanceDate      time.Time         `json:"issuanceDate"`
	ExpirationDate    *time.Time        `json:"expirationDate,omitempty"`
	CredentialSubject map[string]any    `json:"credentialSubject"`
	CredentialStatus  *CredentialStatus `json:"credentialStatus,omitempty"`
	Proof             *Proof            `json:"proof,omitempty"`
}

// CredentialStatus allows revocation / suspension tracking
type CredentialStatus struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// Proof represents the cryptographic signature over the credential
type Proof struct {
	Type               string    `json:"type"`
	Created            time.Time `json:"created"`
	ProofPurpose       string    `json:"proofPurpose"`
	VerificationMethod string    `json:"verificationMethod"`
	JWS                string    `json:"jws,omitempty"`
	SignatureValue     string    `json:"signatureValue,omitempty"`
}

type VCFilter struct {
	Issuer      string `json:"issuer,omitempty"`
	SubjectID   string `json:"subjectId,omitempty"`
	CredType    string `json:"credType,omitempty"`
	ActiveOnly  bool   `json:"activeOnly,omitempty"`
	ExpiredOnly bool   `json:"expiredOnly,omitempty"`
}

type VPFilter struct {
	Issuer      string `json:"issuer,omitempty"`
	SubjectID   string `json:"subjectId,omitempty"`
	CredType    string `json:"credType,omitempty"`
	ActiveOnly  bool   `json:"activeOnly,omitempty"`
	ExpiredOnly bool   `json:"expiredOnly,omitempty"`
}
