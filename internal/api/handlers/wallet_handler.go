package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harishmurkal/6g-digi-wallet/internal/models"
	"github.com/harishmurkal/6g-digi-wallet/internal/service/wallet"
)

type WalletHandler struct {
	WalletdidSvc wallet.DIDService
	WalletvcSvc  wallet.VCService
	WalletvpSvc  wallet.VPService
}

// GET /wallet/help
func (h *WalletHandler) Help(w http.ResponseWriter, r *http.Request) {
	help := map[string]string{
		"/wallet/help":        "Show this help message",
		"/wallet/did/store":   "POST: Store a DID Document (body: DIDDocument)",
		"/wallet/did/{id}":    "GET: Fetch DID Document by ID",
		"/wallet/did/list":    "GET: List all stored DIDs",
		"/wallet/vc/store":    "POST: Store a Verifiable Credential (body: VerifiableCredential)",
		"/wallet/vc/{id}":     "GET: Fetch VC by ID",
		"/wallet/vc/list":     "GET: List all stored VCs (optional filters later)",
		"/wallet/vp/build":    "POST: Build a Verifiable Presentation from VC IDs",
		"/wallet/verify":      "POST: Verify a VC by ID",
		"/verifier/vp/verify": "POST: Verify Verifiable Presentation (verifier side)",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(help)
}

// ---- DID Section ----

// POST /wallet/did/store
func (h *WalletHandler) StoreDID(w http.ResponseWriter, r *http.Request) {
	logInfo("WalletHandler.StoreDID called")
	var didDoc models.DIDDocument
	if err := json.NewDecoder(r.Body).Decode(&didDoc); err != nil {
		logError("invalid DID Doc: %v", err)
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.WalletdidSvc.StoreDID(&didDoc); err != nil {
		logError("Failed to store DID Doc: %v", err)
		http.Error(w, "failed to store DID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": "stored",
		"id":     didDoc.ID,
	})
	logInfo("WalletHandler.StoreDID responded successfully with DID: %s", didDoc.ID)
}

// GET /wallet/did/{id}
func (h *WalletHandler) GetDID(w http.ResponseWriter, r *http.Request) {
	logInfo("WalletHandler.GetDID called")
	id := mux.Vars(r)["id"]
	didDoc, err := h.WalletdidSvc.GetDID(id)
	if err != nil {
		logError("Failed to retrieve DID Doc: %v", err)
		http.Error(w, "DID not found: "+err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(didDoc)
	logInfo("WalletHandler.GetDID responded successfully with DID: %s", id)
}

// GET /wallet/did/list
func (h *WalletHandler) ListDID(w http.ResponseWriter, r *http.Request) {
	logInfo("WalletHandler.ListDID called")
	dids, err := h.WalletdidSvc.ListDIDs()
	if err != nil {
		logError("Failed to List DIDs: %v", err)
		http.Error(w, "failed to list DIDs: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Log the JSON string, not the Go slice
	didsJSON, jsonErr := json.MarshalIndent(dids, "", "  ")
	if jsonErr != nil {
		// Fallback in case marshaling fails
		logError("Failed to marshal DIDs for logging: %v", jsonErr)
		logInfo("WalletHandler.ListDID responded successfully with %d did", len(dids))
	} else {
		// This will print the full JSON array
		logInfo("WalletHandler.ListDID responded successfully with retrieved did: %s", string(didsJSON))
	}

	json.NewEncoder(w).Encode(dids)
	//logInfo("WalletHandler.ListDID responded successfully with retrieved DIDs: %v", dids)
}

// ---- VC Section ----

// POST /wallet/vc/store
func (h *WalletHandler) StoreVC(w http.ResponseWriter, r *http.Request) {
	logInfo("WalletHandler.StoreVC called")
	var vc models.VerifiableCredential
	if err := json.NewDecoder(r.Body).Decode(&vc); err != nil {
		logError("Invalid VC to store: %v", err)
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.WalletvcSvc.StoreVC(&vc); err != nil {
		logError("Failed to store VC: %v", err)
		http.Error(w, "failed to store VC: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": "stored",
		"id":     vc.ID,
	})
	logInfo("WalletHandler.StoreVC responded successfully with stored VC: %v", vc.ID)
}

// GET /wallet/vc/{id}
func (h *WalletHandler) GetVC(w http.ResponseWriter, r *http.Request) {
	logInfo("WalletHandler.GetVC called")
	id := mux.Vars(r)["id"]
	vc, err := h.WalletvcSvc.GetVC(id)
	if err != nil {
		logError("VC not found: %v", err)
		http.Error(w, "VC not found: "+err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(vc)
	logInfo("WalletHandler.GetVC responded successfully with retrieved vc id: %s", id)
}

// GET /wallet/vc/list
func (h *WalletHandler) ListVC(w http.ResponseWriter, r *http.Request) {
	logInfo("WalletHandler.ListVC called")
	filter := models.VCFilter{} // later: parse query params like ?issuer=...&activeOnly=true
	vcs, err := h.WalletvcSvc.ListVCs(filter)
	if err != nil {
		logError("Failed to List VCs: %v", err)
		http.Error(w, "failed to list VCs: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Log the JSON string, not the Go slice
	vcsJSON, jsonErr := json.MarshalIndent(vcs, "", "  ")
	if jsonErr != nil {
		// Fallback in case marshaling fails
		logError("Failed to marshal VCs for logging: %v", jsonErr)
		logInfo("WalletHandler.ListVC responded successfully with %d vcs", len(vcs))
	} else {
		// This will print the full JSON array
		logInfo("WalletHandler.ListVC responded successfully with retrieved vcs: %s", string(vcsJSON))
	}

	json.NewEncoder(w).Encode(vcs)
	//logInfo("WalletHandler.ListVC responded successfully with retrieved vcs: %v", vcs)
}

// POST /wallet/vc/verify
func (h *WalletHandler) VerifyVC(w http.ResponseWriter, r *http.Request) {
	logInfo("WalletHandler.VerifyVC called")
	var req struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logError("Invalid VC: %v", err)
		http.Error(w, "invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	vc, err := h.WalletvcSvc.GetVC(req.ID)
	if err != nil {
		logError("VC ID, %s not found: %v", req.ID, err)
		http.Error(w, "VC not found: "+err.Error(), http.StatusNotFound)
		return
	}

	valid, err := h.WalletvcSvc.VerifyVC(vc)
	if err != nil {
		logError("VC ID, %s verification failed: %v", req.ID, err)
		http.Error(w, "verification error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"id":     req.ID,
		"valid":  valid,
		"status": "verified",
	})
	logInfo("WalletHandler.VerifyVC responded successfully with verified vc: %v", req.ID)
}

// ---- VP Section ----

// POST /wallet/vp/store
func (h *WalletHandler) StoreVP(w http.ResponseWriter, r *http.Request) {
	logInfo("WalletHandler.StoreVP called")
	var vp models.VerifiablePresentation
	if err := json.NewDecoder(r.Body).Decode(&vp); err != nil {
		logError("Invalid VP to store: %v", err)
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Assuming a WalletvpSvc similar to your WalletvcSvc
	if err := h.WalletvpSvc.StoreVP(&vp); err != nil {
		logError("Failed to store VP: %v", err)
		http.Error(w, "failed to store VP: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": "stored",
		"id":     vp.Holder, // Assuming vp has an Holder field
	})
	logInfo("WalletHandler.StoreVP responded successfully with stored VP: %v", vp.Holder)
}

// GET /wallet/vp/{id:.+}
func (h *WalletHandler) GetVP(w http.ResponseWriter, r *http.Request) {
	logInfo("WalletHandler.GetVP called")
	id := mux.Vars(r)["id"] // Using {id:.+} allows IDs with slashes

	// Assuming a WalletvpSvc
	vp, err := h.WalletvpSvc.GetVP(id)
	if err != nil {
		logError("VP not found: %v", err)
		http.Error(w, "VP not found: "+err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(vp)
	logInfo("WalletHandler.GetVP responded successfully with retrieved vp id: %s", id)
}

// GET /wallet/vp/list
func (h *WalletHandler) ListVP(w http.ResponseWriter, r *http.Request) {
	logInfo("WalletHandler.ListVP called")

	// Assuming a models.VPFilter struct
	filter := models.VPFilter{} // later: parse query params like ?verifier=...&domain=...

	// Assuming a WalletvpSvc
	vps, err := h.WalletvpSvc.ListVPs(filter)
	if err != nil {
		logError("Failed to List VPs: %v", err)
		http.Error(w, "failed to list VPs: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(vps)
	logInfo("WalletHandler.ListVP responded successfully with retrieved vps: %v", vps)
}

// POST /wallet/vp/verify
// Note: This implies the wallet is verifying its own stored VP (e.g., a "dry run" before sending)
func (h *WalletHandler) VerifyVP(w http.ResponseWriter, r *http.Request) {
	logInfo("WalletHandler.VerifyVP called")
	var req struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logError("Invalid VP verify request: %v", err)
		http.Error(w, "invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Assuming a WalletvpSvc
	vp, err := h.WalletvpSvc.GetVP(req.ID)
	if err != nil {
		logError("VP ID, %s not found: %v", req.ID, err)
		http.Error(w, "VP not found: "+err.Error(), http.StatusNotFound)
		return
	}

	// This service method would perform a self-check (e.g., check signatures, expiration)
	valid, err := h.WalletvpSvc.VerifyVP(vp)
	if err != nil {
		logError("VP ID, %s verification failed: %v", req.ID, err)
		http.Error(w, "verification error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"id":     req.ID,
		"valid":  valid,
		"status": "verified",
	})
	logInfo("WalletHandler.VerifyVP responded successfully with verified vp: %v", req.ID)
}

// POST /wallet/vp/build
func (h *WalletHandler) BuildVP(w http.ResponseWriter, r *http.Request) {
	logInfo("WalletHandler.BuildVP called")
	var req struct {
		VCIDs        []string            `json:"vc_ids"`
		RevealFields map[string][]string `json:"reveal_fields"`
		Nonce        string              `json:"nonce"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logError("Invalid VP JSON request to build: %v", err)
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	vp, err := h.WalletvcSvc.BuildVP(req.VCIDs, req.RevealFields, req.Nonce)
	if err != nil {
		logError("Build VP failed: %v", err)
		http.Error(w, "failed to build VP: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Use MarshalIndent for a "pretty-printed" log output
	vpJSON, jsonErr := json.MarshalIndent(vp, "", "  ") // prefix, indent
	if jsonErr != nil {
		logError("Failed to marshal VP for logging: %v", jsonErr)
		logInfo("WalletHandler.BuildVP responded successfully")
	} else {
		// Add a newline \n to make the log easy to read
		logInfo("WalletHandler.BuildVP responded successfully with vp:\n%s", string(vpJSON))
	}

	json.NewEncoder(w).Encode(vp)
	//logInfo("WalletHandler.BuildVP responded successfully with vp: %v", vp)
}
