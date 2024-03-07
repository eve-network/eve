package airdrop

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/go-github/v60/github"
)

type SnapshotOpt struct {
	Since *time.Time `json:"since"`
	Util  *time.Time `json:"util"`
	Repo  string     `json:"repo"`
}

func IsAllTime(opt SnapshotOpt) bool {
	return opt.Since == nil && opt.Util == nil
}

type ContributorSnapshot interface {
	// Run runs snapshot process store all contributors by given SnapshotOpts
	Run(ctx context.Context, opts ...SnapshotOpt)
}

type githubSnapshot struct {
	logger       *slog.Logger
	githubClient *github.Client
	owner        string
	storage      Storage
}

var _ ContributorSnapshot = &githubSnapshot{}

func NewGithubSnapshot(owner string) (ContributorSnapshot, error) {
	return &githubSnapshot{
		logger:       slog.Default(),
		githubClient: github.NewClient(nil),
		owner:        owner,
	}, nil
}

// Run implements ContributorSnapshot.
func (gs *githubSnapshot) Run(ctx context.Context, opts ...SnapshotOpt) {
	for _, opt := range opts {
		commits, err := gs.listCommits(ctx, opt)
		if err != nil {
			gs.logger.Error("Run: error when get list commits", "repo", opt.Repo)
			continue
		}
		contributors := gs.buildContributors(commits)
		if err := gs.storage.Store(opt.Repo, contributors); err != nil {
			gs.logger.Error("Run: error when store contributor", "repo", opt.Repo)
			continue
		}
	}
}

func (gs *githubSnapshot) buildContributors(_ []*github.RepositoryCommit) []Contributor {
	return []Contributor{}
}

func (gs *githubSnapshot) listCommits(ctx context.Context, opt SnapshotOpt) ([]*github.RepositoryCommit, error) {
	// TODO: handle query all pages

	commitsOpt := &github.CommitsListOptions{}
	var results []*github.RepositoryCommit
	result, _, err := gs.githubClient.Repositories.ListCommits(
		ctx,
		gs.owner,
		opt.Repo,
		commitsOpt,
	)
	if err != nil {
		return nil, err
	}
	results = append(results, result...)

	if len(results) == 0 {
		return nil, fmt.Errorf("no commits found")
	}

	return results, nil
}
