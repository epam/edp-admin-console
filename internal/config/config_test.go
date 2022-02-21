package config

import (
	"context"
	"testing"

	"github.com/astaxie/beego"
	"github.com/stretchr/testify/assert"
)

func TestSetupConfig(t *testing.T) {
	mustSetAppConfig(t, AuthEnable, "false")

	ctx := context.Background()
	emptyStr := ""
	config, err := SetupConfig(ctx, emptyStr)
	assert.NoError(t, err)
	assert.False(t, config.AuthEnable)

}

func mustSetAppConfig(t *testing.T, k, v string) {
	t.Helper()
	err := beego.AppConfig.Set(k, v)
	if err != nil {
		t.Fatal(err)
	}
}
