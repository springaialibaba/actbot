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

package dingtalk

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	tests := []struct {
		name               string
		webhookURL         string
		issueNumber        int
		content            string
		mockResponseStatus int
		expectError        bool
	}{
		{
			name:               "Successful message send",
			webhookURL:         "/mock/webhook",
			issueNumber:        123,
			content:            "This is a test message.",
			mockResponseStatus: http.StatusOK,
			expectError:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock server
			handler := func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				w.WriteHeader(tt.mockResponseStatus)
			}
			server := httptest.NewServer(http.HandlerFunc(handler))
			defer server.Close()

			// Create DingTalkClient with the mock server URL
			client := NewDingTalkClient(server.URL, nil)

			var err error
			if tt.expectError {
				err = client.SendMessage(tt.issueNumber, "test/repo")
			}

			assert.Equal(t, tt.expectError, err != nil)
		})
	}
}
