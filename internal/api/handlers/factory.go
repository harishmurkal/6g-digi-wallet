// internal/api/handlers/factory.go
package handlers

import (
	"github.com/harishmurkal/6g-digi-wallet/internal/service/issuer"
	"github.com/harishmurkal/6g-digi-wallet/internal/service/verifier"
	"github.com/harishmurkal/6g-digi-wallet/internal/service/wallet"
)

func NewIssuerHandler(svc issuer.IssuerService) *IssuerHandler {
	return &IssuerHandler{IssuerService: svc}
}

func NewWalletHandler(didSvc wallet.DIDService, vcSvc wallet.VCService, vpSvc wallet.VPService) *WalletHandler {
	return &WalletHandler{WalletdidSvc: didSvc, WalletvcSvc: vcSvc, WalletvpSvc: vpSvc}
}

func NewVerifierHandler(svc verifier.VerifierService) *VerifierHandler {
	return &VerifierHandler{VerifierService: svc}
}
