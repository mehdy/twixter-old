// nolint: goerr113, funlen
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
		profileMock := mocks.GenerateProfile()
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

		profileMock := mocks.GenerateProfile()
		twitterAPIMock.EXPECT().Profile(gomock.Eq(profileMock.Username)).Return(profileMock, nil)
		storeMock.EXPECT().SaveProfile(gomock.Eq(profileMock)).Return(expectedErr)

		if err := twitter.UpdateProfile(profileMock.Username); !errors.Is(err, expectedErr) {
			t.Logf("Failed to handle failed store: %s", err)
			t.Fail()
		}
	})
}
