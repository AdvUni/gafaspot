package constants

const (
	// TimeLayout is the layout which must follow all time strings stored to and retrieved from database.
	// It can be used with time.ParseInLocation() and time.Format(). Such time strings are perfectly comparable.
	// As they don't contain any time zone information, it should be always parsed with time.Local.
	// So, each time value inside gafaspot are interpreted in the local time zone of the running server.
	TimeLayout = "2006-01-02 15:04"
)
