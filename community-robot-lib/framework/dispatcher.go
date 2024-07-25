package framework

import (
	"net/http"
	"sync"

	"community-robot-lib/config"

	"github.com/sirupsen/logrus"
)

const (
	UserAgentHeader = "Robot-Gateway-Access"
)

type dispatcher struct {
	agent *config.ConfigAgent

	h handlers

	// Tracks running handlers for graceful shutdown
	wg sync.WaitGroup

	// secret usage
	hmac func() []byte
}

func (d *dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ge := d.h.reqHandler(w, r)
	if ge == nil {
		return
	}

	lgr := logrus.WithFields(ge.CollectLogFiled())

	d.Dispatch(ge, lgr)
}

func (d *dispatcher) Dispatch(event *GenericEvent, lgr *logrus.Entry) {
	if event.EventType < AccessEvent || event.EventType > OtherEvent {
		lgr.Error("Ignoring unknown event type")
	} else {
		d.wg.Add(1)
		go d.handleEvent(event, lgr)
	}
}

func (d *dispatcher) Wait() {
	d.wg.Wait() // Handle remaining requests
}

var eventHandlerList []GenericHandlerFunc
var once sync.Once

// Event-Type Value
const (
	AccessEvent = iota
	PushEvent
	IssueEvent
	PullRequestEvent
	IssueCommentEvent
	PullRequestCommentEvent
	OtherEvent
)

func (h *handlers) indexEventHandler() {
	// slice element order must same to Event-Type Value
	eventHandlerList = []GenericHandlerFunc{
		h.accessHandler,
		h.pushCodeBranchTagHandler,
		h.issueHandlers,
		h.pullRequestHandler,
		h.issueCommentHandler,
		h.pullRequestCommentHandler,
		h.otherHandler,
	}
}

func buildDispatcherHandler(h *handlers) {
	once.Do(h.indexEventHandler)
}

func (d *dispatcher) getConfig() config.Config {
	_, c := d.agent.GetConfig()

	return c
}

// handleAccessEvent access robot handle request that come form webhook
func (d *dispatcher) handleEvent(evt *GenericEvent, lgr *logrus.Entry) {
	defer d.wg.Done()

	fn := eventHandlerList[evt.EventType]
	if err := fn(evt, d.getConfig(), lgr); err != nil {
		lgr.Error(err)
	} else {
		lgr.Info()
	}
}
