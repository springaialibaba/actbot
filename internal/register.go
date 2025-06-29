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
	"github.com/google/go-github/v72/github"
	"github.com/gookit/slog"

	"github.com/ShyunnY/actbot/internal/actors"
	"github.com/ShyunnY/actbot/internal/actors/assign"
	"github.com/ShyunnY/actbot/internal/actors/labeler"
	"github.com/ShyunnY/actbot/internal/actors/retest"
	"github.com/ShyunnY/actbot/internal/actors/sync"
)

type GitHubEventType string

type RegisterFn = func(ghClient *github.Client, logger *slog.Logger, opts *actors.Options) actors.Actor

const (
	IssueComment GitHubEventType = "issue_comment"
)

var actorMap = map[GitHubEventType][]RegisterFn{
	IssueComment: {
		assign.NewAssignActor,
		retest.NewRetestActor,
		sync.NewSyncActor,
		labeler.NewLabelerActor,
	},
}
