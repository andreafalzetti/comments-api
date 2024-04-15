package comments

type Request struct {
	ID       string
	Data     string
	ClientID string
}

type Response struct {
	ID   string
	Data string
}

const (
	ActionSignIn  = "SIGN_IN"
	ActionSignOut = "SIGN_OUT"
	ActionWhoami  = "WHOAMI"
)
