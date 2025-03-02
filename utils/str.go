package utils

const (
	inboxLen = len("inbox")
)

/* Strips "inbox" prefix from the string */
func WithoutInboxPrefix(path string) string {
	return path[inboxLen:]
}
