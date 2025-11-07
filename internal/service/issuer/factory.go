// internal/service/issuer/factory.go
package issuer

import (
	"crypto/ed25519"

	"github.com/harishmurkal/6g-digi-wallet/internal/models"
	"github.com/harishmurkal/6g-digi-wallet/internal/service/crypto6g"
	"github.com/harishmurkal/6g-digi-wallet/internal/storage"
)

type IssuerService interface {
	GenerateDID(method string, opts map[string]any) (*models.DIDDocument, error)
	ResolveDID(did string) (*models.DIDDocument, error)
	ListDID() ([]*models.DIDDocument, error)
	CreateVC(req *models.VCRequest) (*models.VerifiableCredential, error)
}

// issuerService is the concrete implementation of the IssuerService interface.
type issuerService struct {
	store      storage.Store
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
	cryptoSvc  crypto6g.CryptoService
}

// NewIssuerService creates and returns a new IssuerService instance.
func NewIssuerService(store storage.Store, cSvc crypto6g.CryptoService) *issuerService {
	pub, priv, _ := ed25519.GenerateKey(nil)
	return &issuerService{
		store:      store,
		privateKey: priv,
		publicKey:  pub,
		cryptoSvc:  cSvc,
	}
}
