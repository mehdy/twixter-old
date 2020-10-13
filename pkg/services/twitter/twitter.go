package twitter

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/mehdy/twixter/pkg/entities"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	// nolint: gosec // no credentials
	twitterTokenURL = "https://api.twitter.com/oauth2/token"
)

type Twitter struct {
	api    *twitter.Client
	logger entities.Logger
}

func NewTwitter(config entities.ConfigGetter, logger entities.Logger) *Twitter {
	creds := &clientcredentials.Config{
		ClientID:     config.GetString("TWITTER_CONSUMER_KEY"),
		ClientSecret: config.GetString("TWITTER_CONSUMER_SECRET"),
		TokenURL:     twitterTokenURL,
	}

	return &Twitter{
		api:    twitter.NewClient(creds.Client(context.TODO())),
		logger: logger,
	}
}
