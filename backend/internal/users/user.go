package users

import "regexp"

// ---------- Enums ----------

type AccountType string
type UserRole string
type PermissionLevel string

const (
	AccountFree AccountType = "free"
	AccountPaid AccountType = "paid"
)

const (
	RoleStudent UserRole = "student"
	RoleTeacher UserRole = "teacher"
)

const (
	PermUser  PermissionLevel = "user"
	PermAdmin PermissionLevel = "admin"
)

// ---------- User model ----------

type User struct {
	ID           uint64          `json:"id"`
	PhoneNumber  string          `json:"phone_number"` // +976XXXXXXXX
	PasswordHash string          `json:"-"`

	AccountType AccountType     `json:"account_type"`
	Role        UserRole        `json:"role"`
	Permission  PermissionLevel `json:"permission"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

// ---------- Validation ----------

var mongoliaPhoneRegex = regexp.MustCompile(`^\+976\d{8}$`)

func IsValidMongoliaPhone(phone string) bool {
	return mongoliaPhoneRegex.MatchString(phone)
}
