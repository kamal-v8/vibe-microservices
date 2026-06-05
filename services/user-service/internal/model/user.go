// Package model defines the domain types for users and the DTOs used in API
// request/response cycles. Validation tags are applied so Gin's ShouldBindJSON
// can enforce constraints before any database call.
package model

import (
	"time"
)

// User is the domain model that maps 1-to-1 with the "users" table in PostgreSQL.
// The `json` tags control API serialization; `db` tags document the column mapping
// even though we use raw database/sql (handy for future migration to sqlx).
type User struct {
	ID        string    `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Bio       string    `json:"bio" db:"bio"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateUserRequest is the payload expected for POST /api/v1/users.
// Username and Email are required; Bio is optional and defaults to empty.
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Bio      string `json:"bio" binding:"max=1000"`
}

// UpdateUserRequest is the payload expected for PUT /api/v1/users/:id.
// All fields are pointers so we can distinguish "field not sent" (nil) from
// "field sent as empty string" — this enables true partial updates.
type UpdateUserRequest struct {
	Username *string `json:"username" binding:"omitempty,min=3,max=50"`
	Email    *string `json:"email" binding:"omitempty,email,max=255"`
	Bio      *string `json:"bio" binding:"omitempty,max=1000"`
}
