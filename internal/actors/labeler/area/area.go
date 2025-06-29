package area

import (
	"regexp"
	"strings"

	"github.com/google/go-github/v72/github"
	"github.com/gookit/slog"

	"github.com/ShyunnY/actbot/internal/actors"
)

const (
	areaLabelerActorName = "AreaLabelerActor"
	areaPrefix           = "area/"
)

var (
	areaRegexp   = regexp.MustCompile(`^/area\s+(.+)$`)
	unareaRegexp = regexp.MustCompile(`^/unarea\s+(.+)$`)
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
	)
	a.logger.Infof("actor %s started processing events, issue number: #%d", a.Name(), issue.GetNumber())

	var err error
	if areaMatch := areaRegexp.FindStringSubmatch(body); areaMatch != nil {
		labels := strings.Fields(areaMatch[1])
		for _, label := range labels {
			label = areaPrefix + label
			err = actors.CheckAndAddLabel(a.ghClient, repo.GetFullName(), issue.GetNumber(), label)
			if err != nil {
				return err
			}
		}
	} else if unareaMatch := unareaRegexp.FindStringSubmatch(body); unareaMatch != nil {
		labels := strings.Fields(unareaMatch[1])
		for _, label := range labels {
			label = areaPrefix + label
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

	if areaRegexp.MatchString(commentEvent.Comment.GetBody()) ||
		unareaRegexp.MatchString(commentEvent.Comment.GetBody()) {

		a.event = commentEvent
		return true
	}

	return false
}

func (a *actor) Name() string {
	return areaLabelerActorName
}
