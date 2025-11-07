package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/harishmurkal/6g-digi-wallet/internal/models"
	"github.com/harishmurkal/6g-digi-wallet/internal/service/verifier"
)

type VerifierHandler struct {
	VerifierService verifier.VerifierService
}

// POST /verifier/verify
func (h *VerifierHandler) Verify(w http.ResponseWriter, r *http.Request) {
	logInfo("VerifierHandler.Verify called")
	var verifier models.VerifiablePresentation
	if err := json.NewDecoder(r.Body).Decode(&verifier); err != nil {
		logError("Invalid VP: %v", err)
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	ok, err := h.VerifierService.VerifyVP(&verifier)
	if err != nil {
		logError("Verification of VP Failed: %v", err)
		http.Error(w, "verification failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"verified": ok,
	})

	logInfo("VerifierHandler.Verify responded successfully")
}
