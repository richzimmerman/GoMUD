package interfaces

type AccountInterface interface {
	AccountName() string
	LastIPLogged() string
	EmailAddress() string
	GetCharacters() []string
	GetPassword() string
	ChangePassword(s string) error
	LoggedInStatus() bool
	SetLoggedInStatus(b bool)
}
