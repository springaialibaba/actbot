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

package sync

import (
	"regexp"

	"github.com/google/go-github/v72/github"
	"github.com/gookit/slog"

	"github.com/ShyunnY/actbot/internal/actors"
	"github.com/ShyunnY/actbot/internal/options/dingtalk"
)

const (
	syncActorName = "SyncActor"

	// Sync Label: GitHub issues that have been synced
	// to the DingTalk group will be marked with this label.
	syncLabel = "sync"
)

type actor struct {
	ghClient *github.Client
	logger   *slog.Logger

	// DingTalk Client
	dingTalk *dingtalk.DingTalkClient

	// event is the GitHub issue comment event that triggered this actor.
	event github.IssueCommentEvent
}

var syncRegexp = regexp.MustCompile(`^/sync\s*$`)

func NewSyncActor(ghClient *github.Client, logger *slog.Logger, opts *actors.Options) actors.Actor {
	return &actor{
		dingTalk: opts.DingTalkClient,
		ghClient: ghClient,
		logger:   logger,
	}
}

func (a *actor) Handler() error {
	var (
		repo  = a.event.GetRepo()
		issue = a.event.GetIssue()
	)
	a.logger.Infof("actor %s started processing events, issue number: #%d", a.Name(), repo.GetFullName(), issue.GetNumber())

	// check if the issue is already labeled with syncLabel, return.
	err, has := actors.HasLabel(a.ghClient, repo.GetFullName(), syncLabel, issue.GetNumber())
	if err != nil {
		a.logger.Infof("failed to check if issue #%d has label %s, err: %v", issue.GetNumber(), syncLabel, err)
		return err
	}

	if has {
		a.logger.Infof("issue #%d has label %s, skip sending message", issue.GetNumber(), syncLabel)
		return nil
	}

	// send msg
	if err := a.dingTalk.SendMessage(issue.GetNumber(), repo.GetFullName()); err != nil {
		a.logger.Errorf("failed to send message to DingTalk by err: %v", err)
		return err
	}

	// Add sync label to the issue
	err := actors.AddLabelToIssue(a.ghClient, repo.GetFullName(), issue.GetNumber(), syncLabel)
	a.logger.Warnf("add label %s to issue #%d, err: %v", syncLabel, issue.GetNumber(), err)
	if err != nil {
		return err
	}

	return nil
}

// Capture checks if the event is a GitHub issue comment `/sync` event.
// If it is, exec handler func.
func (a *actor) Capture(event actors.GenericEvent) bool {
	// Get issue comment event and type check.
	commentEvent, ok := event.Event.(github.IssueCommentEvent)
	if !ok {
		a.logger.Debugf("event %T is not a github.IssueCommentEvent", event.Event)
		return false
	}

	// pull request is essentially an issue, and the current actor does not handle this situation.
	if commentEvent.Issue.IsPullRequest() {
		return false
	}

	// do not handle closed issues
	if !commentEvent.Issue.GetClosedAt().IsZero() || commentEvent.Issue.ClosedBy != nil {
		return false
	}

	comment := commentEvent.GetComment()
	if comment == nil || len(comment.GetBody()) == 0 {
		return false
	}

	// Check if the comment body matches the `/sync` command.
	matches := syncRegexp.FindAllStringSubmatch(comment.GetBody(), -1)
	if matches == nil {
		// the comment does not match the `/sync` command.
		return false
	}

	// If the command is `/sync`, set the event msg.
	a.event = commentEvent

	return true
}

func (a *actor) Name() string {
	return syncActorName
}
