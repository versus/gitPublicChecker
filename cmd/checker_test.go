package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChecker(t *testing.T) {

	assert.True(t, true, "True is true!")

}

func TestShowRepo(t *testing.T) {
	config := new(Config)
	config.Repos = []string{"aaa", "bbb"}
	ShowRepos(*config)
}

func TestCheckRepo(t *testing.T) {
	url := "https://github.com/stretchr/testify"
	CheckRepo(url)
}
