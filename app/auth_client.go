package app

type AuthClient interface {
	SubjectFrom(string) (string, error)
}
