package comments

import (
	"strings"
)

type Request struct {
	ID           string
	Data         string
	ClientID     string
	Reference    string
	Comment      string
	DiscussionId string
}

type Response struct {
	ID   string
	Data string
}

// Unmarshall converts a raw string into a Request object
// the message is divided into parts separated by pipe |
// first 7 chars are the request id
// the following segment is the data (aka the request type)
// following it can either be the client_id, or different info depending on the request type
func Unmarshall(raw string) *Request {
	r := &Request{}
	raw = strings.TrimSuffix(raw, "\n")
	parts := strings.Split(raw, "|")

	r.ID = parts[0]

	if len(parts) > 1 {
		r.Data = parts[1]
	}

	switch r.Data {
	case ActionSignIn:
		// <request_id>|SIGN_IN|<client_id>
		r.ClientID = parts[2]

	case ActionWhoami:
		// <request_id>|WHOAMI

	case ActionSignOut:
		// <request_id>|SIGN_OUT

	case ActionCreateDiscussion:
		// <request_id>|CREATE_DISCUSSION|<reference>|<comment>
		r.Reference = parts[2]
		r.Comment = parts[3]

	case ActionCreateReply:
		// <request_id>|CREATE_REPLY|<discussion_id>|<comment>
		r.DiscussionId = parts[2]
		r.Comment = parts[3]

	case ActionGetDiscussion:
		// <request_id>|GET_DISCUSSION|<discussion_id>
		r.DiscussionId = parts[2]

	case ActionListDiscussions:
		// <request_id>|LIST_DISCUSSIONS|<reference_prefix>
		r.Reference = parts[2]
	default:
		// no-op
	}

	return r
}
