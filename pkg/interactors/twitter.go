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

func (t *Twitter) updateProfiles(
	profile *entities.TwitterProfile,
	listOfProfilesCh chan []*entities.TwitterProfile,
	following bool,
) (int, int) {
	successfulCount, failedCount := 0, 0

	for listOfProfiles := range listOfProfilesCh {
		if err := t.store.SaveProfiles(listOfProfiles); err != nil {
			t.logger.As("W").WithError(err).Logf("Failed to update a batch of profiles")
			failedCount++

			continue
		}

		var err error
		if following {
			err = t.store.AddFollowings(profile, listOfProfiles)
		} else {
			err = t.store.AddFollowers(profile, listOfProfiles)
		}

		if err != nil {
			t.logger.As("W").WithError(err).Logf("Failed to update the relations")
			failedCount++

			continue
		}

		t.logger.As("D").Logf("Updated a batch of profiles successfully")
		successfulCount++
	}

	return successfulCount, failedCount
}

func (t *Twitter) UpdateFollowings(username string) error {
	profile, err := t.store.GetProfile(username)
	if err != nil {
		t.logger.As("E").WithError(err).WithField("username", username).Logf("Failed to fetch profile from store")

		return newError(err, "failed to fetch profile")
	}

	followingsCh, err := t.twitter.Followings(username)
	if err != nil {
		t.logger.As("E").WithError(err).WithField("username", username).
			Logf("Failed to fetch followings from the twitter API")

		return newError(err, "failed to fetch followings")
	}

	t.logger.As("D").WithField("username", username).Logf("Fetched the followings from the twitter API successfully")

	successfulCount, failedCount := t.updateProfiles(profile, followingsCh, true)
	t.logger.As("I").
		WithField("username", username).
		WithField("successful", successfulCount).
		WithField("failed", failedCount).
		Logf("Updated profile followings")

	return nil
}

func (t *Twitter) UpdateFollowers(username string) error {
	profile, err := t.store.GetProfile(username)
	if err != nil {
		t.logger.As("E").WithError(err).WithField("username", username).Logf("Failed to fetch profile from store")

		return newError(err, "failed to fetch profile")
	}

	followersCh, err := t.twitter.Followers(username)
	if err != nil {
		t.logger.As("E").WithError(err).WithField("username", username).
			Logf("Failed to fetch followers from the twitter API")

		return newError(err, "failed to fetch followers")
	}

	t.logger.As("D").WithField("username", username).Logf("Fetched the followers from the twitter API successfully")

	successfulCount, failedCount := t.updateProfiles(profile, followersCh, false)
	t.logger.As("I").
		WithField("username", username).
		WithField("successful", successfulCount).
		WithField("failed", failedCount).
		Logf("Updated profile followers")

	return nil
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

func (t *Twitter) updateFollowingNetwork(username string, depth int) {
	usernames := []string{username}

	for i := 0; i < depth; i++ {
		for _, un := range usernames {
			if err := t.UpdateFollowings(un); err != nil {
				t.logger.As("W").
					WithError(err).
					WithField("username", un).
					WithField("depth", i).
					Logf("Failed to update followings")
			} else {
				t.logger.As("D").
					WithField("username", un).
					WithField("depth", i).
					Logf("Updated followings")
			}

			followings, err := t.store.GetFollowings(un)
			if err != nil {
				t.logger.As("W").
					WithError(err).
					WithField("username", username).
					WithField("depth", i).
					Logf("Failed to get followings from store")

				return
			}

			usernames = []string{}
			for _, profile := range followings {
				usernames = append(usernames, profile.Username)
			}
		}
	}
}

func (t *Twitter) updateFollowerNetwork(username string, depth int) {
	usernames := []string{username}

	for i := 0; i < depth; i++ {
		for _, un := range usernames {
			if err := t.UpdateFollowers(un); err != nil {
				t.logger.As("W").
					WithError(err).
					WithField("username", un).
					WithField("depth", i).
					Logf("Failed to update followers")
			} else {
				t.logger.As("D").
					WithField("username", un).
					WithField("depth", i).
					Logf("Updated followings")
			}

			followings, err := t.store.GetFollowers(un)
			if err != nil {
				t.logger.As("W").
					WithError(err).
					WithField("username", username).
					WithField("depth", i).
					Logf("Failed to get followers from store")

				return
			}

			usernames = []string{}
			for _, profile := range followings {
				usernames = append(usernames, profile.Username)
			}
		}
	}
}

func (t *Twitter) UpdateNetwork(username string, followings bool, followers bool, depth int) {
	if followings {
		t.updateFollowingNetwork(username, depth)
	}

	if followers {
		t.updateFollowerNetwork(username, depth)
	}
}

func (t *Twitter) GetTopFollowingsByFollowers(username string, limit int) ([]*entities.TwitterProfile, error) {
	profiles, err := t.store.GetTopFollowingsByFollowers(username, limit)
	if err != nil {
		t.logger.As("E").
			WithError(err).
			WithField("username", username).
			WithField("limit", limit).
			Logf("Failed to get top followings by followers")
		return nil, newError(err, "failed to get top followings by followers")
	}
	return profiles, nil
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
