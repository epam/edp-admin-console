package util

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"log"
)

func IsGitRepoAvailable(repo string, user string, pass string) bool {
	r, _ := git.Init(memory.NewStorage(), nil)
	remote, _ := r.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{repo},
	})
	rfs, err := remote.List(&git.ListOptions{
		Auth: &http.BasicAuth{
			Username: user,
			Password: pass,
		}})
	if err != nil {
		log.Println(err)
		return false
	}
	return len(rfs) != 0
}
