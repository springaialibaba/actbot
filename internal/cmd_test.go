// Copyright 2024-2025 the original author or authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
