// nolint: goerr113, funlen, dupl, goconst, scopelint, gocognit
package interactors_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mehdy/twixter/pkg/entities"
	"github.com/mehdy/twixter/pkg/interactors"
	"github.com/mehdy/twixter/pkg/mocks"
	"github.com/mehdy/twixter/pkg/services/logrus"
	"github.com/mehdy/twixter/pkg/services/viper"
)

func TestTiwtterIntractor(suiteT *testing.T) {
	run := func(name string, testFunc func(
		t *testing.T,
		twitter entities.TwitterInteractor,
		twitterAPIMock *mocks.MockTwitterAPI,
		storeMock *mocks.MockStore),
	) {
		config := viper.NewConfig()
		config.Set("log.level", "F")
		logger := logrus.NewLogger(config)

		suiteT.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			twitterAPIMock := mocks.NewMockTwitterAPI(ctrl)
			storeMock := mocks.NewMockStore(ctrl)
			twitter := interactors.NewTwitter(logger, twitterAPIMock, storeMock)

			testFunc(t, twitter, twitterAPIMock, storeMock)
		})
	}

	run("UpdateProfileSuccessfully", func(
		t *testing.T,
		twitter entities.TwitterInteractor,
		twitterAPIMock *mocks.MockTwitterAPI,
		storeMock *mocks.MockStore,
	) {
		profileMock := mocks.GenerateProfile("A")
		twitterAPIMock.EXPECT().Profile(gomock.Eq(profileMock.Username)).Return(profileMock, nil)
		storeMock.EXPECT().SaveProfile(gomock.Eq(profileMock)).Return(nil)

		if err := twitter.UpdateProfile(profileMock.Username); err != nil {
			t.Logf("Failed to update profile successfully: %s", err)
			t.Fail()
		}
	})

	run("UpdateProfileNonExistent", func(
		t *testing.T,
		twitter entities.TwitterInteractor,
		twitterAPIMock *mocks.MockTwitterAPI,
		storeMock *mocks.MockStore,
	) {
		username := "NonExistent"
		expectedErr := errors.New("profile not found")

		twitterAPIMock.EXPECT().Profile(gomock.Eq(username)).Return(nil, expectedErr)
		storeMock.EXPECT().SaveProfile(nil).Times(0)

		if err := twitter.UpdateProfile(username); !errors.Is(err, expectedErr) {
			t.Logf("Failed to handle non-existent profile: %s", err)
			t.Fail()
		}
	})

	run("UpdateProfileFailedStore", func(
		t *testing.T,
		twitter entities.TwitterInteractor,
		twitterAPIMock *mocks.MockTwitterAPI,
		storeMock *mocks.MockStore,
	) {
		expectedErr := errors.New("failed to save profile")

		profileMock := mocks.GenerateProfile("A")
		twitterAPIMock.EXPECT().Profile(gomock.Eq(profileMock.Username)).Return(profileMock, nil)
		storeMock.EXPECT().SaveProfile(gomock.Eq(profileMock)).Return(expectedErr)

		if err := twitter.UpdateProfile(profileMock.Username); !errors.Is(err, expectedErr) {
			t.Logf("Failed to handle failed store: %s", err)
			t.Fail()
		}
	})

	run("UpdateFollowingsSuccessfully", func(
		t *testing.T,
		twitter entities.TwitterInteractor,
		twitterAPIMock *mocks.MockTwitterAPI,
		storeMock *mocks.MockStore,
	) {
		const (
			totalProfiles = 10
			batchSize     = 2
		)

		profile := mocks.GenerateProfile("A")
		profileMocks := mocks.GenerateProfileBatches(totalProfiles, batchSize)

		storeMock.EXPECT().GetProfile(gomock.Eq(profile.Username)).Return(profile, nil)

		storeCalls := []*gomock.Call{}
		for _, profileBatch := range profileMocks {
			storeCalls = append(storeCalls,
				storeMock.EXPECT().SaveProfiles(gomock.Eq(profileBatch)),
				storeMock.EXPECT().AddFollowings(gomock.Eq(profile), gomock.Eq(profileBatch)),
			)
		}
		gomock.InOrder(storeCalls...)

		followingsCh := make(chan []*entities.TwitterProfile)
		twitterAPIMock.EXPECT().Followings(gomock.Eq(profile.Username)).Return(followingsCh, nil)
		go func() {
			for _, profileBatch := range profileMocks {
				followingsCh <- profileBatch
			}
			close(followingsCh)
		}()

		if err := twitter.UpdateFollowings(profile.Username); err != nil {
			t.Logf("Failed to update profile's followings successfully: %s", err)
			t.Fail()
		}
	})

	run("UpdateFollowingsOfNonExistentProfile", func(
		t *testing.T,
		twitter entities.TwitterInteractor,
		twitterAPIMock *mocks.MockTwitterAPI,
		storeMock *mocks.MockStore,
	) {
		username := "NonExistent"
		expectedErr := errors.New("profile not found")

		storeMock.EXPECT().GetProfile(gomock.Eq(username)).Return(nil, expectedErr)

		if err := twitter.UpdateFollowings(username); !errors.Is(err, expectedErr) {
			t.Logf("Failed to handle non-existent profile: %s", err)
			t.Fail()
		}
	})

	run("UpdateFollowingsFailedTwitterAPI", func(
		t *testing.T,
		twitter entities.TwitterInteractor,
		twitterAPIMock *mocks.MockTwitterAPI,
		storeMock *mocks.MockStore,
	) {
		profile := mocks.GenerateProfile("A")
		storeMock.EXPECT().GetProfile(gomock.Eq(profile.Username)).Return(profile, nil)

		expectedErr := errors.New("failed to fetch followings")
		twitterAPIMock.EXPECT().Followings(gomock.Eq(profile.Username)).Return(nil, expectedErr)

		if err := twitter.UpdateFollowings(profile.Username); !errors.Is(err, expectedErr) {
			t.Logf("Failed to handle failed twitter API: %s", err)
			t.Fail()
		}
	})

	run("UpdateFollowersSuccessfully", func(
		t *testing.T,
		twitter entities.TwitterInteractor,
		twitterAPIMock *mocks.MockTwitterAPI,
		storeMock *mocks.MockStore,
	) {
		const (
			totalProfiles = 10
			batchSize     = 2
		)

		profile := mocks.GenerateProfile("A")
		profileMocks := mocks.GenerateProfileBatches(totalProfiles, batchSize)

		storeMock.EXPECT().GetProfile(gomock.Eq(profile.Username)).Return(profile, nil)

		storeCalls := []*gomock.Call{}
		for _, profileBatch := range profileMocks {
			storeCalls = append(storeCalls,
				storeMock.EXPECT().SaveProfiles(gomock.Eq(profileBatch)),
				storeMock.EXPECT().AddFollowers(gomock.Eq(profile), gomock.Eq(profileBatch)),
			)
		}
		gomock.InOrder(storeCalls...)

		followingsCh := make(chan []*entities.TwitterProfile)
		twitterAPIMock.EXPECT().Followers(gomock.Eq(profile.Username)).Return(followingsCh, nil)
		go func() {
			for _, profileBatch := range profileMocks {
				followingsCh <- profileBatch
			}
			close(followingsCh)
		}()

		if err := twitter.UpdateFollowers(profile.Username); err != nil {
			t.Logf("Failed to update profile successfully: %s", err)
			t.Fail()
		}
	})

	run("UpdateFollowersOfNonExistentProfile", func(
		t *testing.T,
		twitter entities.TwitterInteractor,
		twitterAPIMock *mocks.MockTwitterAPI,
		storeMock *mocks.MockStore,
	) {
		username := "NonExistent"
		expectedErr := errors.New("profile not found")

		storeMock.EXPECT().GetProfile(gomock.Eq(username)).Return(nil, expectedErr)

		if err := twitter.UpdateFollowers(username); !errors.Is(err, expectedErr) {
			t.Logf("Failed to handle non-existent profile: %s", err)
			t.Fail()
		}
	})

	run("UpdateFollowersFailedTwitterAPI", func(
		t *testing.T,
		twitter entities.TwitterInteractor,
		twitterAPIMock *mocks.MockTwitterAPI,
		storeMock *mocks.MockStore,
	) {
		profile := mocks.GenerateProfile("A")
		storeMock.EXPECT().GetProfile(gomock.Eq(profile.Username)).Return(profile, nil)

		expectedErr := errors.New("failed to fetch Followers")
		twitterAPIMock.EXPECT().Followers(gomock.Eq(profile.Username)).Return(nil, expectedErr)

		if err := twitter.UpdateFollowers(profile.Username); !errors.Is(err, expectedErr) {
			t.Logf("Failed to handle failed twitter API: %s", err)
			t.Fail()
		}
	})

	updateNetowrkTestTable := []struct {
		title      string
		suffix     string
		followings bool
		followers  bool
		depth      int
	}{
		{"UpdateNetworkOfFollowingsDepth1Successfully", "A", true, false, 1},
		{"UpdateNetworkOfFollowingsDepth2Successfully", "A", true, false, 2},
		{"UpdateNetworkOfFollowingsDepth3Successfully", "A", true, false, 3},
		{"UpdateNetworkOfFollowersDepth1Successfully", "B", false, true, 1},
		{"UpdateNetworkOfFollowersDepth2Successfully", "B", false, true, 2},
		{"UpdateNetworkOfFollowersDepth3Successfully", "B", false, true, 3},
		{"UpdateNetworkOfBothDepth1Successfully", "C", true, true, 1},
		{"UpdateNetworkOfBothDepth2Successfully", "C", true, true, 2},
		{"UpdateNetworkOfBothDepth3Successfully", "C", true, true, 3},
	}
	for _, test := range updateNetowrkTestTable {
		run(test.title, func(
			t *testing.T,
			twitter entities.TwitterInteractor,
			twitterAPIMock *mocks.MockTwitterAPI,
			storeMock *mocks.MockStore,
		) {
			profile := mocks.GenerateProfile(test.suffix)

			if test.followings {
				followingsMap := map[*entities.TwitterProfile][]*entities.TwitterProfile{
					profile: mocks.GenerateProfileBatches(10, 1)[0],
				}
				for i := 0; i < test.depth; i++ {
					for p, followings := range followingsMap {
						setupFollowings(p, [][]*entities.TwitterProfile{followings}, storeMock, twitterAPIMock)
						storeMock.EXPECT().GetFollowings(gomock.Eq(p.Username)).Return(followingsMap[p], nil)
					}

					profiles := []*entities.TwitterProfile{}
					for _, followings := range followingsMap {
						profiles = append(profiles, followings...)
					}

					followingsMap = map[*entities.TwitterProfile][]*entities.TwitterProfile{}
					for _, p := range profiles {
						followingsMap[p] = mocks.GenerateProfileBatches(10, 1)[0]
					}
				}
			}
			if test.followers {
				followersMap := map[*entities.TwitterProfile][]*entities.TwitterProfile{
					profile: mocks.GenerateProfileBatches(10, 1)[0],
				}
				for i := 0; i < test.depth; i++ {
					for p, followers := range followersMap {
						setupFollowers(p, [][]*entities.TwitterProfile{followers}, storeMock, twitterAPIMock)
						storeMock.EXPECT().GetFollowers(gomock.Eq(p.Username)).Return(followersMap[p], nil)
					}

					profiles := []*entities.TwitterProfile{}
					for _, followers := range followersMap {
						profiles = append(profiles, followers...)
					}

					followersMap = map[*entities.TwitterProfile][]*entities.TwitterProfile{}
					for _, p := range profiles {
						followersMap[p] = mocks.GenerateProfileBatches(10, 1)[0]
					}
				}
			}
			twitter.UpdateNetwork(profile.Username, test.followings, test.followers, test.depth)
		})
	}
}

