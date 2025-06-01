package models

import "time"

type User struct {
	ID        int       `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	NoHp      string    `json:"no_hp"`
	Alamat    string    `json:"alamat"`
	RoleID    int       `json:"role_id"`
	CreatedBy int       `json:"created_by"`
	UpdatedBy int       `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted string    `json:"isDeleted"`
}

type CreateUserInput struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	NoHp      string `json:"no_hp"`
	Alamat    string `json:"alamat"`
	RoleID    int    `json:"role_id"`
	CreatedBy int    `json:"created_by"`
}
