package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/inconshreveable/go-update"
)

var (
	githubUrl = "https://api.github.com/repos/XescuGC/scrap/releases/latest"
)

type GithubRepository struct {
	URL             string    `json:"url"`
	HTMLURL         string    `json:"html_url"`
	AssetsURL       string    `json:"assets_url"`
	UploadURL       string    `json:"upload_url"`
	TarballURL      string    `json:"tarball_url"`
	ZipballURL      string    `json:"zipball_url"`
	ID              int       `json:"id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Body            string    `json:"body"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Author          struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	Assets []struct {
		URL                string    `json:"url"`
		BrowserDownloadURL string    `json:"browser_download_url"`
		ID                 int       `json:"id"`
		Name               string    `json:"name"`
		Label              string    `json:"label"`
		State              string    `json:"state"`
		ContentType        string    `json:"content_type"`
		Size               int       `json:"size"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		Uploader           struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"uploader"`
	} `json:"assets"`
}

func doUpdate() error {
	githubRepository, err := getGithubRepository()

	if githubRepository.TagName == version {
		fmt.Printf("Nothing to update it's the same version: %v :)\n", version)
		return nil
	} else {
		fmt.Prinf("Updating from version %v to %v\n", version, githubRepository.TagName)
	}

	if err != nil {
		return err
	}

	resp, err := http.Get("https://github.com/XescuGC/scrap/releases/download/" + githubRepository.TagName + "/scrap")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = update.Apply(resp.Body, update.Options{})

	if err != nil {
		return err
	}

	fmt.Printf("Successfully updated to version: %v\n", githubRepository.TagName)
	return nil
}

func getGithubRepository() (*GithubRepository, error) {
	resp, err := httpClient.Get(githubUrl)

	checkError("", err)

	defer resp.Body.Close()

	githubRepository := new(GithubRepository)

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(githubRepository)

	if err != nil {
		return nil, err
	}

	return githubRepository, nil
}
