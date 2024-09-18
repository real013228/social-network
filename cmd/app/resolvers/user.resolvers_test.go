package resolvers

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/real013228/social-network/internal/model"
	mock_resolvers "github.com/real013228/social-network/internal/services/mock"
	"testing"
)

func TestMutationResolver_CreateUser(t *testing.T) {
	type mockBehaviour func(s *mock_resolvers.MockuserService, user model.User)

	testTable := []struct {
		name               string
		inputBody          string
		inputUser          model.User
		mockBehaviour      mockBehaviour
		expectedResultBody string
	}{
		{
			name:      "okay",
			inputBody: `{"username":"Test", "email":"test@mail.ru"}`,
			inputUser: model.User{
				Username: "Test",
				Email:    "test@mail.ru",
			},
			mockBehaviour: func(s *mock_resolvers.MockuserService, user model.User) {
				s.EXPECT().CreateUser(context.TODO(), user)
			},
			expectedResultBody: `{"username":"Test","email":"test@mail.ru"}`,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			service := mock_resolvers.NewMockuserService(ctl)
			got, err := service.CreateUser(context.Background(), model.CreateUserInput{
				Username: tt.inputUser.Username,
				Email:    tt.inputUser.Email,
			})
			if err != nil {
				t.Error(err)
			}
			if _, err := uuid.Parse(got); err != nil {
				t.Error(err)
			}
		})
	}
}
