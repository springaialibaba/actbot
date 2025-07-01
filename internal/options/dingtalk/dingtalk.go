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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gookit/slog"
)

type DingTalkClient struct {
	logger                 *slog.Logger
	ChatGroupRobotEndPoint string
}

const (
	DingTalk        = "DingTalk"
	DefaultEndPoint = "https://oapi.dingtalk.com/robot/send?access_token=%s"
)

func NewDingTalkClient(chatGroupRobotEndPoint string, logger *slog.Logger) *DingTalkClient {
	return &DingTalkClient{
		logger:                 logger,
		ChatGroupRobotEndPoint: chatGroupRobotEndPoint,
	}
}

func (dt *DingTalkClient) Name() string {
	return DingTalk
}

// SendMessage sends a message to the DingTalk chat group.
// In the envisionary design, we synchronize the status of issues on
// GitHub by sending a message to the DingTalk group.
// so that subsequent community-related contributors can
// pay attention to and deal with it.
// Supports sending text in markdown format.
func (dt *DingTalkClient) SendMessage(issueNumber int, content string) error {
	if dt.ChatGroupRobotEndPoint == "" {
		return fmt.Errorf("chat group robot endpoint cannot be empty")
	}

	if issueNumber == 0 {
		return fmt.Errorf("issue number cannot be zero")
	}

	// DingTalk markdown content
	message := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": "Issue #" + strconv.Itoa(issueNumber),
			"text":  content,
		},
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Send post request to DingTalk API
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(DefaultEndPoint, dt.ChatGroupRobotEndPoint),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}

	// Parse resp body to debug.
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	var respContent any
	err = json.Unmarshal(respBody, &respContent)
	// when debug, open it.
	// dt.logger.Debugf("Response from DingTalk: %s", respBody)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return nil
}
