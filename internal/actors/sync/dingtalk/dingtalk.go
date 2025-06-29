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

func (D *DingTalkClient) Name() string {
	return DingTalk
}

// SendMessage sends a message to the DingTalk chat group.
// In the envisionary design, we synchronize the status of issues on
// GitHub by sending a message to the DingTalk group.
// so that subsequent community-related contributors can
// pay attention to and deal with it.
// Supports sending text in markdown format.
func (D *DingTalkClient) SendMessage(issueNumber int) error {

	if D.ChatGroupRobotEndPoint == "" {
		return fmt.Errorf("chat group robot endpoint cannot be empty")
	}

	// DingTalk markdown content
	message := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": "Issue #" + strconv.Itoa(issueNumber),
			"text":  fmt.Sprintf("### Issue: [#%d](https://github.com/alibaba/spring-ai-alibaba/issues/%d), Please pay attention. 👀", issueNumber, issueNumber),
		},
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Send post request to DingTalk API
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(DefaultEndPoint, D.ChatGroupRobotEndPoint),
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
	D.logger.Debugf("Response from DingTalk: %s", respBody)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return nil

}
