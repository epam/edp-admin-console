package query

type CodebaseDockerStream struct {
	Id                int             `json:"id" orm:"column(id)"`
	OcImageStreamName string          `json:"ocImageStreamName" orm:"column(oc_image_stream_name)"`
	CodebaseBranch    *CodebaseBranch `json:"-" orm:"rel(fk)"`
}

func (c *CodebaseDockerStream) TableName() string {
	return "codebase_docker_stream"
}
