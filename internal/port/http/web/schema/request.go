package webschema

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ChangeUserRequest struct {
	Email    string `json:"email"`
	State    string `json:"state"`
	Status   string `json:"status"`
	Password string `json:"password"`
}

type ConfirmEmailRequest struct {
	Email string `json:"email"`
}

type ConfirmNewEmailRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ConfirmResetPasswordRequest struct {
	Email string `json:"email"`
}

type JWTRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RegisterUserRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}
