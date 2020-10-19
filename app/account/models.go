package account

type Account struct {
	username string `json:"username"`
}

type validateAccount struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
