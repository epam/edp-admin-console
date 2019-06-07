package query

type Status string

const (
	Active   Status = "active"
	Inactive Status = "inactive"
	Failed   Status = "failed"
)
