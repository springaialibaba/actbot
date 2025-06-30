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

package kind

import (
	"regexp"
	"strings"

	"github.com/google/go-github/v72/github"
	"github.com/gookit/slog"

	"github.com/ShyunnY/actbot/internal/actors"
)

const (
	kindLabelerActorName = "KindLabelerActor"
	kindPrefix           = "kind/"
)

var (
	kindRegexp   = regexp.MustCompile(`^/kind\s+(.+)$`)
	unkindRegexp = regexp.MustCompile(`^/unkind\s+(.+)$`)
)

type actor struct {
	ghClient *github.Client
	logger   *slog.Logger

	event github.IssueCommentEvent
}

func NewLabelerActor(ghClient *github.Client, logger *slog.Logger, _ *actors.Options) actors.Actor {
	return &actor{
		ghClient: ghClient,
		logger:   logger,
	}
}

func (a *actor) Handler() error {
	var (
		issue   = a.event.GetIssue()
		repo    = a.event.GetRepo()
		comment = a.event.GetComment()
		body    = comment.GetBody()
		err     error
	)

	if kindMatch := kindRegexp.FindStringSubmatch(body); kindMatch != nil {
		labels := strings.Fields(kindMatch[1])
		for _, label := range labels {
			label = kindPrefix + label
			err = actors.CheckAndAddLabel(a.ghClient, repo.GetFullName(), issue.GetNumber(), label)
			if err != nil {
				return err
			}
		}
	} else if unkindMatch := unkindRegexp.FindStringSubmatch(body); unkindMatch != nil {
		labels := strings.Fields(unkindMatch[1])
		for _, label := range labels {
			label = kindPrefix + label
			err = actors.RemoveLabelToIssue(a.ghClient, repo.GetFullName(), issue.GetNumber(), label)
			if err != nil {
				return err
			}
		}
	}

	// Regardless of whether it is successful or not,
	// remove the 'needs-triage' tag to prove that the issue has been handled by the maintainers.
	err = actors.RemoveLabelToIssue(a.ghClient, repo.GetFullName(), issue.GetNumber(), actors.NeedsTriageLabel)
	if err != nil {
		a.logger.Error("failed to remove 'needs-triage' label", "error", err)
		return err
	}

	return nil
}

func (a *actor) Capture(event actors.GenericEvent) bool {
	genericEvent := event.Event
	commentEvent, ok := genericEvent.(github.IssueCommentEvent)
	if !ok {
		a.logger.Error("cannot extract event to github.IssueCommentEvent, please check event type")
		return false
	}

	if commentEvent.Issue.IsPullRequest() || len(commentEvent.Comment.GetBody()) == 0 {
		return false
	}
	if commentEvent.Issue.GetClosedBy() != nil || !commentEvent.Issue.GetClosedAt().IsZero() {
		return false
	}

	if kindRegexp.MatchString(commentEvent.Comment.GetBody()) ||
		unkindRegexp.MatchString(commentEvent.Comment.GetBody()) {

		a.event = commentEvent
		return true
	}

	return false
}

func (a *actor) Name() string {
	return kindLabelerActorName
}
