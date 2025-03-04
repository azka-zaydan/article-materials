package dto

type CreateUserReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
