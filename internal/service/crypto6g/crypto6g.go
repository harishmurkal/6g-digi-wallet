package crypto6g

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/harishmurkal/6g-digi-wallet/internal/models"
)

// CryptoService defines the methods required for key generation and digital signatures.
type CryptoService interface {
	// GenerateKeyPair creates a public/private key pair (e.g., Ed25519) and returns the private key
	// (as raw bytes) and the public key (in JWK format) for use in a DID Document.
	GenerateKeyPair(keyType string) ([]byte, map[string]any, error)

	// SignVC generates a standardized signature (e.g., JWS) over a Verifiable Credential payload.
	SignVC(vc *models.VerifiableCredential, privateKey ed25519.PrivateKey) (string, error)

	// VerifySignature verifies the cryptographic proof on a Verifiable Presentation (VP) or a VC.
	// The payload is the data being verified, which holds the Proof field.
	VerifySignature(proof *models.Proof, payload any) (bool, error)
}

// cryptoService is a concrete implementation using standard Go libraries.
type cryptoService struct {
	// Dependencies can be added here if needed (e.g., key vault reference)
}

// NewCryptoService creates a new instance of the crypto service.
func NewCryptoService() CryptoService {
	return &cryptoService{}
}

// --- Implementation of CryptoService Interface ---

// GenerateKeyPair generates an Ed25519 key pair.
func (s *cryptoService) GenerateKeyPair(keyType string) ([]byte, map[string]any, error) {
	if keyType != "Ed25519VerificationKey2018" && keyType != "Ed25519VerificationKey2020" {
		return nil, nil, fmt.Errorf("unsupported key type: %s", keyType)
	}

	// Generate the Ed25519 key pair
	publicKey, privateKey, err := ed25519.GenerateKey(rand.New(rand.NewSource(time.Now().UnixNano())))
	if err != nil {
		return nil, nil, err
	}

	// Convert Public Key to JWK Format (required for DID Document)
	publicKeyJWK := map[string]any{
		"kty": "OKP", // Octet Key Pair
		"crv": "Ed25519",
		"x":   base64.RawURLEncoding.EncodeToString(publicKey),
	}

	// Return the raw private key bytes and the public JWK
	return privateKey, publicKeyJWK, nil
}

// SignVC is a placeholder for the complex, standards-compliant signing process.
func (s *cryptoService) SignVC(vc *models.VerifiableCredential, privateKey ed25519.PrivateKey) (string, error) {
	// 1. Remove the 'Proof' field from the VC (if it exists).
	vcToSign := *vc
	vcToSign.Proof = nil

	// 2. Canonicalize the VC JSON structure (MUST use JCS or similar canonicalization).
	// STUB: Using standard Go marshal is NOT secure/compliant but serves as a placeholder.
	canonicalVC, err := json.Marshal(vcToSign)
	if err != nil {
		return "", err
	}

	// 3. Sign the canonicalized bytes.
	signature := ed25519.Sign(privateKey, canonicalVC)

	// 4. Encode the signature (often JWS or LDP)
	return base64.RawURLEncoding.EncodeToString(signature), nil
}

// VerifySignature is a placeholder for verifying the signature against the payload.
func (s *cryptoService) VerifySignature(proof *models.Proof, payload any) (bool, error) {
	// 1. Retrieve the appropriate Public Key based on proof.VerificationMethod (Requires DID Resolution).
	// 2. Canonicalize the payload (same method used in signing, removing the proof).
	// 3. Verify the proof.SignatureValue (or proof.JWS) against the canonicalized payload using the public key.

	// STUB: Always return true for development purposes.
	return true, nil
}
