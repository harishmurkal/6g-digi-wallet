package wallet

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/harishmurkal/6g-digi-wallet/internal/models"
)

// ----------------------
// DID Service Functions
// ----------------------

func (s *WalletService) StoreDID(doc *models.DIDDocument) error {
	if doc.ID == "" {
		return errors.New("DID must have an ID")
	}
	return s.store.Save(doc.ID, doc)
}

func (s *WalletService) GetDID(id string) (*models.DIDDocument, error) {
	if id == "" {
		return nil, fmt.Errorf("empty DID ID")
	}

	var doc models.DIDDocument
	if err := s.store.Load(id, &doc); err != nil {
		return nil, fmt.Errorf("failed to load DID: %w", err)
	}

	return &doc, nil
}

func (s *WalletService) ListDIDs() ([]*models.DIDDocument, error) {
	keys, err := s.store.ListKeys("")
	if err != nil {
		return nil, err
	}

	var dids []*models.DIDDocument
	for _, key := range keys {
		if !strings.HasPrefix(key, "did:") {
			continue
		}
		var doc models.DIDDocument
		if err := s.store.Load(key, &doc); err == nil {
			dids = append(dids, &doc)
		}
	}

	return dids, nil
}

// ----------------------
// VC Service Functions
// ----------------------

func (s *WalletService) StoreVC(vc *models.VerifiableCredential) error {
	if vc.ID == "" {
		return errors.New("VC must have an ID")
	}
	return s.store.Save(vc.ID, vc)
}

func (s *WalletService) GetVC(id string) (*models.VerifiableCredential, error) {
	if id == "" {
		return nil, fmt.Errorf("empty VC ID")
	}

	var vc models.VerifiableCredential
	if err := s.store.Load(id, &vc); err != nil {
		return nil, fmt.Errorf("failed to load VC: %w", err)
	}

	return &vc, nil
}

func (s *WalletService) ListVCs(filter models.VCFilter) ([]*models.VerifiableCredential, error) {
	keys, err := s.store.ListKeys("")
	if err != nil {
		return nil, err
	}

	var vcs []*models.VerifiableCredential
	for _, key := range keys {
		if !strings.HasPrefix(key, "vc:") {
			continue
		}

		var vc models.VerifiableCredential
		if err := s.store.Load(key, &vc); err == nil {
			if matchVCFilter(&vc, filter) {
				vcs = append(vcs, &vc)
			}
		}
	}

	return vcs, nil
}

func (s *WalletService) VerifyVC(vc *models.VerifiableCredential) (bool, error) {
	if vc.Proof == nil {
		return false, errors.New("missing proof")
	}
	return true, nil
}

// Helper: filter function
func matchVCFilter(vc *models.VerifiableCredential, f models.VCFilter) bool {
	if f.Issuer != "" && vc.Issuer != f.Issuer {
		return false
	}
	if f.CredType != "" {
		found := slices.Contains(vc.Type, f.CredType)
		if !found {
			return false
		}
	}
	if f.ActiveOnly && vc.ExpirationDate != nil && vc.ExpirationDate.Before(time.Now()) {
		return false
	}
	if f.ExpiredOnly && (vc.ExpirationDate == nil || vc.ExpirationDate.After(time.Now())) {
		return false
	}
	return true
}

// ----------------------
// VP Service Functions
// ----------------------

func (s *WalletService) StoreVP(vp *models.VerifiablePresentation) error {
	if vp.Holder == "" {
		return errors.New("VP must have an Holder")
	}
	// Store with a "vp:" prefix to distinguish from VCs
	return s.store.Save("vp:"+vp.Holder, vp)
}

func (s *WalletService) GetVP(id string) (*models.VerifiablePresentation, error) {
	if id == "" {
		return nil, fmt.Errorf("empty VP ID")
	}

	var vp models.VerifiablePresentation
	if err := s.store.Load("vp:"+id, &vp); err != nil {
		return nil, fmt.Errorf("failed to load VP: %w", err)
	}

	return &vp, nil
}

func (s *WalletService) ListVPs(filter models.VPFilter) ([]*models.VerifiablePresentation, error) {
	keys, err := s.store.ListKeys("")
	if err != nil {
		return nil, err
	}

	var vps []*models.VerifiablePresentation
	for _, key := range keys {
		// Ensure we only load VPs
		if !strings.HasPrefix(key, "vp:") {
			continue
		}

		var vp models.VerifiablePresentation
		if err := s.store.Load(key, &vp); err == nil {
			// You'll need to implement this helper function
			if matchVPFilter(&vp, filter) {
				vps = append(vps, &vp)
			}
		}
	}

	return vps, nil
}

