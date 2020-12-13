package models

import (
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
)

//Create Struct

type Assignee struct {
	ID         string `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string `json:"name" bson:"name"`
	Department string `json:"department" bson:"department"`
}

type Candidate struct {
	ID           string             `json:"_id" bson:"_id" validate:"required"`
	FirstName    string             `json:"first_name" bson:"first_name" validate:"required"`
	LastName     string             `json:"last_name" bson:"last_name" validate:"required"`
	Email        string             `json:"email" bson:"email" validate:"required,email"`
	Department   string             `json:"department" bson:"department" validate:"required"`
	University   string             `json:"university" bson:"university" validate:"required"`
	Experience   bool               `json:"experience" bson:"experience" validate:"required"`
	Status       string             `json:"status" bson:"status"`
	MeetingCount int32              `json:"meeting_count" bson:"meeting_count"`
	NextMeeting  primitive.DateTime `json:"next_meeting" bson:"next_meeting"`
	Assignee     string             `json:"assignee" bson:"assignee"`
	ApplicationDate primitive.DateTime `json:"application_date" bson:"application_date"`
}

func (c *Candidate) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("email", validateEmail)

	return validate.Struct(c)
}

func validateEmail(fl validator.FieldLevel) bool {
	var emailRegex = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	re := regexp.MustCompile(emailRegex)

	if len(fl.Field().String()) < 3 && len(fl.Field().String()) > 254 {
		return false
	}
	return re.MatchString(fl.Field().String())
}
