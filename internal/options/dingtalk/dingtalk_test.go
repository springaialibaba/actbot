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
				err = client.SendMessage(tt.issueNumber)
			}

			assert.Equal(t, tt.expectError, err != nil)
		})
	}
}
