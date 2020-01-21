package db_errors

type CodebaseIsUsedByCDPipeline struct {
	Status   StatusReason
	Message  string
	Codebase string
	Pipeline string
}

func (e CodebaseIsUsedByCDPipeline) Error() string {
	return string(e.Status)
}
