package repository

import (
	"testing"

	"github.com/ilhamgepe/tengahin/internal/model"
)

func TestCreateuser(t *testing.T) {
	testCases := []struct {
		name string
		arg  model.RegisterDTO
	}{
		{
			name: "success",
			arg: model.RegisterDTO{
				Email:    "a@b.com",
				Username: "username",
				Fullname: "fullname",
				Password: "password",
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		logger.Info().Msgf("test case: %+v", tc)
	}
}
