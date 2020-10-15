package postgres

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type TwitterProfile struct {
	ID                  uuid.UUID `gorm:"default:uuid_generate_v1mc()"`
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
	Entities            datatypes.JSON
	FollowersCount      int
	FollowingsCount     int
	FavouritesCount     int
	ListedCount         int
	TweetsCount         int
	JoinedAt            time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
	// nolint: lll // Multiline tag is not supported
	Followings []TwitterProfile `json:"followings" gorm:"many2many:twitter_relations;foreignKey:ID;joinForeignKey:ProfileA;References:ID;joinReferences:ProfileB"`
	// nolint: lll // Multiline tag is not supported
	Followers []TwitterProfile `json:"followers" gorm:"many2many:twitter_relations;foreignKey:ID;joinForeignKey:ProfileB;References:ID;joinReferences:ProfileA"`
}
