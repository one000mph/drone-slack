package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	version = "2.1.0"
	build   = "5"
)

func main() {

	app := cli.NewApp()
	app.Name = "slack plugin"
	app.Usage = "slack plugin"
	app.Action = run
	app.Version = fmt.Sprintf("%s+%s", version, build)
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "webhook",
			Usage:   "slack webhook url",
			EnvVars: []string{"SLACK_WEBHOOK", "PLUGIN_SLACK_WEBHOOK,"},
		},
		&cli.StringFlag{
			Name:    "ghtoken",
			Usage:   "github access token",
			EnvVars: []string{"GITHUB_ACCESS_TOKEN", "PLUGIN_GITHUB_ACCESS_TOKEN"},
		},
		&cli.StringFlag{
			Name:    "ghtoslacks",
			Usage:   "github_slack_lookup",
			EnvVars: []string{"GITHUB_SLACK_LOOKUP", "PLUGIN_GITHUB_SLACK_LOOKUP"},
		},
		&cli.StringFlag{
			Name:    "channel",
			Usage:   "slack channel",
			EnvVars: []string{"PLUGIN_CHANNEL"},
		},
		&cli.StringFlag{
			Name:    "recipient",
			Usage:   "slack recipient",
			EnvVars: []string{"PLUGIN_RECIPIENT"},
		},
		&cli.StringFlag{
			Name:    "username",
			Usage:   "slack username",
			EnvVars: []string{"PLUGIN_USERNAME"},
		},
		&cli.StringFlag{
			Name:    "template",
			Usage:   "slack template",
			EnvVars: []string{"PLUGIN_TEMPLATE"},
		},
		&cli.BoolFlag{
			Name:    "link-names",
			Usage:   "slack link names",
			EnvVars: []string{"PLUGIN_LINK_NAMES"},
		},
		&cli.StringFlag{
			Name:    "image",
			Usage:   "slack image url",
			EnvVars: []string{"PLUGIN_IMAGE_URL"},
		},
		&cli.StringFlag{
			Name:    "icon.url",
			Usage:   "slack icon url",
			EnvVars: []string{"PLUGIN_ICON_URL"},
		},
		&cli.StringFlag{
			Name:    "icon.emoji",
			Usage:   "slack emoji url",
			EnvVars: []string{"PLUGIN_ICON_EMOJI"},
		},
		&cli.StringFlag{
			Name:    "repo.owner",
			Usage:   "repository owner",
			EnvVars: []string{"DRONE_REPO_OWNER"},
		},
		&cli.StringFlag{
			Name:    "repo.name",
			Usage:   "repository name",
			EnvVars: []string{"DRONE_REPO_NAME"},
		},
		&cli.StringFlag{
			Name:    "commit.sha",
			Usage:   "git commit sha",
			EnvVars: []string{"DRONE_COMMIT_SHA"},
			Value:   "00000000",
		},
		&cli.StringFlag{
			Name:    "commit.link",
			Usage:   "git commit link",
			EnvVars: []string{"DRONE_COMMIT_LINK"},
		},
		&cli.StringFlag{
			Name:    "commit.ref",
			Value:   "refs/heads/master",
			Usage:   "git commit ref",
			EnvVars: []string{"DRONE_COMMIT_REF"},
		},
		&cli.StringFlag{
			Name:    "commit.branch",
			Value:   "master",
			Usage:   "git commit branch",
			EnvVars: []string{"DRONE_COMMIT_BRANCH"},
		},
		&cli.StringFlag{
			Name:    "commit.author",
			Usage:   "git author name",
			EnvVars: []string{"DRONE_COMMIT_AUTHOR"},
		},
		&cli.StringFlag{
			Name:    "commit.pull",
			Usage:   "git pull request",
			EnvVars: []string{"DRONE_PULL_REQUEST"},
		},
		&cli.StringFlag{
			Name:    "commit.message",
			Usage:   "commit message",
			EnvVars: []string{"DRONE_COMMIT_MESSAGE"},
		},
		&cli.StringFlag{
			Name:    "build.event",
			Value:   "push",
			Usage:   "build event",
			EnvVars: []string{"DRONE_BUILD_EVENT"},
		},
		&cli.IntFlag{
			Name:    "build.number",
			Usage:   "build number",
			EnvVars: []string{"DRONE_BUILD_NUMBER"},
		},
		&cli.StringFlag{
			Name:    "build.status",
			Usage:   "build status",
			Value:   "success",
			EnvVars: []string{"DRONE_BUILD_STATUS"},
		},
		&cli.StringFlag{
			Name:    "build.prevstatus",
			Usage:   "build previous status",
			Value:   "success",
			EnvVars: []string{"DRONE_PREV_BUILD_STATUS"},
		},
		&cli.StringFlag{
			Name:    "build.link",
			Usage:   "build link",
			EnvVars: []string{"DRONE_BUILD_LINK"},
		},
		&cli.Int64Flag{
			Name:    "build.started",
			Usage:   "build started",
			EnvVars: []string{"DRONE_BUILD_STARTED"},
		},
		&cli.Int64Flag{
			Name:    "build.created",
			Usage:   "build created",
			EnvVars: []string{"DRONE_BUILD_CREATED"},
		},
		&cli.StringFlag{
			Name:    "build.tag",
			Usage:   "build tag",
			EnvVars: []string{"DRONE_TAG"},
		},
		&cli.StringFlag{
			Name:    "build.deployTo",
			Usage:   "environment deployed to",
			EnvVars: []string{"DRONE_DEPLOY_TO"},
		},
		&cli.StringFlag{
			Name:    "build.deployID",
			Usage:   "github deployment id",
			EnvVars: []string{"DRONE_DEPLOY_ID"},
		},
		&cli.Int64Flag{
			Name:    "job.started",
			Usage:   "job started",
			EnvVars: []string{"DRONE_JOB_STARTED"},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Repo: Repo{
			Owner: c.String("repo.owner"),
			Name:  c.String("repo.name"),
		},
		Build: Build{
			Tag:        c.String("build.tag"),
			Number:     c.Int("build.number"),
			Event:      c.String("build.event"),
			Status:     c.String("build.status"),
			PrevStatus: c.String("build.prevstatus"),
			Commit:     c.String("commit.sha"),
			CommitLink: c.String("commit.link"),
			Ref:        c.String("commit.ref"),
			Branch:     c.String("commit.branch"),
			Author:     c.String("commit.author"),
			Pull:       c.String("commit.pull"),
			Message:    c.String("commit.message"),
			DeployTo:   c.String("build.deployTo"),
			DeployID:   c.Int64("build.deployID"),
			Link:       c.String("build.link"),
			Started:    c.Int64("build.started"),
			Created:    c.Int64("build.created"),
		},
		Job: Job{
			Started: c.Int64("job.started"),
		},
		Config: Config{
			GhToken:       c.String("ghtoken"),
			GhToSlackJSON: c.String("ghtoslacks"),
			Webhook:       c.String("webhook"),
			Channel:       c.String("channel"),
			Recipient:     c.String("recipient"),
			Username:      c.String("username"),
			Template:      c.String("template"),
			ImageURL:      c.String("image"),
			IconURL:       c.String("icon.url"),
			IconEmoji:     c.String("icon.emoji"),
			LinkNames:     c.Bool("link_names"),
		},
	}

	return plugin.Exec()
}
