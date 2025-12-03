package crypto6g

import (
	"crypto/ed25519"
	"encoding/base64"
	"testing"
	"time"

	"github.com/harishmurkal/6g-digi-wallet/internal/models"
)

// TestGenerateKeyPair_Valid ensures Ed25519 key generation works correctly.
func TestGenerateKeyPair_Valid(t *testing.T) {
	svc := NewCryptoService()

	privateKey, publicKeyJWK, err := svc.GenerateKeyPair("Ed25519VerificationKey2018")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if privateKey == nil || len(privateKey) != ed25519.PrivateKeySize {
		t.Errorf("invalid private key length: got %d", len(privateKey))
	}

	if publicKeyJWK["kty"] != "OKP" || publicKeyJWK["crv"] != "Ed25519" {
		t.Errorf("invalid public key JWK fields: %+v", publicKeyJWK)
	}

	if _, ok := publicKeyJWK["x"]; !ok {
		t.Error("public key JWK missing 'x' field")
	}
}

// TestGenerateKeyPair_Invalid ensures invalid key types are rejected.
func TestGenerateKeyPair_Invalid(t *testing.T) {
	svc := NewCryptoService()

	_, _, err := svc.GenerateKeyPair("RSA")
	if err == nil {
		t.Error("expected error for unsupported key type, got nil")
	}
}

// TestSignVC ensures a Verifiable Credential can be signed correctly.
func TestSignVC(t *testing.T) {
	svc := NewCryptoService()

	// Generate keypair first
	priv, pubJWK, err := svc.GenerateKeyPair("Ed25519VerificationKey2018")
	if err != nil {
		t.Fatalf("failed to generate keypair: %v", err)
	}

	// Create a mock Verifiable Credential
	issuanceDate, _ := time.Parse(time.RFC3339, "2025-11-07T12:00:00Z")
	vc := &models.VerifiableCredential{
		ID:           "urn:uuid:1234",
		Context:      []string{"https://www.w3.org/2018/credentials/v1"},
		Type:         []string{"VerifiableCredential"},
		Issuer:       "did:example:issuer",
		IssuanceDate: issuanceDate,
		CredentialSubject: map[string]any{
			"id":   "did:example:subject",
			"name": "Alice",
		},
	}

	signature, err := svc.SignVC(vc, ed25519.PrivateKey(priv))
	if err != nil {
		t.Fatalf("SignVC failed: %v", err)
	}

	if _, err := base64.RawURLEncoding.DecodeString(signature); err != nil {
		t.Errorf("signature is not valid base64: %v", err)
	}

	if len(pubJWK) == 0 {
		t.Error("expected non-empty public key JWK")
	}
}

// TestVerifySignature ensures the stub returns expected values.
func TestVerifySignature(t *testing.T) {
	svc := NewCryptoService()

	ok, err := svc.VerifySignature(&models.Proof{}, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !ok {
		t.Error("expected true from VerifySignature stub")
	}
}
