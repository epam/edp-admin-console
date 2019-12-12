package db_errors

type StatusReason string

const (
	StatusReasonNotFound StatusReason = "NotFound"
	StatusReasonUnknown  StatusReason = "Unknown"
)

func IsNotFound(err error) bool {
	return reasonForError(err) == StatusReasonNotFound
}

func reasonForError(err error) StatusReason {
	switch t := err.(type) {
	case StatusError:
		return t.Status
	}
	return StatusReasonUnknown
}
