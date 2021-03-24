package util

import (
	"github.com/astaxie/beego"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateNativeProjectLinkMethod_ShouldBeExecutedSuccessfully(t *testing.T) {
	beego.AppConfig.Set("projectMaskUrl", "/console/project/{namespace}/overview")
	link := CreateNativeProjectLink("https://stub-domain", "stub-project")
	assert.Equal(t, "https://stub-domain/console/project/stub-project/overview", link)
}

func TestCreateDockerStreamLinkMethod_ShouldBeExecutedSuccessfully(t *testing.T) {
	beego.AppConfig.Set("imageStreamMaskUrl", "/console/project/{namespace}/browse/images/{stream}")
	link := CreateDockerStreamLink("https://stub-domain", "stub-project", "stub-stream")
	assert.Equal(t, "https://stub-domain/console/project/stub-project/browse/images/stub-stream", link)
}

func TestCreateGitlabCILinkMethod_ShouldBeExecutedSuccessfully(t *testing.T) {
	link := CreateGitlabCILink("stub-domain", "stub-relative-path")
	assert.Equal(t, "https://stub-domainstub-relative-path/pipelines?scope=branches&page=1", link)
}
