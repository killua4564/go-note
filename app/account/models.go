package account

type Account struct {
	Username string `json:"username"`
}

type validateAccount struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
