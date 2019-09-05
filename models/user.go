package models

type User struct {
	ID int64 `json:"id"`
	Email string `json:"email"`
	Fullname string `json:"fullname"`
	Password string `json:"password"`
	UpdatedAt int64 `json:"updatedAt"`
	CreatedAt int64 `json:"createdAt"`
}