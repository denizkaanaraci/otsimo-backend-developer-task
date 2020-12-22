package storage

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"otsimo-backend-developer-task/models"
	"time"
)

// BaseHandler will hold everything that controller needs
type Handler struct {
	db *mongo.Database
}

// NewBaseHandler returns a new BaseHandler
func NewHandler(db *mongo.Database) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) CreateCandidate(candidate *models.Candidate) error {
	assigneeID, err := h.GetAssigneeIdByDepartment(candidate.Department)

	// if there is no assignee with this given department
	if err != nil {
		return err
	}

	candidate.Assignee = assigneeID
	candidate.ID = h.CreateUniqueId()
	candidate.Status = "Pending"
	candidate.MeetingCount = 0
	now := time.Now()
	candidate.ApplicationDate = &now
	_, err = h.db.Collection("Candidates").InsertOne(context.Background(), candidate)

	return err
}

func (h *Handler) CreateUniqueId() string {

	var _id string
	var candidate models.Candidate
	for true {
		_id = primitive.NewObjectID().Hex()

		filter := bson.M{"_id": _id}
		err := h.db.Collection("Candidates").FindOne(context.Background(), filter).Decode(&candidate)
		if err != nil {
			break
		}
	}
	return _id
}

func (h *Handler) ReadCandidate(_id string) (*models.Candidate, error) {

	var candidate models.Candidate

	filter := bson.M{"_id": _id}
	err := h.db.Collection("Candidates").FindOne(context.Background(), filter).Decode(&candidate)
	return &candidate, err
}

func (h *Handler) GetAssigneeIdByDepartment(department string) (string, error) {

	var assignee models.Assignee

	filter := bson.M{"department": department}
	err := h.db.Collection("Assignees").FindOne(context.Background(), filter).Decode(&assignee)

	return assignee.ID, err
}

func (h *Handler) DeleteCandidate(_id string) error {

	filter := bson.M{"_id": _id}
	_, err := h.db.Collection("Candidates").DeleteOne(context.Background(), filter)
	return err
}

func (h *Handler) ArrangeMeeting(_id string, nextMeetingTime *time.Time) error {

	_, err := h.ReadCandidate(_id)

	if err != nil {
		log.Println("ID not found in DB")
		return err
	}
	filter := bson.M{"_id": _id}
	opt := bson.D{
		{"$set", bson.D{{"next_meeting", primitive.DateTime(nextMeetingTime.Unix())}}},
	}
	_, err = h.db.Collection("Candidates").UpdateOne(context.Background(), filter, opt)
	//TODO: if id is not in db
	return err
}

func (h *Handler) CompleteMeeting(_id string) error {

	var candidate models.Candidate

	filter := bson.M{"_id": _id}
	err := h.db.Collection("Candidates").FindOne(context.Background(), filter).Decode(&candidate)

	// if object is empty
	if err != nil {
		return err
	}

	if candidate.MeetingCount == 0 {
		candidate.Status = "Pending"
		candidate.Status = "In Progress"
	}

	candidate.MeetingCount += 1

	if candidate.MeetingCount == 3 {
		assigneeId, err := h.GetAssigneeIdByDepartment("CEO")
		// if CEO obj is not in db
		if err != nil {
			return err
		}
		candidate.Assignee = assigneeId
	}

	_, err = h.db.Collection("Candidates").ReplaceOne(context.Background(), filter, candidate)

	return err
}

func (h *Handler) DenyCandidate(_id string) error {

	var candidate models.Candidate

	filter := bson.M{"_id": _id}

	err := h.db.Collection("Candidates").FindOne(context.Background(), filter).Decode(&candidate)

	// if object is empty
	if err != nil {
		return errors.New("ID not found in DB, candidate object is empty!")
	}

	filter = bson.M{"_id": _id}
	opt := bson.D{
		{"$set", bson.D{{"status", "Denied"}}},
	}
	_, err = h.db.Collection("Candidates").UpdateOne(context.Background(), filter, opt)

	return err
}

func (h *Handler) AcceptCandidate(_id string) error {

	var candidate models.Candidate

	filter := bson.M{"_id": _id}
	err := h.db.Collection("Candidates").FindOne(context.Background(), filter).Decode(&candidate)

	// if object is empty
	if err != nil {
		return errors.New("ID not found in DB, candidate object is empty!")
	}

	if candidate.MeetingCount == 4 {
		opt := bson.D{
			{"$set", bson.D{{"status", "Accepted"}}},
		}
		_, err = h.db.Collection("Candidates").UpdateOne(context.Background(), filter, opt)
	} else {
		return errors.New("Cannot accept candidate! Meeting count is less than 4")
	}

	return err
}

func (h *Handler) FindAssigneeIDByName(name string) (*models.Assignee, error) {

	var assignee models.Assignee

	filter := bson.M{"name": name}
	err := h.db.Collection("Assignees").FindOne(context.Background(), filter).Decode(&assignee)
	return &assignee, err
}

//bonus
func (h *Handler) FindAssigneesCandidates(_id string) ([]*models.Candidate, error) {

	filter := bson.M{"assignee": _id}
	cursor, err := h.db.Collection("Candidates").Find(context.Background(), filter)

	var candidates []*models.Candidate

	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.Background(), &candidates); err != nil {
		return nil, err
	}
	return candidates, err
}
