package kind

import (
	"regexp"
	"strings"

	"github.com/ShyunnY/actbot/internal/actors"
	"github.com/google/go-github/v72/github"
	"github.com/gookit/slog"
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
