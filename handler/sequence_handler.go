package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/miledxz/salesforge/model"
	"github.com/miledxz/salesforge/service"
)

func CreateSequence(w http.ResponseWriter, r *http.Request) {
	var seq model.Sequence
	if err := json.NewDecoder(r.Body).Decode(&seq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	createdSeq, err := service.CreateSequence(&seq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdSeq)
}

func UpdateStep(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var step model.SequenceStep
	if err := json.NewDecoder(r.Body).Decode(&step); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = service.UpdateStep(id, &step)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Step updated"})
}

func DeleteStep(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	err = service.DeleteStep(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func UpdateTracking(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var input struct {
		OpenTrackingEnabled  *bool `json:"openTrackingEnabled"`
		ClickTrackingEnabled *bool `json:"clickTrackingEnabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = service.UpdateTracking(id, input.OpenTrackingEnabled, input.ClickTrackingEnabled)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Tracking updated"})
}
