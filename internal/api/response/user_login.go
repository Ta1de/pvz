package response

type LoginPostRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
