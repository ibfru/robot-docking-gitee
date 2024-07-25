package framework

import (
	"bytes"
	"community-robot-lib/config"
	"encoding/gob"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

type EventHeader struct {
	EventType    int    // dispatcher.go constants defined
	PlatformName string // "github", "gitee", "gitlab", "atomgit", "gitcode" etc
	EventName    string // "Push Hook" => PushEvent , "Issue Hook" => IssueEvent,
	// "Merge Request Hook" => PullRequestEvent, "Note Hook" => IssueCommentEvent & PullRequestCommentEvent
	EventUUID string // event id
}

// PushPayload data come from PushEvent
type PushPayload struct {
	Base string
	Head string
}

type IssuePayload struct {
	IssueNumber    string
	IssueAuthor    string
	IssueComment   string
	IssueCommenter string
	IssueLabels    []string
}

type PullRequestPayload struct {
	PRNumber    string
	PRAuthor    string
	PRComment   string
	PRCommenter string
	IssueLabels []string
}

type EventPayload struct {
	Action  string
	Org     string
	Repo    string
	HtmlURL string
	PushPayload
	IssuePayload
	PullRequestPayload
}

type GenericEvent struct {
	EventHeader
	EventPayload
	SourcePayload []byte
}

type GenericHandlerFunc func(evt *GenericEvent, cnf config.Config, lgr *logrus.Entry) error

type eventHandler struct {
	accessHandler             GenericHandlerFunc
	pushCodeBranchTagHandler  GenericHandlerFunc
	issueHandlers             GenericHandlerFunc
	pullRequestHandler        GenericHandlerFunc
	issueCommentHandler       GenericHandlerFunc
	pullRequestCommentHandler GenericHandlerFunc
	otherHandler              GenericHandlerFunc
}

type PreEventHandlerFunc func(w http.ResponseWriter, r *http.Request) *GenericEvent

type preEventHandler struct {
	reqHandler PreEventHandlerFunc
}

type postEventHandler struct {
	// expansion
}

type handlers struct {
	preEventHandler
	eventHandler
	postEventHandler
}

func (h *handlers) RegisterPreEventHandler(fn PreEventHandlerFunc) {
	h.reqHandler = fn
}

// RegisterAccessHandler registers a plugin's AccessEvent handler.
func (h *handlers) RegisterAccessHandler(fn GenericHandlerFunc) {
	h.accessHandler = fn
}

// RegisterPushCodeBranchTagHandler registers a plugin's PushEvent handler.
// source code push event、branch push/delete event、tag push/delete event
func (h *handlers) RegisterPushCodeBranchTagHandler(fn GenericHandlerFunc) {
	h.pushCodeBranchTagHandler = fn
}

// RegisterIssueHandler registers a plugin's IssueEvent handler.
// issue create/delete event、issue status change event、issue reviewer event
func (h *handlers) RegisterIssueHandler(fn GenericHandlerFunc) {
	h.issueHandlers = fn
}

// RegisterPullRequestHandler registers a plugin's PullRequestEvent handler.
// PR create/update/merge/close event、PR label create/update/delete event、PR associate(or cancel) issue event
func (h *handlers) RegisterPullRequestHandler(fn GenericHandlerFunc) {
	h.pullRequestHandler = fn
}

// RegisterIssueCommentHandler registers a plugin's IssueCommentEvent handler.
// issue comment add event
func (h *handlers) RegisterIssueCommentHandler(fn GenericHandlerFunc) {
	h.issueCommentHandler = fn
}

// RegisterPullRequestCommentHandler registers a plugin's PullRequestCommentEvent handler.
// PR comment add event
func (h *handlers) RegisterPullRequestCommentHandler(fn GenericHandlerFunc) {
	h.pullRequestCommentHandler = fn
}

// RegisterOtherHandler registers a plugin's OtherEvent handler.
func (h *handlers) RegisterOtherHandler(fn GenericHandlerFunc) {
	h.otherHandler = fn
}

func (ge *GenericEvent) CollectLogFiled() map[string]interface{} {

	m := make(map[string]interface{})

	if ge.Repo == "" {
		return m
	}

	if ge.EventName != "" {
		m["event-type"] = ge.EventName
	}
	if ge.EventUUID != "" {
		m["event-uuid"] = ge.EventUUID
	}
	if ge.Action != "" {
		m["event-action"] = ge.Action
	}
	if ge.Org != "" {
		m["org"] = ge.Org
	}
	if ge.Repo != "" {
		m["repo"] = ge.Repo
	}
	if ge.HtmlURL != "" {
		m["url"] = ge.HtmlURL
	}

	return m
}

func (ge *GenericEvent) ConvertToBytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(ge); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (ge *GenericEvent) ConvertFromBytes(b []byte) error {
	if b == nil {
		return errors.New("no data to convert")
	}

	buf := new(bytes.Buffer)
	buf.Write(b)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(ge); err != nil {
		return err
	}

	return nil
}
