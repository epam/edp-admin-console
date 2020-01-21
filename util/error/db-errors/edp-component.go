package db_errors

type StatusReason string

const (
	StatusReasonNotFound                   StatusReason = "NotFound"
	StatusReasonUnknown                    StatusReason = "Unknown"
	StatusReasonCodebaseIsUsedByCDPipeline StatusReason = "Used"
)

func IsNotFound(err error) bool {
	return reasonForError(err) == StatusReasonNotFound
}

func CodebaseIsUsed(err error) bool {
	return reasonForError(err) == StatusReasonCodebaseIsUsedByCDPipeline
}

func reasonForError(err error) StatusReason {
	switch t := err.(type) {
	case StatusError:
		return t.Status
	case CodebaseIsUsedByCDPipeline:
		return t.Status
	}
	return StatusReasonUnknown
}
