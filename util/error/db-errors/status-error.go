package db_errors

type StatusError struct {
	Status StatusReason
}

func (e StatusError) Error() string {
	return string(e.Status)
}
