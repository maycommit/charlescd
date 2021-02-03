package git

import (
	"fmt"
	"os"

	"github.com/maycommit/circlerr/internal/k8s/controller/cache/project"
	"github.com/maycommit/circlerr/internal/k8s/controller/env"
	apperror "github.com/maycommit/circlerr/internal/k8s/controller/error"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func CloneAndOpenRepository(repoUrl string) (*git.Repository, apperror.Error) {
	gitDirOut := fmt.Sprintf("%s/%s", env.Get("GIT_DIR"), repoUrl)

	r, err := git.PlainClone(gitDirOut, false, &git.CloneOptions{
		URL:      repoUrl,
		Progress: os.Stdout,
	})
	if err != nil && err != git.ErrRepositoryAlreadyExists {
		return nil, apperror.New("Clone and open repository failed", err.Error()).AddOperation("git.CloneAndOpenRepository.PlainClone")
	}

	r, err = git.PlainOpen(gitDirOut)
	if err != nil {
		return nil, apperror.New("Clone and open repository failed", err.Error()).AddOperation("git.CloneAndOpenRepository.PlainOpen")
	}

	return r, nil
}

func SyncRepository(projectName string, projectCache *project.ProjectCache) (string, apperror.Error) {
	r, cloneAndOpenError := CloneAndOpenRepository(projectCache.RepoURL)
	if cloneAndOpenError != nil {
		return "", cloneAndOpenError
	}

	w, err := r.Worktree()
	if err != nil {
		return "", apperror.New("Sync repository failed", err.Error()).AddOperation("git.SyncRepository.Worktree")
	}

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return "", apperror.New("Sync repository failed", err.Error()).AddOperation("git.SyncRepository.Pull")
	}

	h, err := r.ResolveRevision(plumbing.Revision("HEAD"))
	if err != nil {
		return "", nil
	}

	return h.String(), nil
}
