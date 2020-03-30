package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/drone/drone-template-lib/template"
	"github.com/google/go-github/github"
	"github.com/wraithgar/slack"
	"golang.org/x/oauth2"
)

type (
	Repo struct {
		Owner string
		Name  string
	}

	Build struct {
		Tag        string
		Event      string
		Number     int
		Commit     string
		CommitLink string
		Ref        string
		Branch     string
		Author     string
		Pull       string
		Message    string
		DeployTo   string
		DeployID   int64
		Status     string
		PrevStatus string
		Link       string
		Started    int64
		Created    int64
	}

	Config struct {
		Webhook       string
		Channel       string
		Recipient     string
		Username      string
		Template      string
		ImageURL      string
		IconURL       string
		IconEmoji     string
		LinkNames     bool
		GhToken       string
		GhToSlackJSON string
	}

	Job struct {
		Started int64
	}

	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
		Job    Job
	}

	//Custom payload object from andbang deploy event
	DeploymentPayload struct {
		ResponseURL string `json:"response_url,omitempty"`
		Person      string `json:"person,omitempty"`
	}

	GhToSlack struct {
		Github string `json:"github"`
		Slack  string `json:"slack"`
	}
)

func (p Plugin) Exec() error {
	fmt.Printf("talky custom drone-slack version %s build %s\n", version, build)
	attachment := slack.Attachment{
		Text:       message(p.Repo, p.Build),
		Fallback:   fallback(p.Repo, p.Build),
		Color:      color(p.Build),
		MarkdownIn: []string{"text", "fallback"},
		ImageURL:   p.Config.ImageURL,
	}

	if p.Config.Template != "" {
		txt, err := template.RenderTrim(p.Config.Template, p)

		if err != nil {
			return err
		}

		attachment.Text = txt
	}

	fmt.Printf("Build event: %s\n", p.Build.Event)
	if p.Build.Event == "promote" {
		//Send a delayed response if we have a response url
		deployPayload, err := getDeploymentPayload(p.Config, p.Repo, p.Build)
		if err != nil {
			return err
		}
		response := slack.ResponsePayload{}

		response.Attachments = []*slack.Attachment{&attachment}
		if p.Config.LinkNames == true {
			response.LinkNames = "1"
		}
		fmt.Printf("%+v\n", response.Attachments)
		// return err
		delayedClient := slack.NewDelayedResponse(deployPayload.ResponseURL)
		return delayedClient.PostDelayedResponse(&response)
	}

	//Send a webhook message
	payload := slack.WebHookPostPayload{}
	payload.Username = p.Config.Username
	payload.Attachments = []*slack.Attachment{&attachment}
	payload.IconUrl = p.Config.IconURL
	payload.IconEmoji = p.Config.IconEmoji
	payload.Attachments = []*slack.Attachment{&attachment}

	//These are the only cases we care about
	if p.Build.Event == "tag" || p.Build.Status == "failure" || p.Build.PrevStatus == "failure" {
		if p.Config.GhToSlackJSON != "" {
			var ghtoslacks []GhToSlack
			jsonerr := json.Unmarshal([]byte(p.Config.GhToSlackJSON), &ghtoslacks)
			if jsonerr != nil {
				fmt.Println(jsonerr)
			} else {
				for _, v := range ghtoslacks {
					if v.Github == p.Build.Author {
						payload.Channel = prepend("@", v.Slack)
						fmt.Printf("%s", payload.Channel)
					}
				}
			}
		} else {
			fmt.Println("GhToSlackJSON not found")
		}

		if payload.Channel == "" {
			if p.Config.Recipient != "" {
				fmt.Printf("%s", p.Config.Recipient)
				payload.Channel = prepend("@", p.Config.Recipient)
			} else if p.Config.Channel != "" {
				payload.Channel = prepend("#", p.Config.Channel)
			}
		}

		if p.Config.LinkNames == true {
			payload.LinkNames = "1"
		}

		fmt.Printf("sending webhook message to %s\n", payload.Channel)
		client := slack.NewWebHook(p.Config.Webhook)
		return client.PostMessage(&payload)
	}

	fmt.Println("Not sending this one to slack")
	fmt.Println(attachment.Text)
	return nil
}

func message(repo Repo, build Build) string {
	return fmt.Sprintf("*%s* <%s|%s/%s#%s> (%s) by %s",
		build.Status,
		build.Link,
		repo.Owner,
		repo.Name,
		build.Commit[:8],
		build.Branch,
		build.Author,
	)
}

func fallback(repo Repo, build Build) string {
	return fmt.Sprintf("%s %s/%s#%s (%s) by %s",
		build.Status,
		repo.Owner,
		repo.Name,
		build.Commit[:8],
		build.Branch,
		build.Author,
	)
}

func color(build Build) string {
	switch build.Status {
	case "success":
		return "good"
	case "failure", "error", "killed":
		return "danger"
	default:
		return "warning"
	}
}

func prepend(prefix, s string) string {
	if !strings.HasPrefix(s, prefix) {
		return prefix + s
	}

	return s
}

func getDeploymentPayload(config Config, repo Repo, build Build) (DeploymentPayload, error) {
	var payload DeploymentPayload
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GhToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	parts := strings.Split(build.Link, "/")
	owner := parts[3]
	name := parts[4]
	id := build.DeployID

	deployment, _, gherr := client.Repositories.GetDeployment(ctx, owner, name, id)
	if gherr != nil {
		fmt.Println(gherr)
		return payload, gherr
	}
	jsonerr := json.Unmarshal(deployment.Payload, &payload)
	if jsonerr != nil {
		fmt.Println(jsonerr)
		return payload, jsonerr
	}
	return payload, nil
}
