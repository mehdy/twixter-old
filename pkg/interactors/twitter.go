package interactors

import (
	"github.com/mehdy/twixter/pkg/entities"
)

// Ensure Twitter implements the entities.TwitterInteractor interface.
var _ entities.TwitterInteractor = &Twitter{}

type Twitter struct {
	logger  entities.Logger
	twitter TwitterAPI
	store   Store
}

func NewTwitter(logger entities.Logger, twitter TwitterAPI, store Store) *Twitter {
	return &Twitter{logger: logger, twitter: twitter, store: store}
}

func (t *Twitter) UpdateFollowings(username string) error {
	panic("not implemented") // TODO: Implement
}

func (t *Twitter) UpdateFollowers(username string) error {
	panic("not implemented") // TODO: Implement
}

func (t *Twitter) UpdateProfile(username string) error {
	profile, err := t.twitter.Profile(username)
	if err != nil {
		t.logger.As("E").WithError(err).Logf("Failed to fetch profile from the twitter API")

		return newError(err, "failed to update profile")
	}

	t.logger.As("D").WithField("username", username).Logf("Fetched the profile from the twitter API successfully")

	err = t.store.SaveProfile(profile)
	if err != nil {
		t.logger.As("E").WithError(err).Logf("Failed to store profile")

		return newError(err, "failed to update profile")
	}

	t.logger.As("D").WithField("username", username).Logf("Stored the profile successfully")

	return nil
}

func (t *Twitter) UpdateNetwork(username string, followings bool, followers bool, depth int) error {
	panic("not implemented") // TODO: Implement
}

func (t *Twitter) GetTopFollowingsByFollowers(username string, limit int) ([]*entities.TwitterProfile, error) {
	panic("not implemented") // TODO: Implement
}

func (t *Twitter) GetTopFollowersByFollowers(username string, limit int) ([]*entities.TwitterProfile, error) {
	panic("not implemented") // TODO: Implement
}

func (t *Twitter) GetTopFollowedByFollowings(username string, followed bool, limit int) (
	[]*entities.TwitterProfile, error) {
	panic("not implemented") // TODO: Implement
}

func (t *Twitter) GetTopFollowedByFollowers(username string, followed bool, limit int) (
	[]*entities.TwitterProfile, error) {
	panic("not implemented") // TODO: Implement
}

func (t *Twitter) GetVerifiedFollowers(username string) ([]*entities.TwitterProfile, error) {
	panic("not implemented") // TODO: Implement
}
