package git

import (
	"charlescd/internal/env"
	"charlescd/internal/manager/circle"
	"fmt"
	"github.com/go-git/go-git/v5"
	"os"
)

func CloneAndOpenRepository(project circle.Project) (*git.Repository, error) {
	gitDirOut := fmt.Sprintf("%s/%s", env.Get("GIT_DIR"), project.Name)

	r, err := git.PlainClone(gitDirOut, false, &git.CloneOptions{
		URL:      project.Repository,
		Progress: os.Stdout,
	})
	if err != nil && err != git.ErrRepositoryAlreadyExists {
		return nil, err
	}

	r, err = git.PlainOpen(gitDirOut)
	if err != nil {
		return nil, err
	}

	return r, nil
}
