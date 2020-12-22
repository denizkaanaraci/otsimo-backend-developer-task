package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"otsimo-backend-developer-task/helper"
	"otsimo-backend-developer-task/models"
	"otsimo-backend-developer-task/storage"
	"strconv"
	"time"
)

// Handler passes database connection to storage package
type Handler struct {
	cHandler *storage.Handler
}

// NewBaseHandler returns a new BaseHandler
func NewHandler(cHandler *storage.Handler) *Handler {
	return &Handler{
		cHandler: cHandler,
	}
}

func (h *Handler) CreateCandidate(w http.ResponseWriter, r *http.Request) {

	var candidate models.Candidate
	err := json.NewDecoder(r.Body).Decode(&candidate)

	// decode error
	if err != nil {
		log.Println(err.Error())
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = candidate.Validate()

	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		log.Println("[ERROR] validating email", err)
		return
	}
	// insert candidate model.
	err = h.cHandler.CreateCandidate(&candidate)

	// create error
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		log.Println("[ERROR] error occured at create candidate", err)
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, candidate)
}

func (h *Handler) ReadCandidate(w http.ResponseWriter, r *http.Request) {

	var id, _ = mux.Vars(r)["id"]

	candidate, err := h.cHandler.ReadCandidate(id)

	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		log.Println(err.Error())
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, candidate)
}

func (h *Handler) DeleteCandidate(w http.ResponseWriter, r *http.Request) {
	var id, _ = mux.Vars(r)["id"]

	err := h.cHandler.DeleteCandidate(id)

	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		log.Println(err.Error())
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, "success")
}

func (h *Handler) ArrangeMeeting(w http.ResponseWriter, r *http.Request) {

	var payload map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		log.Println(err.Error())
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	_id := fmt.Sprint(payload["_id"])
	_time, _ := strconv.ParseInt(fmt.Sprint(payload["nextMeetingTime"]), 10, 64)

	nextMeetingTime := time.Unix(_time, 0)
	err = h.cHandler.ArrangeMeeting(_id, &nextMeetingTime)

	if err != nil {
		log.Println(err.Error())
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, "success")
}

func (h *Handler) CompleteMeeting(w http.ResponseWriter, r *http.Request) {

	var id, _ = mux.Vars(r)["id"]

	err := h.cHandler.CompleteMeeting(id)

	if err != nil {
		log.Println(err.Error())
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, "success")
}

func (h *Handler) DenyCandidate(w http.ResponseWriter, r *http.Request) {
	var params = mux.Vars(r)

	_id, _ := params["id"]

	err := h.cHandler.DenyCandidate(_id)

	if err != nil {
		log.Println(err.Error())
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, "success")
}

func (h *Handler) AcceptCandidate(w http.ResponseWriter, r *http.Request) {
	var params = mux.Vars(r)

	_id, _ := params["id"]

	err := h.cHandler.AcceptCandidate(_id)

	if err != nil {
		log.Println(err.Error())
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, "success")
}

func (h *Handler) FindAssigneeIDByName(w http.ResponseWriter, r *http.Request) {

	var params = mux.Vars(r)
	name, _ := params["name"]

	assignee, err := h.cHandler.FindAssigneeIDByName(name)

	if err != nil {
		log.Println(err.Error())
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, assignee)
}

func (h *Handler) FindAssigneesCandidates(w http.ResponseWriter, r *http.Request) {
	var id, _ = mux.Vars(r)["id"]

	candidates, err := h.cHandler.FindAssigneesCandidates(id)

	if err != nil {
		log.Println(err.Error())
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Println(candidates, err)

	helper.RespondWithJSON(w, http.StatusOK, candidates)
}
