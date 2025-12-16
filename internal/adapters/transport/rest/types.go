package rest

import "time"

// DTOs kept in a single place to make transport shapes easy to find.
type customerInput struct {
	Name       string            `json:"name"`
	Email      string            `json:"email"`
	TaxID      string            `json:"tax_id"`
	Phone      string            `json:"phone"`
	Country    string            `json:"country"`
	IsActive   bool              `json:"is_active"`
	Attributes map[string]string `json:"attributes"`
}

type customerResponse struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Email      string            `json:"email"`
	TaxID      string            `json:"tax_id"`
	Phone      string            `json:"phone"`
	Country    string            `json:"country"`
	IsActive   bool              `json:"is_active"`
	Attributes map[string]string `json:"attributes"`
	CreatedAt  time.Time         `json:"created_at,omitempty"`
	UpdatedAt  time.Time         `json:"updated_at,omitempty"`
}

type pageResponse struct {
	NextPageToken string `json:"next_page_token"`
	TotalSize     int32  `json:"total_size"`
}

type listResponse struct {
	Customers []customerResponse `json:"customers"`
	Page      pageResponse       `json:"page"`
}
