package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitGitHubClient(t *testing.T) {
	cases := []struct {
		caseName string
		token    string
		expect   bool
	}{
		{
			caseName: "Provide github token to initialize the github client",
			token:    "fake token",
			expect:   true,
		},
		{
			caseName: "Provide empty github token to initialize the github client",
			token:    "",
			expect:   false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.caseName, func(t *testing.T) {
			ghClient, err := InitGitHubClient(tc.token)
			if tc.expect {
				assert.NotNil(t, ghClient)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, ghClient)
				assert.Error(t, err)
			}
		})
	}
}
