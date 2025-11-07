package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/harishmurkal/6g-digi-wallet/internal/api/handlers"
	"github.com/harishmurkal/6g-digi-wallet/internal/service/crypto6g"
	"github.com/harishmurkal/6g-digi-wallet/internal/service/issuer"
	"github.com/harishmurkal/6g-digi-wallet/internal/service/verifier"
	"github.com/harishmurkal/6g-digi-wallet/internal/service/wallet"
)

// logMiddleware logs requests and responses with method, path, and timing.
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("→ %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("← %s %s (%v)", r.Method, r.URL.Path, time.Since(start))
	})
}

// NewRouter constructs and returns a configured router.
func NewRouter(
	issuerSvc issuer.IssuerService,
	walletdidSvc wallet.DIDService,
	walletvcSvc wallet.VCService,
	walletvpSvc wallet.VPService,
	verifierSvc verifier.VerifierService,
	cryptoSvc crypto6g.CryptoService,
) *mux.Router {
	r := mux.NewRouter()

	// Initialize handlers
	issuerHandler := handlers.NewIssuerHandler(issuerSvc)
	walletHandler := handlers.NewWalletHandler(walletdidSvc, walletvcSvc, walletvpSvc)
	verifierHandler := handlers.NewVerifierHandler(verifierSvc)

	// ==== ISSUER ROUTES ====
	r.HandleFunc("/issuer/did/generate", issuerHandler.GenerateDID).Methods("POST")
	r.HandleFunc("/issuer/did/{id:.+}", issuerHandler.ResolveDID).Methods("GET")
	r.HandleFunc("/issuer/vc/create", issuerHandler.CreateVC).Methods("POST")

	// ==== WALLET ROUTES ====
	r.HandleFunc("/wallet/help", walletHandler.Help).Methods("GET")

	r.HandleFunc("/wallet/did/store", walletHandler.StoreDID).Methods("POST")
	r.HandleFunc("/wallet/did/list", walletHandler.ListDID).Methods("GET")
	r.HandleFunc("/wallet/did/{id:.+}", walletHandler.GetDID).Methods("GET")

	r.HandleFunc("/wallet/vc/store", walletHandler.StoreVC).Methods("POST")
	r.HandleFunc("/wallet/vc/list", walletHandler.ListVC).Methods("GET")
	r.HandleFunc("/wallet/vc/{id:.+}", walletHandler.GetVC).Methods("GET")
	r.HandleFunc("/wallet/vc/verify", walletHandler.VerifyVC).Methods("POST")

	r.HandleFunc("/wallet/vp/store", walletHandler.StoreVP).Methods("POST")
	r.HandleFunc("/wallet/vp/list", walletHandler.ListVP).Methods("GET")
	r.HandleFunc("/wallet/vp/{id:.+}", walletHandler.GetVP).Methods("GET")
	r.HandleFunc("/wallet/vp/verify", walletHandler.VerifyVP).Methods("POST")

	r.HandleFunc("/wallet/vp/build", walletHandler.BuildVP).Methods("POST")

	// ==== VERIFIER ROUTES ====
	r.HandleFunc("/verifier/vp/verify", verifierHandler.Verify).Methods("POST")

	r.Use(mux.MiddlewareFunc(logMiddleware))
	return r
}
