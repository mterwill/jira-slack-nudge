package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/mterwill/jira-slack-nudge"
)

func fromEnv(k string) string {
	v, ok := os.LookupEnv(k)

	if !ok {
		log.Fatalf("Environment variable %s must be set", k)
	}

	return v
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: ./jira-slack-nudge <path-to-config.yaml>")
	}

	filename := os.Args[1]

	rawCfg, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	cfg := nudge.Config{}
	err = yaml.Unmarshal(rawCfg, &cfg)

	p := nudge.New(
		fromEnv("JIRA_SERVER"),
		fromEnv("JIRA_USERNAME"),
		fromEnv("JIRA_PASSWORD"),
		"", // use default config URL in messages
	)

	ctx := context.Background()

	err = p.Run(ctx, &cfg)

	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
