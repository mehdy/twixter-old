package mocks

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mehdy/twixter/pkg/entities"
)

const (
	DayDuration  = 24 * time.Hour
	YearDuration = 365 * DayDuration

	DefaultFollowersCount  = 10
	DefaultFollowingsCount = 10
	DefaultFavouritesCount = 10
	DefaultListedCount     = 10
	DefaultTweetsCount     = 10
)

func GenerateProfile(suffix string) *entities.TwitterProfile {
	id := uuid.New()

	return &entities.TwitterProfile{
		TwitterID:        id.String(),
		ID:               id,
		Name:             fmt.Sprintf("TEST_NAME_%s", suffix),
		Username:         fmt.Sprintf("TEST_USERNAME_%s", suffix),
		Location:         fmt.Sprintf("TEST_LOCATION_%s", suffix),
		Bio:              fmt.Sprintf("TEST_BIO_%s", suffix),
		URL:              "https://test.tld",
		Email:            "test@mail.tld",
		ProfileBannerURL: "https://test.tld/fake/banner",
		ProfileImageURL:  "https://test.tld/fake/image",
		FollowersCount:   DefaultFollowersCount,
		FollowingsCount:  DefaultFollowingsCount,
		FavouritesCount:  DefaultFavouritesCount,
		ListedCount:      DefaultListedCount,
		TweetsCount:      DefaultTweetsCount,
		Entities: map[string]interface{}{
			"TEST_KEY": map[string]string{
				"TEST_SUB_KEY": "TEST_VALUE",
			},
		},
		JoinedAt: time.Now().Add(-YearDuration),
	}
}
