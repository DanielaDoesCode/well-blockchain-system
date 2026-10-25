package aggregator

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"

	"Thesis/well-blockchain/well-backend/internal/db"

	"github.com/ethereum/go-ethereum/common"
)

type Handler struct {
	agg *Aggregator
	db  *db.DB
}

type submitRequest struct {
	PA0            string `json:"pA_0"`
	PA1            string `json:"pA_1"`
	PB00           string `json:"pB_0_0"`
	PB01           string `json:"pB_0_1"`
	PB10           string `json:"pB_1_0"`
	PB11           string `json:"pB_1_1"`
	PC0            string `json:"pC_0"`
	PC1            string `json:"pC_1"`
	Commitment     string `json:"commitment"`
	Nullifier      string `json:"nullifier"`
	MerkleRoot     string `json:"merkle_root"`
	ConditionGroup string `json:"condition_group"`
}

type submitResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type countResult struct {
	ConditionGroup string `json:"condition_group"`
	TotalCount     int64  `json:"total_count"`
	PublishedCount int64  `json:"published_count"`
}

type countersResponse struct {
	Counters []countResult `json:"counters"`
}

func NewHandler(agg *Aggregator, database *db.DB) *Handler {
	return &Handler{agg: agg, db: database}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/api/v1/contribute":
		h.handleContribute(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/api/v1/counters":
		h.handleCounters(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) handleContribute(w http.ResponseWriter, r *http.Request) {
	var req submitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, submitResponse{Status: "error", Message: "invalid JSON"})
		return
	}

	if err := validateSubmit(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, submitResponse{Status: "error", Message: err.Error()})
		return
	}

	commitment := common.HexToHash(req.Commitment)
	nullifier  := common.HexToHash(req.Nullifier)

	registered, err := h.agg.IsRegistered(r.Context(), commitment)
	if err != nil {
		log.Printf("aggregator: isRegistered: %v", err)
		writeJSON(w, http.StatusInternalServerError, submitResponse{Status: "error", Message: "verification failed"})
		return
	}
	if !registered {
		writeJSON(w, http.StatusBadRequest, submitResponse{Status: "error", Message: "commitment not registered"})
		return
	}

	pA, pB, pC, err := parseProof(&req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, submitResponse{Status: "error", Message: err.Error()})
		return
	}

	merkleRoot := new(big.Int)
	if _, ok := merkleRoot.SetString(req.MerkleRoot, 10); !ok {
		// also accept hex with 0x prefix
		if _, ok := merkleRoot.SetString(strings.TrimPrefix(req.MerkleRoot, "0x"), 16); !ok {
			writeJSON(w, http.StatusBadRequest, submitResponse{Status: "error", Message: "invalid merkle_root"})
			return
		}
	}

	valid, err := h.agg.VerifyProof(r.Context(), pA, pB, pC, commitment, nullifier, merkleRoot)
	if err != nil {
		log.Printf("aggregator: verifyProof: %v", err)
		writeJSON(w, http.StatusInternalServerError, submitResponse{Status: "error", Message: "proof verification failed"})
		return
	}
	if !valid {
		writeJSON(w, http.StatusBadRequest, submitResponse{Status: "error", Message: "invalid proof"})
		return
	}

	contrib := &db.Contribution{
		Commitment:     req.Commitment,
		Nullifier:      req.Nullifier,
		MerkleRoot:     merkleRoot.String(), // stored as decimal string
		PA0:            req.PA0,
		PA1:            req.PA1,
		PB00:           req.PB00,
		PB01:           req.PB01,
		PB10:           req.PB10,
		PB11:           req.PB11,
		PC0:            req.PC0,
		PC1:            req.PC1,
		ConditionGroup: req.ConditionGroup,
	}

	if err := h.db.InsertPendingContribution(r.Context(), contrib); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			writeJSON(w, http.StatusConflict, submitResponse{Status: "error", Message: "duplicate contribution"})
			return
		}
		log.Printf("aggregator: insert: %v", err)
		writeJSON(w, http.StatusInternalServerError, submitResponse{Status: "error", Message: "storage error"})
		return
	}

	if err := h.db.IncrementCounter(r.Context(), req.ConditionGroup); err != nil {
		log.Printf("aggregator: increment counter: %v", err)
	}

	writeJSON(w, http.StatusOK, submitResponse{Status: "accepted"})
}

func (h *Handler) handleCounters(w http.ResponseWriter, r *http.Request) {
	counters, err := h.db.GetCounters(r.Context())
	if err != nil {
		log.Printf("aggregator: get counters: %v", err)
		writeJSON(w, http.StatusInternalServerError, countersResponse{Counters: []countResult{}})
		return
	}

	resp := countersResponse{}
	for _, c := range counters {
		resp.Counters = append(resp.Counters, countResult{
			ConditionGroup: c.ConditionGroup,
			TotalCount:     c.TotalCount,
			PublishedCount: c.PublishedCount,
		})
	}
	if resp.Counters == nil {
		resp.Counters = []countResult{}
	}

	writeJSON(w, http.StatusOK, resp)
}

func validateSubmit(req *submitRequest) error {
	if req.Commitment == "" {
		return fmt.Errorf("commitment is required")
	}
	if req.Nullifier == "" {
		return fmt.Errorf("nullifier is required")
	}
	if req.MerkleRoot == "" {
		return fmt.Errorf("merkle_root is required")
	}
	if req.ConditionGroup == "" {
		return fmt.Errorf("condition_group is required")
	}
	if !common.IsHexAddress(req.Commitment) && !strings.HasPrefix(req.Commitment, "0x") {
		return fmt.Errorf("commitment must be a hex string")
	}
	return nil
}

func parseProof(req *submitRequest) ([2]*big.Int, [2][2]*big.Int, [2]*big.Int, error) {
	var pA [2]*big.Int
	var pB [2][2]*big.Int
	var pC [2]*big.Int

	for i, s := range []struct {
		val *string
		dst **big.Int
	}{
		{&req.PA0, &pA[0]}, {&req.PA1, &pA[1]},
		{&req.PB00, &pB[0][0]}, {&req.PB01, &pB[0][1]},
		{&req.PB10, &pB[1][0]}, {&req.PB11, &pB[1][1]},
		{&req.PC0, &pC[0]}, {&req.PC1, &pC[1]},
	} {
		if *s.val == "" {
			return pA, pB, pC, fmt.Errorf("missing proof field at index %d", i)
		}
		*s.dst = new(big.Int)
		if _, ok := (*s.dst).SetString(*s.val, 10); !ok {
			return pA, pB, pC, fmt.Errorf("invalid proof field at index %d", i)
		}
	}

	return pA, pB, pC, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

var _ = context.Background
