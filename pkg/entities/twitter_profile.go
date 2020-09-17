package entities

import (
	"time"

	"github.com/google/uuid"
)

// TwitterProfile represents a user's profile on Twitter.
type TwitterProfile struct {
	ID                  uuid.UUID
	TwitterID           string
	Name                string
	Username            string
	Location            string
	Bio                 string
	URL                 string
	Email               string
	ProfileBannerURL    string
	ProfileImageURL     string
	Verified            bool
	Protected           bool
	DefaultProfile      bool
	DefaultProfileImage bool
	FollowersCount      int
	FollowingsCount     int
	FavouritesCount     int
	ListedCount         int
	TweetsCount         int
	Entities            map[string]interface{}
	JoinedAt            time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Followings          []TwitterProfile
	Followers           []TwitterProfile
}
