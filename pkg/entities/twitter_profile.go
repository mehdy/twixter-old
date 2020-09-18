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

// TwitterIntractor defines the core functionalities on TwitterProfile.
type TwitterInteractor interface {
	// UpdateFollowings updates the followings of the given username.
	UpdateFollowings(username string) error
	// UpdateFollowers updates the followers of the given username.
	UpdateFollowers(username string) error
	// UpdateProfile updates the TwitterProfile of the given username.
	UpdateProfile(username string) error
	// UpdateNetwork updates the followings and followers of the given username as requested
	// and theirs as well until the given depth.
	UpdateNetwork(username string, followings, followers bool, depth int) error

	// GetTopFollowingsByFollowers returns a limited list of TwitterProfiles followed by the given username
	// sorted by followers count.
	GetTopFollowingsByFollowers(username string, limit int) ([]*TwitterProfile, error)
	// GetTopFollowersByFollowers returns a limited list of TwitterProfiles following the given username
	// sorted by followers count.
	GetTopFollowersByFollowers(username string, limit int) ([]*TwitterProfile, error)
	// GetTopFollowedByFollowings returns a limited list of TwitterProfiles followed by the followings of
	// the given username filtered by whether it's followed or not.
	GetTopFollowedByFollowings(username string, followed bool, limit int) ([]*TwitterProfile, error)
	// GetTopFollowedByFollowers returns a limited list of TwitterProfiles followed by the followers of
	// the given username filtered by whether it's followed or not.
	GetTopFollowedByFollowers(username string, followed bool, limit int) ([]*TwitterProfile, error)

	// GetVerifiedFollowers returns a list of verified TwitterProfiles following the given username.
	GetVerifiedFollowers(username string) ([]*TwitterProfile, error)
}
