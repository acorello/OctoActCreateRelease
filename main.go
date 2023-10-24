package main

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/google/go-github/v38/github"
	"golang.org/x/oauth2"
)

type RepoInfo struct {
	Token, Name, Owner string
}

type Assets string

func main() {
	var repo RepoInfo
	flag.StringVar(&repo.Token, "auth-token", "", "authentication token")
	flag.StringVar(&repo.Owner, "repo-owner", "", "repository owner name")
	flag.StringVar(&repo.Name, "repo", "", "repository")

	var assetsDir string
	flag.StringVar(&assetsDir, "assets-dir", "", "folder containing files to upload to this release and nothing else")

	var r = newRelease()
	flag.StringVar(r.TagName, "tag-name", "", "")
	flag.StringVar(r.Name, "release-name", "", "name for this release")
	flag.BoolVar(r.Draft, "is-draft", true, "is this a draft-release?")
	flag.BoolVar(r.Prerelease, "is-pre-release", true, "is this a pre-release?")

	flag.Parse()

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: repo.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// the github client sucks
	r, _, err := client.Repositories.CreateRelease(ctx, repo.Owner, repo.Name, r)
	if err != nil {
		fmt.Println("Failed to create the release:", err)
		return
	}

	t.Execute(os.Stdout, r)

	if !strings.HasSuffix(assetsDir, "/") {
		assetsDir = assetsDir + "/"
	}

	// this passing errors around seems verbose!
	err = filepath.WalkDir(assetsDir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error at %q: %v", path, err)
		}
		if info.IsDir() {
			return nil
		}
		uploadError := UploadAsset(ctx, client, *r.ID, assetsDir, path, &repo)
		if uploadError != nil {
			return fmt.Errorf("error uploading %q: %v", path, uploadError)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Failed to upload the assets:", err)
	}
}

var t = Must(template.New("report").Parse(`Release Created
	WEB:	{{ .GetHTMLURL }}
	API:	{{ .GetURL }}
`))

func UploadAsset(ctx context.Context, client *github.Client, releaseId int64, basePath, fullpath string, repo *RepoInfo) error {
	file, err := os.Open(fullpath)
	if err != nil {
		return err
	}
	defer file.Close()
	uploadOpts := github.UploadOptions{
		Name: strings.TrimPrefix(fullpath, basePath),
	}
	_, _, err = client.Repositories.UploadReleaseAsset(ctx, repo.Owner, repo.Name, releaseId, &uploadOpts, file)
	return err
}

func newRelease() *github.RepositoryRelease {
	return &github.RepositoryRelease{
		TagName:                new(string),
		TargetCommitish:        new(string),
		Name:                   new(string),
		Body:                   new(string),
		DiscussionCategoryName: new(string),
		Draft:                  new(bool),
		Prerelease:             new(bool),
	}
}

func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}
