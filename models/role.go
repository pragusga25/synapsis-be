package models

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleCustomer Role = "customer"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleCustomer:
		return true
	}
	return false
}

func (r Role) String() string {
	return string(r)
}
