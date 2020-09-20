package interactors

import "github.com/mehdy/twixter/pkg/entities"

// TwitterAPI defines core functionalities for Twitter API.
type TwitterAPI interface {
	// Profile returns the TwitterProfile of the given username.
	Profile(username string) (*entities.TwitterProfile, error)
	// Followings returns the list of TwitterProfiles followed by the given username.
	Followings(username string) ([]*entities.TwitterProfile, error)
	// Followers returns the list of TwitterProfiles following the given username.
	Followers(username string) ([]*entities.TwitterProfile, error)
}
