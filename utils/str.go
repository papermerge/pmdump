package utils

const (
	inboxLen = len("inbox")
	homeLen  = len("home")
)

/* Strips "inbox" prefix from the string */
func WithoutInboxPrefix(path string) string {
	return path[inboxLen:]
}

/* Strips "home" prefix from the string */
func WithoutHomePrefix(path string) string {
	return path[homeLen:]
}
