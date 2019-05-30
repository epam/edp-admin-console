package models

import (
	"errors"
)

var (
	ErrCDPipelineIsExists    = errors.New("cd pipeline is already exists")
	ErrNonValidRelatedBranch = errors.New("application has non valid related branch")
)
