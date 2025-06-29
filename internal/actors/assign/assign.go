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

package assign

import (
	"context"
	"fmt"
	"regexp"

	"github.com/google/go-github/v72/github"
	"github.com/gookit/slog"

	"github.com/ShyunnY/actbot/internal/actors"
)

const (
	assignActorName = "AssignActor"
)

var assignRegexp = regexp.MustCompile(`^/(un)?assign\b`)

type actor struct {
	ghClient *github.Client
	logger   *slog.Logger

	event github.IssueCommentEvent
	add   bool
}

func NewAssignActor(ghClient *github.Client, logger *slog.Logger, _ *actors.Options) actors.Actor {
	return &actor{
		ghClient: ghClient,
		logger:   logger,
		add:      true,
	}
}

func (a *actor) Handler() error {
	var (
		issue     = a.event.GetIssue()
		comment   = a.event.GetComment()
		loginUser = comment.GetUser()
		repo      = a.event.GetRepo()
		assignees = issue.Assignees
	)
	a.logger.Infof("actor %s started processing events, issue number: #%d", a.Name(), issue.GetNumber())

	owner, repoName := actors.GetOwnerRepo(repo.GetFullName())
	if a.add {
		// if it has been assigned to the login user, we will write back a comment
		if isAssignLoginUser(loginUser, assignees) {
			err := actors.AddComment(
				a.ghClient,
				fmt.Sprintf("@%s %s", loginUser.GetLogin(), "The issue has been assigned to you. Please do not attempt to assign it"),
				repo.GetFullName(),
				issue.GetNumber(),
			)
			return err
		}

		if _, _, err := a.ghClient.Issues.AddAssignees(
			context.Background(),
			owner,
			repoName,
			issue.GetNumber(),
			[]string{loginUser.GetLogin()},
		); err != nil {
			return err
		}
		a.logger.Infof("assigned issue to '%s'", loginUser.GetLogin())

		if err := actors.AddReaction(a.ghClient, actors.CommendReaction, repo.GetFullName(), comment.GetID()); err != nil {
			return err
		}
		a.logger.Infof("add a reaction '%s' to comment %d of issue #%d", actors.CommendReaction, comment.GetID(), issue.GetNumber())

		if err := actors.RemoveLabelToIssue(a.ghClient, repo.GetFullName(), issue.GetNumber(), actors.HelpWantedLabel); err != nil {
			return err
		}
		a.logger.Infof("remove '%s' label from issue #%d", actors.HelpWantedLabel, issue.GetNumber())
	} else {
		// if it has been unassigned to the login user, we will write back a comment
		if !isAssignLoginUser(loginUser, assignees) {
			err := actors.AddComment(
				a.ghClient,
				fmt.Sprintf("@%s %s", loginUser.GetLogin(), "This issue is no assigned to you. Please do not try to unassign it again"),
				repo.GetFullName(),
				issue.GetNumber(),
			)
			return err
		}

		if _, _, err := a.ghClient.Issues.RemoveAssignees(
			context.Background(),
			owner,
			repoName,
			issue.GetNumber(),
			[]string{loginUser.GetLogin()},
		); err != nil {
			return err
		}

		// This is left to the contributors, so don't include the 'help-wanted' tag
		// if err := actors.AddLabelToIssue(a.ghClient, repo.GetFullName(), issue.GetNumber(), actors.HelpWantedLabel); err != nil {
		//	return err
		// }
		// a.logger.Infof("add '%s' label from issue #%d", actors.HelpWantedLabel, issue.GetNumber())

		a.logger.Infof("unassigned issue to '%s'", loginUser.GetLogin())
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

	matches := assignRegexp.FindAllStringSubmatch(comment.GetBody(), -1)
	if matches == nil {
		return false
	}
	for _, match := range matches {
		if match[1] == "un" {
			a.add = false
		}
	}
	a.event = commentEvent

	return true
}

func (a *actor) Name() string {
	return assignActorName
}

func isAssignLoginUser(user *github.User, assignees []*github.User) bool {
	if len(assignees) == 0 {
		return false
	}

	for _, assignee := range assignees {
		if assignee.GetID() == user.GetID() {
			return true
		}
	}

	return false
}
