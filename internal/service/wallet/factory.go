package wallet

import (
	"github.com/harishmurkal/6g-digi-wallet/internal/models"
	"github.com/harishmurkal/6g-digi-wallet/internal/service/crypto6g"
	"github.com/harishmurkal/6g-digi-wallet/internal/storage"
)

// ---- DID Service Interface ----
type DIDService interface {
	StoreDID(doc *models.DIDDocument) error
	GetDID(id string) (*models.DIDDocument, error)
	ListDIDs() ([]*models.DIDDocument, error)
}

// ---- VC Service Interface ----
type VCService interface {
	StoreVC(vc *models.VerifiableCredential) error
	GetVC(id string) (*models.VerifiableCredential, error)
	ListVCs(filter models.VCFilter) ([]*models.VerifiableCredential, error)
	VerifyVC(vc *models.VerifiableCredential) (bool, error)
	BuildVP(vcIDs []string, revealFields map[string][]string, nonce string) (*models.VerifiablePresentation, error)
}

// ---- VP Service Interface ----
type VPService interface {
	StoreVP(vc *models.VerifiablePresentation) error
	GetVP(id string) (*models.VerifiablePresentation, error)
	ListVPs(filter models.VPFilter) ([]*models.VerifiablePresentation, error)
	VerifyVP(vc *models.VerifiablePresentation) (bool, error)
}

// ---- Combined WalletService struct (implements both) ----
type WalletService struct {
	store     storage.Store
	cryptoSvc crypto6g.CryptoService
}

// Constructors
func NewDIDService(store storage.Store, cSvc crypto6g.CryptoService) DIDService {
	return &WalletService{store: store, cryptoSvc: cSvc}
}

func NewVCService(store storage.Store, cSvc crypto6g.CryptoService) VCService {
	return &WalletService{store: store, cryptoSvc: cSvc}
}

func NewVPService(store storage.Store, cSvc crypto6g.CryptoService) VPService {
	return &WalletService{store: store, cryptoSvc: cSvc}
}
