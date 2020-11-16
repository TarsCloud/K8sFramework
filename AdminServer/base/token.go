package base

type RequestAccount struct {
	Role int
	Id   string
	Name string
}

type TokenHandler interface {
	LoadAccountFormToken(token string) (account *RequestAccount, errMessage string, errCode int)
}
