package jwt

type JWTUser struct {
	Username string
	UserId   int
}

type JWTPayload struct {
	UserId   int
	UserName string
	Exp      int64
}
