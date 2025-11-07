// internal/service/issuer/issuer_service.go
package issuer

import (
	"crypto/ed25519"
	"fmt"
	"log"
	"maps"
	"time"

	"github.com/google/uuid"
	"github.com/harishmurkal/6g-digi-wallet/internal/models"
)

// resolvePrivateKey fetches the appropriate private key and verification method ID
// associated with the given signing DID.
func (s *issuerService) resolvePrivateKey(signingDID string) (ed25519.PrivateKey, string, error) {
	if signingDID == "" {
		return nil, "", fmt.Errorf("signingDID is empty")
	}

	verificationMethodID := signingDID + "#key-1"
	privateKeyStoreKey := "privatekey:" + verificationMethodID

	var rawKey []byte

	// Load the private key bytes from the store.
	// Note: s.store.Load populates the provided pointer.
	if err := s.store.Load(privateKeyStoreKey, &rawKey); err != nil {
		return nil, "", fmt.Errorf("private key not found for DID %s: %w", signingDID, err)
	}

	// Validate the key size before casting
	if len(rawKey) != ed25519.PrivateKeySize {
		return nil, "", fmt.Errorf("invalid private key size (%d) for DID %s", len(rawKey), signingDID)
	}

	// In a real-world implementation, decrypt here if keys are encrypted at rest.
	privateKey := ed25519.PrivateKey(rawKey)

	return privateKey, verificationMethodID, nil
}

// GenerateDID creates a new key pair, constructs the DID Document with the public key,
// securely stores the private key, and returns the public DID Document.
func (s *issuerService) GenerateDID(method string, opts map[string]any) (*models.DIDDocument, error) {
	// 1. Determine the DID ID and construct the base DID string
	customID, _ := opts["id"].(string)
	idPart := customID
	if idPart == "" {
		// Use a UUID suffix for uniqueness
		idPart = uuid.NewString()[:8]
	}
	// Example DID: did:telco:a1b2c3d4
	did := fmt.Sprintf("did:%s:%s", method, idPart)

	// 2. Generate Cryptographic Key Pair
	// We'll use the Ed25519 standard, which is common in VC/DID.
	keyType := "Ed25519VerificationKey2018"

	// Assuming s.cryptoService has a method to generate a key pair
	privateKey, publicKeyJWK, err := s.cryptoSvc.GenerateKeyPair(keyType)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %w", err)
	}

	// 3. Construct the DID Document
	verificationMethodID := did + "#key-1" // Standard way to reference a key in the doc

	doc := &models.DIDDocument{
		ID:      did,
		Context: []string{"https://www.w3.org/ns/did/v1"},
		// Embed the Public Key for Verification
		PublicKey: []models.PublicKeyEntry{
			{
				ID:           verificationMethodID,
				Type:         keyType,
				Controller:   did,
				PublicKeyJWK: publicKeyJWK, // The public key in JWK format
			},
		},
		// You would typically add an 'authentication' or 'assertionMethod' block
		// referencing the key ID here for full spec compliance.
	}

	// 4. Securely Store the Private Key (Crucial for the Issuer)
	// The Issuer needs the private key to sign VCs later.
	// We use a separate key to store the private key, often encrypted.
	privateKeyStoreKey := "privatekey:" + verificationMethodID
	if err := s.store.Save(privateKeyStoreKey, privateKey); err != nil {
		// Log error, but proceed with public doc storage for now (rollback might be better)
		log.Printf("Warning: Failed to store private key for %s: %v", did, err)
	}

	// 5. Store the Public DID Document (for resolution by Verifiers)
	// This is the public record.
	didStoreKey := did // Use a clear key for public doc storage
	if err := s.store.Save(didStoreKey, doc); err != nil {
		return nil, fmt.Errorf("failed to save public DID Document %s: %w", did, err)
	}

	return doc, nil
}

// ResolveDID retrieves a DID document from the store.
func (s *issuerService) ResolveDID(id string) (*models.DIDDocument, error) {
	var doc models.DIDDocument
	if err := s.store.Load(id, &doc); err != nil {
		return nil, fmt.Errorf("DID not found: %s", id)
	}
	return &doc, nil
}

// ListDID returns all DIDs currently stored.
func (s *issuerService) ListDID() ([]*models.DIDDocument, error) {
	keys, err := s.store.ListKeys("did:")
	if err != nil {
		return nil, err
	}

	var dids []*models.DIDDocument
	for _, k := range keys {
		var doc models.DIDDocument
		if err := s.store.Load(k, &doc); err == nil {
			dids = append(dids, &doc)
		}
	}
	return dids, nil
}

// CreateVC constructs and signs a Verifiable Credential from a VCRequest.
func (s *issuerService) CreateVC(req *models.VCRequest) (*models.VerifiableCredential, error) {
	if req == nil {
		return nil, fmt.Errorf("nil VCRequest")
	}
	if req.IssuerDID == "" {
		return nil, fmt.Errorf("issuerDID missing")
	}
	if req.SubjectDID == "" {
		return nil, fmt.Errorf("subjectDID missing")
	}

	// 1. Prepare VC metadata and timestamps
	issuanceTime := time.Now().UTC().Round(time.Second)

	var expiryTime *time.Time
	if req.ValidityDays > 0 {
		tmp := issuanceTime.Add(time.Duration(req.ValidityDays) * 24 * time.Hour)
		expiryTime = &tmp
	}

	// Determine unique VC ID
	idPart, _ := req.Options["id"].(string)
	if idPart == "" {
		idPart = uuid.NewString()
	}
	vcID := fmt.Sprintf("vc:%s:%s", req.SubjectDID, idPart)

	// 2. Construct unsigned VC
	vc := &models.VerifiableCredential{
		Context:        []string{"https://www.w3.org/2018/credentials/v1"},
		ID:             vcID,
		Type:           append([]string{"VerifiableCredential"}, req.CredentialType...),
		Issuer:         req.IssuerDID,
		IssuanceDate:   issuanceTime,
		ExpirationDate: expiryTime,
		CredentialSubject: map[string]any{
			"id": req.SubjectDID,
		},
	}

	// Add user-defined claims into CredentialSubject
	maps.Copy(vc.CredentialSubject, req.Claims)

	// 3. Retrieve signing key (issuer’s private key + method ID)
	privateKey, verificationMethodID, err := s.resolvePrivateKey(req.IssuerDID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve signing key: %w", err)
	}

	// 4. Cryptographically sign the VC via crypto6g service
	signatureJWS, err := s.cryptoSvc.SignVC(vc, ed25519.PrivateKey(privateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to sign VC: %w", err)
	}

	// 5. Attach Proof (Linked Data Proof / JWS)
	vc.Proof = &models.Proof{
		Type:               "JsonWebSignature2020",
		Created:            issuanceTime,
		ProofPurpose:       "assertionMethod",
		VerificationMethod: verificationMethodID,
		JWS:                signatureJWS,
	}

	// 6. Save VC to the store (if persistence is enabled)
	if s.store != nil {
		if err := s.store.Save(vc.ID, vc); err != nil {
			return nil, fmt.Errorf("failed to save VC: %w", err)
		}
	}

	return vc, nil
}
