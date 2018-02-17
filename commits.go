package main

import (
	"context"
	"net/http"
	"os"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func changelogHandler(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: os.Getenv("STEAM_GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	options := github.CommitsListOptions{
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 20,
		},
	}

	commits, _, err := client.Repositories.ListCommits(ctx, "steam-authority", "steam-authority", &options)
	if err != nil {
		logger.Error(err)
		returnErrorTemplate(w, 500, err.Error())
		return
	}

	template := changelogTemplate{}
	template.Commits = commits

	returnTemplate(w, "commits", template)
}

type changelogTemplate struct {
	GlobalTemplate
	Commits []*github.RepositoryCommit
}
