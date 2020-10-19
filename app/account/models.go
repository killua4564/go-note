package account

type Account struct {
	username string `json:"username"`
}

type validateAccount struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
