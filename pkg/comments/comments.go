package comments

import "strings"

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

func Unmarshall(raw string) *Request {
	// the message is divided into parts separated by pipe |
	// first 7 chars are the request id
	// the following segment is the data
	// the last is the client id (optional)
	r := &Request{}
	raw = strings.TrimSuffix(raw, "\n")
	parts := strings.Split(raw, "|")

	r.ID = parts[0]
	if len(parts) == 3 {
		r.Data = parts[1]
		r.ClientID = parts[2]
	} else if len(parts) == 2 {
		r.Data = parts[1]
	}
	return r
}
