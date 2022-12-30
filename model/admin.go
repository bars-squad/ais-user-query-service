package model

import "github.com/bars-squad/ais-user-query-service/entity"

type AdminRegistration struct {
	ID              string           `validate:"required" json:"_id"`
	Name            string           `validate:"required" json:"name"`
	Email           string           `validate:"required,email" json:"email"`
	EmailIsVerified bool             `json:"emailIsVerified"`
	Password        any              `validate:"required" json:"password"`
	Role            string           `validate:"required" json:"role"`
	CreatedBy       entity.CreatedBy `validate:"required" json:"createdBy"`
	CreatedAt       string           `validate:"required" json:"createdAt"`
	UpdatedAt       string           `json:"updatedAt,omitempty"`
}