// Verifies the VP (e.g., checks its proof/signature)
// NOTE: This is a placeholder. A real implementation would involve
// cryptographic verification of the vp.Proof.
func (s *WalletService) VerifyVP(vp *models.VerifiablePresentation) (bool, error) {
	if vp.Proof == nil {
		return false, errors.New("missing proof")
	}

	// TODO: Add actual crypto verification logic here.
	// For example:
	// - Get the signer's public key (from vp.Holder or vp.Proof.VerificationMethod)
	// - Re-create the unsigned VP
	// - Verify the signature (vp.Proof.ProofValue) against the unsigned VP and public key
	//logInfo("Placeholder: VP %s proof verification successful", vp.ID)
	return true, nil
}

// matchVPFilter is a helper to filter VPs based on criteria.
// This is a placeholder implementation.
func matchVPFilter(vp *models.VerifiablePresentation, filter models.VPFilter) bool {
	// TODO: Implement actual filtering logic based on filter fields.
	// Example (if filter had a 'Domain' field):
	// if filter.Domain != "" && vp.Proof.Domain != filter.Domain {
	// 	 return false
	// }
	return true // Default: match all
}

// ----------------------
// VP Builder
// ----------------------

func (s *WalletService) BuildVP(vcIDs []string, revealFields map[string][]string, nonce string) (*models.VerifiablePresentation, error) {
	if len(vcIDs) == 0 {
		return nil, errors.New("no VC IDs provided")
	}
	if nonce == "" {
		return nil, errors.New("nonce is required for Verifiable Presentation")
	}

	var disclosedVCs []*models.VerifiableCredential

	// 1. Load VCs and Apply Selective Disclosure
	for _, id := range vcIDs {
		// IMPORTANT FIX: You must load the VC using the storage key prefix (e.g., "vc:")
		// Assuming your VCs are stored with the key "vc:" + vc.ID
		storeKey := id

		var vc models.VerifiableCredential
		if err := s.store.Load(storeKey, &vc); err != nil {
			return nil, fmt.Errorf("failed to load VC %s (key: %s): %w", id, storeKey, err)
		}

		// --- Selective Disclosure Logic ---
		// In a real wallet, you would clone the VC and remove/hash fields not requested
		// in the 'revealFields' map to create a Disclosed Credential.
		// For this example, we'll just return the full VC but acknowledge the step.

		// NOTE: A robust implementation requires JSON-LD Frame or BBS+ proof generation.
		// Since we're stubbing, we'll just add the full VC.

		disclosedVCs = append(disclosedVCs, &vc)
	}

	// 2. Define the Holder's DID (The wallet's DID)
	// You need to know the DID of the entity presenting the VP. This is crucial.
	// We'll assume the wallet has a default DID stored somewhere.
	// For now, let's use a placeholder.
	holderDID := "vp:" + nonce

	// 3. Construct the VP
	vp := &models.VerifiablePresentation{
		// Context: Use standard W3C context and ensure you include the VCs' contexts
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
			// Add any other necessary contexts (e.g., specific VC type contexts)
		},
		Type:                 []string{"VerifiablePresentation"},
		VerifiableCredential: disclosedVCs,
		Holder:               holderDID,
		Nonce:                nonce,
		Created:              time.Now().UTC().Round(time.Second), // Use UTC and round for consistency
		// Set a unique ID for the VP using the Holder DID and Nonce
		// This is often *not* an official VP field, but useful for storage/lookup
		// FIX: Use a unique ID for storage, which is separate from the VP's 'holder' field
		// We'll use this for the storage key, not the VP JSON body itself.
	}

	// 4. Cryptographically Sign the VP
	// This is the most critical part of the BuildVP function.
	// The wallet signs the VP using the Holder's private key.

	// FIX: The Proof must be calculated over the final VP structure.
	// The 'SignatureValue' is currently a stub. In a real system, you would:
	// 1. Canonicalize and serialize the VP structure (excluding the Proof field)
	// 2. Sign the resulting bytes with the holder's private key
	// 3. Encode the signature (e.g., as a JWS or base64 string)

	vp.Proof = &models.Proof{
		Type:               "Ed25519Signature2018", // Or use JWS
		Created:            time.Now().UTC().Round(time.Second),
		ProofPurpose:       "authentication",                         // Common purpose for a VP
		VerificationMethod: holderDID + "#key-1",                     // Must reference a key in the Holder's DID Document
		SignatureValue:     "REAL-CRYPTOGRAPHIC-SIGNATURE-OF-THE-VP", // STUB - Replace with real signing logic
	}

	// 5. Save the final VP in the wallet
	// FIX: Use a proper storage key for the VP, like "vp:" + nonce
	vpStoreKey := "vp:" + nonce
	if err := s.store.Save(vpStoreKey, vp); err != nil {
		return nil, fmt.Errorf("failed to save VP %s: %w", vpStoreKey, err)
	}

	return vp, nil
}
