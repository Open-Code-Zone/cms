package utils

import (
	"context"
	"fmt"

	"github.com/google/go-github/v53/github"
	"github.com/Open-Code-Zone/cms/config"
	"golang.org/x/oauth2"
)

const (
	Owner = "waishnav"
	Repo  = "danto-hugo"
)

type GitHubClient struct {
	client *github.Client
	ctx    context.Context
}

func NewGitHubClient() (*GitHubClient, error) {
	token := config.Envs.GitHubToken

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &GitHubClient{
		client: client,
		ctx:    ctx,
	}, nil
}

func (c *GitHubClient) CreateFile(path, content, message string) error {
	fileOpts := &github.RepositoryContentFileOptions{
		Message: &message,
		Content: []byte(content),
	}

	_, _, err := c.client.Repositories.CreateFile(c.ctx, Owner, Repo, path, fileOpts)
	return err
}

func (c *GitHubClient) UpdateFile(path, content, message string) error {
	fileContent, _, _, err := c.client.Repositories.GetContents(c.ctx, Owner, Repo, path, nil)
	if err != nil {
		return fmt.Errorf("failed to get file content: %w", err)
	}

	sha := fileContent.GetSHA()

	fileOpts := &github.RepositoryContentFileOptions{
		Message: &message,
		Content: []byte(content),
		SHA:     &sha, // Required for updating the file
	}

	_, _, err = c.client.Repositories.UpdateFile(c.ctx, Owner, Repo, path, fileOpts)
	if err != nil {
		return fmt.Errorf("failed to update file: %w", err)
	}

	return nil
}

func (c *GitHubClient) DeleteFile(path, message string) error {
	fileContent, _, _, err := c.client.Repositories.GetContents(c.ctx, Owner, Repo, path, nil)
	if err != nil {
		return fmt.Errorf("failed to get file content: %w", err)
	}

	sha := fileContent.GetSHA()
	fileOpts := &github.RepositoryContentFileOptions{
		SHA:     &sha,
		Message: &message,
	}

	_, _, err = c.client.Repositories.DeleteFile(c.ctx, Owner, Repo, path, fileOpts)
	return err
}

func (c *GitHubClient) GetContents(path string) ([]*github.RepositoryContent, error) {
	_, contents, _, err := c.client.Repositories.GetContents(c.ctx, Owner, Repo, path, nil)
	return contents, err
}

func (c *GitHubClient) GetFileContent(path string) (string, error) {
	fileContent, _, _, err := c.client.Repositories.GetContents(c.ctx, Owner, Repo, path, nil)
	if err != nil {
		return "", err
	}
	return fileContent.GetContent()
}

func (c *GitHubClient) ListCommits(path string) ([]*github.RepositoryCommit, error) {
	commits, _, err := c.client.Repositories.ListCommits(c.ctx, Owner, Repo, &github.CommitsListOptions{
		Path: path,
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 1,
		},
	})
	return commits, err
}
