package types

import (
	"github.com/google/uuid"
)

// Standard Multi-Tenant Identifiers
type CompanyID = uuid.UUID
type BranchID = uuid.UUID
type UserID = uuid.UUID

// EmptyID returns an empty UUID
func EmptyID() uuid.UUID {
	return uuid.Nil
}

// NewID generates a new random UUIDv4
func NewID() uuid.UUID {
	return uuid.New()
}

// ParseID parses string into UUID
func ParseID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
