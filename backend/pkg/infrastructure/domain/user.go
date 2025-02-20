package domain

type User struct {
	ID       string `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	LastName string `json:"lastName" db:"last_name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}
