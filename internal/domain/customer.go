package domain

import "time"

type Customer struct {
    ID         string
    Name       string
    Email      string
    TaxID      string
    Phone      string
    Country    string
    IsActive   bool
    Attributes map[string]string
    CreatedAt  time.Time
    UpdatedAt  time.Time
}
