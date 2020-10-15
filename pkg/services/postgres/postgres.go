package postgres

import (
	"encoding/json"
	"fmt"

	"github.com/mehdy/twixter/pkg/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Twitter struct {
	db     *gorm.DB
	logger entities.Logger
}

func NewTwitter(config entities.ConfigGetter, logger entities.Logger) *Twitter {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		config.GetString("POSTGRES_USER"),
		config.GetString("POSTGRES_PASSWORD"),
		config.GetString("POSTGRES_HOST"),
		config.GetInt("POSTGRES_PORT"),
		config.GetString("POSTGRES_DB"),
		config.GetString("POSTGRES_SSLMODE"),
	)

	logger.As("I").Logf(dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.As("F").WithError(err).Logf("Failed to connect to postgres db")
	}

	return &Twitter{
		db:     db,
		logger: logger,
	}
}

func (t *Twitter) GetProfile(username string) (*entities.TwitterProfile, error) {
	tp := &TwitterProfile{}
	if err := t.db.Where(&TwitterProfile{Username: username}).First(tp).Error; err != nil {
		t.logger.As("E").WithError(err).WithField("username", username).Logf("Failed to get Profile from database")

		return nil, newError(err, "failed to fetch profile from database")
	}

	return t.asTwitterProfile(tp), nil
}

func (t *Twitter) SaveProfile(profile *entities.TwitterProfile) error {
	tp := t.fromTwitterProfile(profile)

	if err := t.db.Create(tp).Error; err != nil {
		t.logger.As("E").WithError(err).WithField("username", profile.Username).Logf("Failed to save Profile")

		return newError(err, "failed to save profile")
	}

	return nil
}

func (t *Twitter) fromTwitterProfile(profile *entities.TwitterProfile) *TwitterProfile {
	entitesJSON, err := json.Marshal(profile.Entities)
	if err != nil {
		t.logger.As("W").WithError(err).Logf("Failed to serialize profile.Entities")
	}

	return &TwitterProfile{
		TwitterID:           profile.TwitterID,
		Name:                profile.Name,
		Username:            profile.Username,
		Location:            profile.Location,
		Bio:                 profile.Bio,
		URL:                 profile.URL,
		Email:               profile.Email,
		ProfileBannerURL:    profile.ProfileBannerURL,
		ProfileImageURL:     profile.ProfileImageURL,
		Verified:            profile.Verified,
		Protected:           profile.Protected,
		DefaultProfile:      profile.DefaultProfile,
		DefaultProfileImage: profile.DefaultProfileImage,
		FollowersCount:      profile.FollowersCount,
		FollowingsCount:     profile.FollowingsCount,
		FavouritesCount:     profile.FavouritesCount,
		ListedCount:         profile.ListedCount,
		TweetsCount:         profile.TweetsCount,
		Entities:            entitesJSON,
		JoinedAt:            profile.JoinedAt,
	}
}

func (t *Twitter) asTwitterProfile(profile *TwitterProfile) *entities.TwitterProfile {
	var ent map[string]interface{}

	err := json.Unmarshal(profile.Entities, &ent)
	if err != nil {
		t.logger.As("W").WithError(err).Logf("Failed to serialize profile.Entities")
	}

	return &entities.TwitterProfile{
		TwitterID:           profile.TwitterID,
		Name:                profile.Name,
		Username:            profile.Username,
		Location:            profile.Location,
		Bio:                 profile.Bio,
		URL:                 profile.URL,
		Email:               profile.Email,
		ProfileBannerURL:    profile.ProfileBannerURL,
		ProfileImageURL:     profile.ProfileImageURL,
		Verified:            profile.Verified,
		Protected:           profile.Protected,
		DefaultProfile:      profile.DefaultProfile,
		DefaultProfileImage: profile.DefaultProfileImage,
		FollowersCount:      profile.FollowersCount,
		FollowingsCount:     profile.FollowingsCount,
		FavouritesCount:     profile.FavouritesCount,
		ListedCount:         profile.ListedCount,
		TweetsCount:         profile.TweetsCount,
		Entities:            ent,
		JoinedAt:            profile.JoinedAt,
	}
}
