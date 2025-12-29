package entity

// Role represents user roles in the system (RBAC)
type Role int

const (
	RoleUser Role = iota + 1
	RoleModerator
	RoleAdmin
)

// String returns string representation of role
func (r Role) String() string {
	switch r {
	case RoleUser:
		return "user"
	case RoleModerator:
		return "moderator"
	case RoleAdmin:
		return "admin"
	default:
		return "unknown"
	}
}

// ParseRole parses string to Role
func ParseRole(s string) Role {
	switch s {
	case "user":
		return RoleUser
	case "moderator":
		return RoleModerator
	case "admin":
		return RoleAdmin
	default:
		return RoleUser
	}
}
