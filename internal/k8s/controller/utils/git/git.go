package git

import (
	"fmt"
	"strings"

	"github.com/maycommit/circlerr/internal/k8s/controller/env"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func GetOutDir(url string) string {
	return fmt.Sprintf("%s/%s", env.Get("GIT_DIR"), strings.Replace(url, "/", "_", -1))
}

func GetRepository(gitOptions git.CloneOptions) (*git.Repository, error) {
	gitDirOut := GetOutDir(gitOptions.URL)

	r, err := git.PlainOpen(gitDirOut)
	if err != nil && err == git.ErrRepositoryNotExists {
		r, err = git.PlainClone(gitDirOut, false, &gitOptions)
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return r, err
}

func GetRevision(r *git.Repository) (string, error) {
	plugingHash, err := r.ResolveRevision(plumbing.Revision("HEAD"))
	return plugingHash.String(), err
}

func Pull(r *git.Repository) error {
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	return nil
}

func SyncRepository(gitOptions git.CloneOptions, revision string) (string, error) {
	r, err := GetRepository(gitOptions)
	if err != nil {
		return "", err
	}

	remoteRevision, err := GetRevision(r)
	if err != nil {
		return "", err
	}

	if revision == remoteRevision {
		return remoteRevision, nil
	}

	return remoteRevision, Pull(r)
}
