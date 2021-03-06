package config_test

import (
	"testing"

	"github.com/git-town/git-town/test"
	"github.com/stretchr/testify/assert"
)

func TestRunner_SetOffline(t *testing.T) {
	repo := test.CreateTestGitTownRepo(t)
	err := repo.Config.SetOffline(true)
	assert.NoError(t, err)
	offline := repo.Config.IsOffline()
	assert.True(t, offline)
	err = repo.Config.SetOffline(false)
	assert.NoError(t, err)
	offline = repo.Config.IsOffline()
	assert.False(t, offline)
}
