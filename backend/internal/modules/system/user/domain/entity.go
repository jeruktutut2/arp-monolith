package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents an identity within the multi-tenant ERP system
type User struct {
	ID                  uuid.UUID  `json:"id"`
	CompanyID           uuid.UUID  `json:"company_id"`
	BranchID            *uuid.UUID `json:"branch_id,omitempty"`
	Email               string     `json:"email"`
	PasswordHash        string     `json:"-"`
	Name                string     `json:"name"`
	Role                string     `json:"role"`
	IsActive            bool       `json:"is_active"`
	FailedLoginAttempts int        `json:"-"`
	LockedUntil         *time.Time `json:"-"`
	CreatedAt           time.Time  `json:"created_at"`
	CreatedBy           *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt           time.Time  `json:"updated_at"`
	UpdatedBy           *uuid.UUID `json:"updated_by,omitempty"`
}

// SetPassword hashes and sets the user's password
func (u *User) SetPassword(password string) error {
	trimmed := strings.TrimSpace(password)
	if trimmed == "" {
		return errors.New("password tidak boleh kosong")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(trimmed), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

// ValidatePassword checks if raw password matches the stored bcrypt hash
func (u *User) ValidatePassword(password string) bool {
	if u.PasswordHash == "" {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// IsLocked checks whether the account is currently locked out
func (u *User) IsLocked(now time.Time) bool {
	if u.LockedUntil == nil {
		return false
	}
	return now.Before(*u.LockedUntil)
}

// RecordFailedLogin increments failed attempts and locks the account if reaching maxAttempts
func (u *User) RecordFailedLogin(now time.Time, maxAttempts int, lockoutDuration time.Duration) {
	u.FailedLoginAttempts++
	if u.FailedLoginAttempts >= maxAttempts {
		lockedUntil := now.Add(lockoutDuration)
		u.LockedUntil = &lockedUntil
	}
}

// ResetLoginAttempts clears failed login attempts and lockout timestamp
func (u *User) ResetLoginAttempts() {
	u.FailedLoginAttempts = 0
	u.LockedUntil = nil
}
