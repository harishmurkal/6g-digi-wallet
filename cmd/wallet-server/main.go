package main

import (
	"log"
	"net/http"
	"os"

	"github.com/harishmurkal/6g-digi-wallet/internal/api"
	"github.com/harishmurkal/6g-digi-wallet/internal/service/crypto6g"
	"github.com/harishmurkal/6g-digi-wallet/internal/service/issuer"
	"github.com/harishmurkal/6g-digi-wallet/internal/service/verifier"
	"github.com/harishmurkal/6g-digi-wallet/internal/service/wallet"
	"github.com/harishmurkal/6g-digi-wallet/internal/storage"
)

func main() {
	// 1️⃣ Detect backend from env
	backendEnv := os.Getenv("STORE_BACKEND")
	if backendEnv == "" {
		backendEnv = string(storage.BackendMemory)
	}

	// 2️⃣ Default options map
	opts := map[string]string{
		"path": os.Getenv("STORE_PATH"),
	}
	if opts["path"] == "" {
		opts["path"] = "./data/wallet_store.jsonl"
	}

	// 3️⃣ Initialize store
	store, err := storage.NewStore(storage.BackendType(backendEnv), opts)
	if err != nil {
		log.Fatalf("❌ Failed to initialize store (backend=%s): %v", backendEnv, err)
	}

	log.Printf("🗄️  Using storage backend: %s (path=%s)", backendEnv, opts["path"])

	crypto := crypto6g.NewCryptoService() // or whatever your constructor is

	// 4️⃣ Create all services sharing the same store
	issuerSvc := issuer.NewIssuerService(store, crypto)
	didSvc := wallet.NewDIDService(store, crypto)
	vcSvc := wallet.NewVCService(store, crypto)
	vpSvc := wallet.NewVPService(store, crypto)
	verifierSvc := verifier.NewVerifierService(store, crypto)

	// 5️⃣ Initialize API router
	r := api.NewRouter(issuerSvc, didSvc, vcSvc, vpSvc, verifierSvc, crypto)

	// 6️⃣ Start HTTP server
	log.Println("🚀 Wallet server running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
