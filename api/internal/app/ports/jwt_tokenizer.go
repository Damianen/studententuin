package ports

type JwtTokenizer interface {
	CreateToken(userID string) (string, error)
	VerifyToken(tokenStr string) (string, error)
}

