package comments

import (
	"github.com/google/uuid"
	"strings"
	"time"
)

type Request struct {
	ID        string
	Data      string
	ClientID  string
	Reference string
	Comment   string
}

type Response struct {
	ID   string
	Data string
}

type Reply struct {
	author    string
	text      string
	createdAt int64
}

type Discussion struct {
	id        string
	ref       string
	replies   []Reply
	createdAt int64
}

const (
	ActionSignIn           = "SIGN_IN"
	ActionSignOut          = "SIGN_OUT"
	ActionWhoami           = "WHOAMI"
	ActionCreateDiscussion = "CREATE_DISCUSSION"
	ActionCreateReply      = "CREATE_REPLY"
	ActionGetDiscussion    = "GET_DISCUSSION"
	ActionListDiscussions  = "LIST_DISCUSSIONS"
)

// Unmarshall converts a raw string into a Request object
func Unmarshall(raw string) *Request {
	// the message is divided into parts separated by pipe |
	// first 7 chars are the request id
	// the following segment is the data
	// the last is the client id (optional)
	r := &Request{}
	raw = strings.TrimSuffix(raw, "\n")
	parts := strings.Split(raw, "|")

	r.ID = parts[0]
	if len(parts) == 4 {
		r.Data = parts[1]
		r.Reference = parts[2]
		r.Comment = parts[3]
	} else if len(parts) == 3 {
		r.Data = parts[1]
		r.ClientID = parts[2]
	} else if len(parts) == 2 {
		r.Data = parts[1]
	}
	return r
}

// NewDiscussion generates a unique id, and stores the timestamp of creation
func NewDiscussion(ref, author, comment string) *Discussion {
	now := time.Now().Unix()
	return &Discussion{
		id:        uuid.New().String(),
		ref:       ref,
		replies:   []Reply{{author: author, text: comment, createdAt: now}},
		createdAt: now,
	}
}

func (d *Discussion) GetId() string {
	return d.id
}

func (d *Discussion) AddReply(author, reply string) {
	now := time.Now().Unix()
	d.replies = append(d.replies, Reply{author: author, text: reply, createdAt: now})
}

// GetReplies returns a string wrapped by parenthesis, each item separated by comma, each item containing author and text separated by pipe. Example: (janedoe|"Hey, folks. What do you think of my video? Does it have enough ""polish""?",johndoe|I think it's great!).
// Additionally, each text is wrapped by double quotes and any double quote within the text is escaped with another double quote.
func (d *Discussion) GetReplies() string {
	var replies []string
	for _, r := range d.replies {
		// Escape double quotes in the text
		escapedText := strings.ReplaceAll(r.text, "\"", "\"\"")
		// Wrap the text with double quotes and combine with author using pipe
		formattedReply := r.author + "|\"" + escapedText + "\""
		replies = append(replies, formattedReply)
	}
	// Join all formatted replies with commas and wrap the result with parentheses
	return "(" + strings.Join(replies, ",") + ")"
}
