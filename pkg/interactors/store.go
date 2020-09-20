package interactors

import "github.com/mehdy/twixter/pkg/entities"

//go:generate mockgen -destination=../mocks/gomock_store.go -package=mocks . Store

// Store defines the core functionalities of a store to persist data.
type Store interface {
	// GetProfile returns the TwitterProfile of the given username.
	GetProfile(username string) (*entities.TwitterProfile, error)
	// SaveProfile stores the given TwitterProfile.
	SaveProfile(profile *entities.TwitterProfile) error
}
