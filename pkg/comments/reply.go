package comments

import (
	"fmt"
	"strings"
)

type Reply struct {
	author    string
	text      string
	createdAt int64
}

func (r *Reply) String() string {
	replyText := r.text
	if strings.Contains(replyText, "\"") || strings.Contains(replyText, ",") {
		replyText = strings.ReplaceAll(replyText, "\"", "\"\"")
		replyText = "\"" + replyText + "\""
	}
	return fmt.Sprintf("%s|%s", r.author, replyText)
}
