package entity

type Admin struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	EmailIsVerified bool      `json:"emailIsVerified"`
	Password        any       `json:"password"`
	Role            string    `json:"role"`
	CreatedBy       CreatedBy `json:"createdBy"`
	CreatedAt       string    `json:"createdAt"`
	UpdatedAt       string    `json:"updatedAt,omitempty"`
}

type CreatedBy struct {
	UserID string `json:"userId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}
