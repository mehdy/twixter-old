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

func (t *Twitter) Profile(username string) (*entities.TwitterProfile, error) {
	user, resp, err := t.api.Users.Show(&twitter.UserShowParams{ScreenName: username})
	if err != nil {
		t.logger.As("E").
			WithError(err).
			WithField("username", username).
			Logf("Failed to fetch profile from twitter API")

		return nil, newError(err, "failed to fetch profile from twitter API")
	}
	defer resp.Body.Close()

	return t.asTwitterProfile(*user), nil
}

func (t *Twitter) asTwitterProfile(user twitter.User) *entities.TwitterProfile {
	joinedAt, err := time.Parse(time.RubyDate, user.CreatedAt)
	if err != nil {
		t.logger.As("W").WithError(err).Logf("Failed to parse joinedAt time")
	}

	jsonEntites, err := json.Marshal(user.Entities)
	if err != nil {
		t.logger.As("W").WithError(err).Logf("Failed to convert user.Entites to JSON")
	}

	var ent map[string]interface{}

	err = json.Unmarshal(jsonEntites, &ent)
	if err != nil {
		t.logger.As("W").WithError(err).Logf("Failed to convert user.Entites to JSON")
	}

	return &entities.TwitterProfile{
		TwitterID:           user.IDStr,
		Name:                user.Name,
		Username:            user.ScreenName,
		Location:            user.Location,
		Bio:                 user.Description,
		URL:                 user.URL,
		Email:               user.Email,
		ProfileBannerURL:    user.ProfileBannerURL,
		ProfileImageURL:     user.ProfileImageURLHttps,
		Verified:            user.Verified,
		Protected:           user.Protected,
		DefaultProfile:      user.DefaultProfile,
		DefaultProfileImage: user.DefaultProfileImage,
		FollowersCount:      user.FollowersCount,
		FollowingsCount:     user.FriendsCount,
		FavouritesCount:     user.FavouritesCount,
		ListedCount:         user.ListedCount,
		TweetsCount:         user.StatusesCount,
		Entities:            ent,
		JoinedAt:            joinedAt,
	}
}
