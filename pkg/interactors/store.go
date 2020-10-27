package interactors

import "github.com/mehdy/twixter/pkg/entities"

//go:generate mockgen -destination=../mocks/gomock_store.go -package=mocks . Store

// Store defines the core functionalities of a store to persist data.
type Store interface {
	// GetProfile returns the TwitterProfile of the given username.
	GetProfile(username string) (*entities.TwitterProfile, error)
	// SaveProfile stores the given TwitterProfile.
	SaveProfile(profile *entities.TwitterProfile) error
	// SaveProfiles stores a batch of TwitterProfiles.
	SaveProfiles(profiles []*entities.TwitterProfile) error
	// AddFollowings appends the list of profiles to the profile's followings
	AddFollowings(profile *entities.TwitterProfile, profiles []*entities.TwitterProfile) error
	// AddFollowers appends the list of profiles to the profile's followers
	AddFollowers(profile *entities.TwitterProfile, profiles []*entities.TwitterProfile) error
	// GetFollowings returns the list of profiles followed by the given username.
	GetFollowings(username string) ([]*entities.TwitterProfile, error)
	// GetFollowers returns the list of profiles following the given username.
	GetFollowers(username string) ([]*entities.TwitterProfile, error)
	// GetTopFollowingsByFollowers returns a limited list of TwitterProfiles followed by the given username
	// sorted by followers count.
	GetTopFollowingsByFollowers(username string, limit int) ([]*entities.TwitterProfile, error)
	// GetTopFollowersByFollowers returns a limited list of TwitterProfiles following the given username
	// sorted by followers count.
	GetTopFollowersByFollowers(username string, limit int) ([]*entities.TwitterProfile, error)
	// GetTopFollowedByFollowings returns a limited list of TwitterProfiles followed by the followings of
	// the given username filtered by whether it's followed or not.
	GetTopFollowedByFollowings(username string, followed bool, limit int) ([]*entities.TwitterProfile, error)
	// GetTopFollowedByFollowers returns a limited list of TwitterProfiles followed by the followers of
	// the given username filtered by whether it's followed or not.
	GetTopFollowedByFollowers(username string, followed bool, limit int) ([]*entities.TwitterProfile, error)
	// GetVerifiedFollowers returns a list of verified TwitterProfiles following the given username.
	GetVerifiedFollowers(username string) ([]*entities.TwitterProfile, error)
}
