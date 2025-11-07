package verifier

import (
	"github.com/harishmurkal/6g-digi-wallet/internal/service/crypto6g"
	"github.com/harishmurkal/6g-digi-wallet/internal/storage"
)

type verifierService struct {
	store     storage.Store
	cryptoSvc crypto6g.CryptoService
}

func NewVerifierService(store storage.Store, cSvc crypto6g.CryptoService) VerifierService {
	return &verifierService{store: store, cryptoSvc: cSvc}
}
