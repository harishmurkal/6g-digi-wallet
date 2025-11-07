// internal/api/handlers/did_handlers.go
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harishmurkal/6g-digi-wallet/internal/models"
	"github.com/harishmurkal/6g-digi-wallet/internal/service/crypto6g"
	"github.com/harishmurkal/6g-digi-wallet/internal/service/issuer"
)

type IssuerHandler struct {
	IssuerService issuer.IssuerService
	CryptoService crypto6g.CryptoService
}

func (h *IssuerHandler) GenerateDID(w http.ResponseWriter, r *http.Request) {
	logInfo("IssuerHandler.GenerateDID called")
	var req struct {
		Method  string         `json:"method"`
		Options map[string]any `json:"options"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logError("Decode error: %v", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	didDoc, err := h.IssuerService.GenerateDID(req.Method, req.Options)
	if err != nil {
		logError("GenerateDID failed: %v", err)
		http.Error(w, "error generating DID: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(didDoc)
	logInfo("IssuerHandler.GenerateDID responded successfully")
}

func (h *IssuerHandler) ResolveDID(w http.ResponseWriter, r *http.Request) {
	logInfo("IssuerHandler.ResolveDID called")
	vars := mux.Vars(r)
	id := vars["id"]

	didDoc, err := h.IssuerService.ResolveDID(id)
	if err != nil {
		logError("ResolveDID failed: %v", err)
		http.Error(w, "error resolving DID: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(didDoc)
	logInfo("IssuerHandler.ResolveDID responded successfully")
}

func (h *IssuerHandler) CreateVC(w http.ResponseWriter, r *http.Request) {
	logInfo("IssuerHandler.CreateVC called")

	var req models.VCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logError("invalid VC request: %v", err)
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	logInfo("Request received for issuerDID=%s, subjectDID=%s", req.IssuerDID, req.SubjectDID)

	vc, err := h.IssuerService.CreateVC(&req)
	if err != nil {
		logError("CreateVC failed: %v", err)
		http.Error(w, "error creating VC: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vc)
	logInfo("IssuerHandler.CreateVC responded successfully with VC ID: %s", vc.ID)
}

// similarly for Resolve, List ...