func setupFollowings(
	profile *entities.TwitterProfile,
	profiles [][]*entities.TwitterProfile,
	storeMock *mocks.MockStore,
	twitterAPIMock *mocks.MockTwitterAPI,
) {
	storeMock.EXPECT().GetProfile(gomock.Eq(profile.Username)).Return(profile, nil)

	storeCalls := []*gomock.Call{}
	for _, profileBatch := range profiles {
		storeCalls = append(storeCalls,
			storeMock.EXPECT().SaveProfiles(gomock.Eq(profileBatch)),
			storeMock.EXPECT().AddFollowings(gomock.Eq(profile), gomock.Eq(profileBatch)),
		)
	}

	gomock.InOrder(storeCalls...)

	followingsCh := make(chan []*entities.TwitterProfile)
	twitterAPIMock.EXPECT().Followings(gomock.Eq(profile.Username)).Return(followingsCh, nil)

	go func() {
		for _, profileBatch := range profiles {
			followingsCh <- profileBatch
		}

		close(followingsCh)
	}()
}

func setupFollowers(
	profile *entities.TwitterProfile,
	profiles [][]*entities.TwitterProfile,
	storeMock *mocks.MockStore,
	twitterAPIMock *mocks.MockTwitterAPI,
) {
	storeMock.EXPECT().GetProfile(gomock.Eq(profile.Username)).Return(profile, nil)

	storeCalls := []*gomock.Call{}
	for _, profileBatch := range profiles {
		storeCalls = append(storeCalls,
			storeMock.EXPECT().SaveProfiles(gomock.Eq(profileBatch)),
			storeMock.EXPECT().AddFollowers(gomock.Eq(profile), gomock.Eq(profileBatch)),
		)
	}

	gomock.InOrder(storeCalls...)

	followingsCh := make(chan []*entities.TwitterProfile)
	twitterAPIMock.EXPECT().Followers(gomock.Eq(profile.Username)).Return(followingsCh, nil)

	go func() {
		for _, profileBatch := range profiles {
			followingsCh <- profileBatch
		}

		close(followingsCh)
	}()
}
