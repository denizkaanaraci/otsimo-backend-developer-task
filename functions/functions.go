package functions

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

// CreateCandidate (candidate Candidate) (Candidate, error)
func (h *Handler) CreateCandidate(candidate models.Candidate) error {
	//TODO: """The email format for the candidate should be example@email.xyz.
	//	Otherwise, the candidate should not be inserted to DB because the only
	//	way to communicate with the candidate is through email."""

	assigneeID, err := h.GetAssigneeIdByDepartment(candidate.Department)

	if err != nil {
		return err
	}

	candidate.Assignee = assigneeID
	candidate.ID = h.CreateUniqueId()
	candidate.Status = "Pending"
	candidate.MeetingCount = 0
	now := time.Now()
	candidate.ApplicationDate = primitive.DateTime(now.Unix())
	//candidate.NextMeeting = primitive.DateTime((now.AddDate(0, 0, 1)).Unix()) // default 1 day
	_, err = h.db.Collection("Candidates").InsertOne(context.TODO(), candidate)

	return err
}

//ReadCandidate (_id string) (Candidate, error)
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

//ReadCandidate (_id string) (Candidate, error)
func (h *Handler) ReadCandidate(_id string) (*models.Candidate, error) {

	var candidate models.Candidate

	filter := bson.M{"_id": _id}
	err := h.db.Collection("Candidates").FindOne(context.Background(), filter).Decode(&candidate)
	return &candidate, err
}

//GetAssigneeByDepartment (_id string) (Candidate, error)
func (h *Handler) GetAssigneeIdByDepartment(department string) (string, error) {

	var assignee models.Assignee

	filter := bson.M{"department": department}
	err := h.db.Collection("Assignees").FindOne(context.Background(), filter).Decode(&assignee)

	return assignee.ID, err
}

//DeleteCandidate (_id string) error
func (h *Handler) DeleteCandidate(_id string) error {

	filter := bson.M{"_id": _id}
	_, err := h.db.Collection("Candidates").DeleteOne(context.Background(), filter)
	return err
}

//ArrangeMeeting (_id string, nextMeetingTime *time.Time) error
func (h *Handler) ArrangeMeeting(_id string, nextMeetingTime time.Time) error {

	//assigneeID, err := h.GetAssigneeIdByDepartment(candidate.Department)

	filter := bson.M{"_id": _id}
	opt := bson.D{
		{"$set", bson.D{{"next_meeting", primitive.DateTime(nextMeetingTime.Unix())}}},
	}
	_, err := h.db.Collection("Candidates").UpdateOne(context.Background(), filter, opt)
	//TODO: if id is not in db
	return err
}

//CompleteMeeting (_id string) error
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

//DenyCandidate (_id string) error
func (h *Handler) DenyCandidate(_id string) error {

	filter := bson.M{"_id": _id}
	opt := bson.D{
		{"$set", bson.D{{"status", "Denied"}, {}}},
	}
	_, err := h.db.Collection("Candidates").UpdateOne(context.Background(), filter, opt)

	return err
}

//AcceptCandidate(_id string) error
func (h *Handler) AcceptCandidate(_id string) error {

	var candidate models.Candidate

	filter := bson.M{"_id": _id}
	err := h.db.Collection("Candidates").FindOne(context.Background(), filter).Decode(&candidate)

	// if object is empty
	if err != nil {
		return err
	}

	if candidate.MeetingCount == 4 {
		opt := bson.D{
			{"$set", bson.D{{"status", "Accepted"}, {}}},
		}
		_, err = h.db.Collection("Candidates").UpdateOne(context.Background(), filter, opt)

	}

	return err
}

//FindAssigneeIDByName (name string) string
func (h *Handler) FindAssigneeIDByName(name string) (*models.Assignee, error) {

	var assignee models.Assignee

	filter := bson.M{"name": name}
	err := h.db.Collection("Assignees").FindOne(context.Background(), filter).Decode(&assignee)
	return &assignee, err
}

//bonus
//FindAssigneesCandidates (id string) ([]Candidate, error)
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
