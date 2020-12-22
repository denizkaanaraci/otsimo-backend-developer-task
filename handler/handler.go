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
		errStr := fmt.Sprintf("[ERROR] Error occured at decoding JSON! Details: %s", err.Error())
		log.Println(errStr)
		helper.RespondWithError(w, http.StatusInternalServerError, errStr)
		return
	}
	err = candidate.Validate()

	if err != nil {
		errStr := fmt.Sprintf("[ERROR] JSON Validation Problem! Details: %s", err.Error())
		log.Println(errStr)
		helper.RespondWithError(w, http.StatusInternalServerError, errStr)
		return
	}
	// insert candidate model.
	err = h.cHandler.CreateCandidate(&candidate)

	// create error
	if err != nil {
		errStr := fmt.Sprintf("[ERROR] Error occured at create candidate! Details: %s", err.Error())
		log.Println(errStr)
		helper.RespondWithError(w, http.StatusInternalServerError, errStr)
		return
	}
	log.Println("CreateCandidate is handled successfully!")
	helper.RespondWithJSON(w, http.StatusOK, candidate)
}

func (h *Handler) ReadCandidate(w http.ResponseWriter, r *http.Request) {

	var id, _ = mux.Vars(r)["id"]

	candidate, err := h.cHandler.ReadCandidate(id)

	if err != nil {
		errStr := fmt.Sprintf("[ERROR] Cannot read candidate! Details: %s", err.Error())
		log.Println(errStr)
		helper.RespondWithError(w, http.StatusInternalServerError, errStr)
		return
	}

	log.Println("ReadCandidate is handled successfully!")

	helper.RespondWithJSON(w, http.StatusOK, candidate)
}

func (h *Handler) DeleteCandidate(w http.ResponseWriter, r *http.Request) {
	var id, _ = mux.Vars(r)["id"]

	err := h.cHandler.DeleteCandidate(id)

	if err != nil {
		errStr := fmt.Sprintf("[ERROR] Error occured at delete candidate! Details: %s", err.Error())
		log.Println(errStr)
		helper.RespondWithError(w, http.StatusInternalServerError, errStr)
		return
	}
	log.Println("DeleteCandidate is handled successfully!")
	helper.RespondWithJSON(w, http.StatusOK, "Candidate is deleted.")
}

func (h *Handler) ArrangeMeeting(w http.ResponseWriter, r *http.Request) {

	var payload map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		errStr := fmt.Sprintf("[ERROR] Error occured at decoding JSON! Details: %s", err.Error())
		log.Println(errStr)
		helper.RespondWithError(w, http.StatusInternalServerError, errStr)
		return
	}

	_id := fmt.Sprint(payload["_id"])
	_time, _ := strconv.ParseInt(fmt.Sprint(payload["nextMeetingTime"]), 10, 64)

	nextMeetingTime := time.Unix(_time, 0)
	err = h.cHandler.ArrangeMeeting(_id, &nextMeetingTime)

	if err != nil {
		errStr := fmt.Sprintf("[ERROR] Error occured at arrange meeting! Details: %s", err.Error())
		log.Println(errStr)
		helper.RespondWithError(w, http.StatusInternalServerError, errStr)
		return
	}

	log.Println("ArrangeMeeting is handled successfully!")
	helper.RespondWithJSON(w, http.StatusOK, "Meeting is arranged.")
}

func (h *Handler) CompleteMeeting(w http.ResponseWriter, r *http.Request) {

	var id, _ = mux.Vars(r)["id"]

	err := h.cHandler.CompleteMeeting(id)

	if err != nil {
		errStr := fmt.Sprintf("[ERROR] Error occured at complete meeting! Details: %s", err.Error())
		log.Println(errStr)
		helper.RespondWithError(w, http.StatusInternalServerError, errStr)
		return
	}

	log.Println("CompleteMeeting is handled successfully!")
	helper.RespondWithJSON(w, http.StatusOK, "Meeting is completed.")
}

func (h *Handler) DenyCandidate(w http.ResponseWriter, r *http.Request) {
	var params = mux.Vars(r)

	_id, _ := params["id"]

	err := h.cHandler.DenyCandidate(_id)

	if err != nil {
		errStr := fmt.Sprintf("[ERROR] Error occured at deny candidate! Details: %s", err.Error())
		log.Println(errStr)
		helper.RespondWithError(w, http.StatusInternalServerError, errStr)
		return
	}
	log.Println("DenyCandidate is handled successfully!")
	helper.RespondWithJSON(w, http.StatusOK, "Candidate is denied.")
}

func (h *Handler) AcceptCandidate(w http.ResponseWriter, r *http.Request) {
	var params = mux.Vars(r)

	_id, _ := params["id"]

	err := h.cHandler.AcceptCandidate(_id)

	if err != nil {
		errStr := fmt.Sprintf("[ERROR] Error occured at accept candidate! Details: %s", err.Error())
		log.Println(errStr)
		helper.RespondWithError(w, http.StatusInternalServerError, errStr)
		return
	}
	log.Println("AcceptCandidate is handled successfully!")
	helper.RespondWithJSON(w, http.StatusOK, "Candidate is accepted.")
}

func (h *Handler) FindAssigneeIDByName(w http.ResponseWriter, r *http.Request) {

	var params = mux.Vars(r)
	name, _ := params["name"]

	assignee, err := h.cHandler.FindAssigneeIDByName(name)

	if err != nil {
		errStr := fmt.Sprintf("[ERROR] Error occured at find assignee id! Details: %s", err.Error())
		log.Println(errStr)
		helper.RespondWithError(w, http.StatusInternalServerError, errStr)
		return
	}
	log.Println("FindAssigneeIDByName is handled successfully!")
	helper.RespondWithJSON(w, http.StatusOK, assignee)
}

func (h *Handler) FindAssigneesCandidates(w http.ResponseWriter, r *http.Request) {
	var id, _ = mux.Vars(r)["id"]

	candidates, err := h.cHandler.FindAssigneesCandidates(id)

	if err != nil {
		errStr := fmt.Sprintf("[ERROR] Error occured at find assigneees candidates! Details: %s", err.Error())
		log.Println(errStr)
		helper.RespondWithError(w, http.StatusInternalServerError, errStr)
		return
	}

	log.Println("FindAssigneesCandidates is handled successfully!")
	helper.RespondWithJSON(w, http.StatusOK, candidates)
}
