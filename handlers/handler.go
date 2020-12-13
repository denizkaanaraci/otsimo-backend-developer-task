package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"otsimo-backend-developer-task/functions"
	"otsimo-backend-developer-task/helper"
	"otsimo-backend-developer-task/models"
	"strconv"
	"time"
)

// Handler passes database connection to functions package
type Handler struct {
	cHandler *functions.Handler
}

// NewBaseHandler returns a new BaseHandler
func NewHandler(cHandler *functions.Handler) *Handler {
	return &Handler{
		cHandler: cHandler,
	}
}

// CreateCandidate (candidate Candidate) (Candidate, error)
func (h *Handler) CreateCandidate(w http.ResponseWriter, r *http.Request) {

	var candidate models.Candidate
	err := json.NewDecoder(r.Body).Decode(&candidate)

	// decode error
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	// insert candidate model.
	err = h.cHandler.CreateCandidate(candidate)

	// create error
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, candidate)
}

//ReadCandidate (_id string) (Candidate, error)
func (h *Handler) ReadCandidate(w http.ResponseWriter, r *http.Request) {

	var id, _ = mux.Vars(r)["id"]

	candidate, err := h.cHandler.ReadCandidate(id)

	log.Println(candidate, err)

	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, candidate)
}

//DeleteCandidate (_id string) error
//func DeleteCandidate(_id string) {
func (h *Handler) DeleteCandidate(w http.ResponseWriter, r *http.Request) {
	return
}

//ArrangeMeeting (_id string, nextMeetingTime *time.Time) error
//func ArrangeMeeting(_id string, nextMeetingTime *time.Time) {
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

	err = h.cHandler.ArrangeMeeting(_id, nextMeetingTime)

	if err != nil {
		log.Println(err.Error())
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, "success")
}

//CompleteMeeting (_id string) error
//func CompleteMeeting(_id string) {
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

//DenyCandidate (_id string) error
func (h *Handler) DenyCandidate(w http.ResponseWriter, r *http.Request) {
	var params = mux.Vars(r)

	_id, _ := params["id"]

	err := h.cHandler.DenyCandidate(_id)

	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, "success")
}

//AcceptCandidate(_id string) error
func (h *Handler) AcceptCandidate(w http.ResponseWriter, r *http.Request) {
	var params = mux.Vars(r)

	_id, _ := params["id"]

	err := h.cHandler.AcceptCandidate(_id)

	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, "success")
}

//FindAssigneeIDByName (name string) string
func (h *Handler) FindAssigneeIDByName(w http.ResponseWriter, r *http.Request) {

	var params = mux.Vars(r)
	name, _ := params["name"]

	assignee, err := h.cHandler.FindAssigneeIDByName(name)

	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, assignee)
}

//bonus
//FindAssigneesCandidates (id string) ([]Candidate, error)
func FindAssigneesCandidates(id string) {
	return
}
