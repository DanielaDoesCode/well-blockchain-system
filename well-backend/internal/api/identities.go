package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"Thesis/well-blockchain/well-backend/internal/db"
)

type provisionIdentityRequest struct {
	NationalID     string `json:"national_id"`
	DisplayName    string `json:"display_name"`
	Role           string `json:"role"`
	Address        string `json:"address"`
	SummaryAddress string `json:"summary_address"`
	Commitment     string `json:"commitment"`
	InternalID     string `json:"internal_id"`
}

type identityResponse struct {
	NationalID     string `json:"national_id"`
	DisplayName    string `json:"display_name"`
	Role           string `json:"role"`
	Address        string `json:"address"`
	SummaryAddress string `json:"summary_address"`
	Commitment     string `json:"commitment"`
	InternalID     string `json:"internal_id"`
	Status         string `json:"status"`
}

func toIdentityResponse(e *db.Identity) identityResponse {
	return identityResponse{
		NationalID:     e.NationalID,
		DisplayName:    e.DisplayName,
		Role:           e.Role,
		Address:        e.Address,
		SummaryAddress: e.SummaryAddress,
		Commitment:     e.Commitment,
		InternalID:     e.InternalID,
		Status:         e.Status,
	}
}

var identityRoles = map[string]bool{
	"patient":        true,
	"doctor":         true,
	"first_responder": true,
	"hospital":       true,
}

// lookupRoles are the roles allowed to resolve identities by name / national ID.
var lookupRoles = map[uint8]bool{
	RoleDoctor:         true,
	RoleFirstResponder: true,
	RoleHospital:       true,
	RoleAdmin:          true,
}

// handleProvisionIdentity is the admin/issuer enrollment step: it binds a real
// national ID and display name to an already-registered on-chain identity.
func (h *Handler) handleProvisionIdentity(w http.ResponseWriter, r *http.Request) {
	if err := h.requireAdmin(r.Context()); err != nil {
		writeJSON(w, http.StatusForbidden, errorResponse{Error: err.Error()})
		return
	}

	var req provisionIdentityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON"})
		return
	}
	if req.NationalID == "" || req.DisplayName == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "national_id and display_name are required"})
		return
	}
	if !addrRe.MatchString(req.Address) {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid address"})
		return
	}
	if !identityRoles[req.Role] {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "role must be one of patient, doctor, first_responder, hospital"})
		return
	}
	if req.Commitment != "" && !uuidRe.MatchString(req.Commitment) {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid commitment"})
		return
	}
	if req.InternalID != "" && !uuidRe.MatchString(req.InternalID) {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid internal_id"})
		return
	}

	e := &db.Identity{
		NationalID:     req.NationalID,
		DisplayName:    req.DisplayName,
		Role:           req.Role,
		Address:        strings.ToLower(req.Address),
		SummaryAddress: strings.ToLower(req.SummaryAddress),
		Commitment:     req.Commitment,
		InternalID:     req.InternalID,
		Status:         "active",
	}
	if err := h.store.UpsertIdentity(r.Context(), e); err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "db: " + err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"identity": toIdentityResponse(e)})
}

// handleDeactivateIdentity tombstones an enrolled identity (admin only).
func (h *Handler) handleDeactivateIdentity(w http.ResponseWriter, r *http.Request) {
	if err := h.requireAdmin(r.Context()); err != nil {
		writeJSON(w, http.StatusForbidden, errorResponse{Error: err.Error()})
		return
	}

	nationalID := strings.TrimPrefix(r.URL.Path, "/api/v1/identities/")
	if nationalID == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "missing national_id"})
		return
	}
	if err := h.store.DeactivateIdentity(r.Context(), nationalID); err != nil {
		if errors.Is(err, db.ErrIdentityNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "identity not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "db: " + err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "revoked"})
}

// handleLookupIdentity resolves an identity by national ID or display name,
// returning the on-chain address (plus summary/commitment) for the requester
// to use in a subsequent transaction. Restricted to care/trusted roles and
// every lookup is audited.
func (h *Handler) handleLookupIdentity(w http.ResponseWriter, r *http.Request) {
	caller, err := ExtractUser(r.Context())
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: err.Error()})
		return
	}
	role, err := h.roles.RoleOf(r.Context(), caller)
	if err != nil {
		writeJSON(w, http.StatusForbidden, errorResponse{Error: "unregistered user"})
		return
	}
	if !lookupRoles[role] {
		writeJSON(w, http.StatusForbidden, errorResponse{Error: "role not allowed to resolve identities"})
		return
	}

	q := r.URL.Query()
	nationalID := q.Get("national_id")
	name := q.Get("name")
	if (nationalID == "") == (name == "") {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "exactly one of national_id or name is required"})
		return
	}

	switch {
	case nationalID != "":
		e, err := h.store.GetIdentityByNationalID(r.Context(), nationalID)
		if err != nil {
			if errors.Is(err, db.ErrIdentityNotFound) {
				_ = h.store.RecordLookup(r.Context(), &db.IdentityLookup{
					CallerAddress: strings.ToLower(caller), QueryKey: "national_id",
					QueryValue: nationalID, ResultCount: 0,
				})
				writeJSON(w, http.StatusNotFound, errorResponse{Error: "identity not found"})
				return
			}
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "db: " + err.Error()})
			return
		}
		if e.Status != "active" {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "identity not found"})
			return
		}
		_ = h.store.RecordLookup(r.Context(), &db.IdentityLookup{
			CallerAddress: strings.ToLower(caller), QueryKey: "national_id",
			QueryValue: nationalID, ResultCount: 1,
		})
		writeJSON(w, http.StatusOK, map[string]any{"identity": toIdentityResponse(e)})

	case name != "":
		matches, err := h.store.GetIdentitiesByName(r.Context(), name)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "db: " + err.Error()})
			return
		}
		_ = h.store.RecordLookup(r.Context(), &db.IdentityLookup{
			CallerAddress: strings.ToLower(caller), QueryKey: "name",
			QueryValue: name, ResultCount: len(matches),
		})
		if len(matches) == 0 {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "identity not found"})
			return
		}
		if len(matches) > 1 {
			writeJSON(w, http.StatusConflict, errorResponse{Error: "multiple identities match; query by national_id instead"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"identity": toIdentityResponse(matches[0])})
	}
}