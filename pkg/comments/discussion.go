package comments

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

type Discussion struct {
	id        string
	ref       string
	replies   []Reply
	createdAt int64
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

func (d *Discussion) String() string {
	return fmt.Sprintf("%s|%s|%s", d.id, d.ref, d.GetReplies())
}

func (d *Discussion) GetId() string {
	return d.id
}

func (d *Discussion) GetReference() string {
	return d.ref
}

func (d *Discussion) AddReply(author, reply string) {
	now := time.Now().Unix()
	d.replies = append(d.replies, Reply{author: author, text: reply, createdAt: now})
}

// GetReplies returns the replies formatted as follows: wrapped by parenthesis, each item separated by comma, each item containing author and text separated by pipe.
// Additionally, each text is wrapped by double quotes and any double quote within the text is escaped with another double quote.
// Example: (janedoe|"Hey, folks. What do you think of my video? Does it have enough ""polish""?",johndoe|I think it's great!).
func (d *Discussion) GetReplies() string {
	var replies []string
	for _, r := range d.replies {
		replies = append(replies, r.String())
	}
	return "(" + strings.Join(replies, ",") + ")"
}

// GetParticipants returns the list of users who wrote comments in the discussion.
func (d *Discussion) GetParticipants() []string {
	var participants []string
	for _, r := range d.replies {
		participants = append(participants, r.author)
	}
	return participants
}
