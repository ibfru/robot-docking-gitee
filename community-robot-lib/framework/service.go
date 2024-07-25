package framework

import (
	"community-robot-lib/options"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

	"community-robot-lib/config"
	"community-robot-lib/interrupts"
)

type HandlerRegister interface {
	RegisterPreEventHandler(PreEventHandlerFunc)
	RegisterAccessHandler(GenericHandlerFunc)
	RegisterPushCodeBranchTagHandler(GenericHandlerFunc)
	RegisterIssueHandler(GenericHandlerFunc)
	RegisterPullRequestHandler(GenericHandlerFunc)
	RegisterIssueCommentHandler(GenericHandlerFunc)
	RegisterPullRequestCommentHandler(GenericHandlerFunc)
	RegisterOtherHandler(GenericHandlerFunc)
}

type Robot interface {
	NewConfig() config.Config
	RegisterEventHandler(HandlerRegister)
}

func Run(bot Robot, servOpt options.ServiceOptions, clientOpt options.ClientOptions) {
	agent := config.NewConfigAgent(bot.NewConfig)
	if err := agent.Start(servOpt.ConfigFile); err != nil {
		logrus.WithError(err).Errorf("start config:%s", servOpt.ConfigFile)
		return
	}

	defer interrupts.WaitForGracefulShutdown()

	// dispatcher not used, custom handle request
	if clientOpt.Handler == nil {
		h := handlers{}
		bot.RegisterEventHandler(&h)
		buildDispatcherHandler(&h)
		d := &dispatcher{agent: &agent, h: h, hmac: clientOpt.TokenGenerator}

		interrupts.OnInterrupt(func() {
			agent.Stop()
			d.Wait()
		})

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// service's healthy check, do nothing
		})

		http.Handle(clientOpt.HandlerPath, d)
	} else {
		interrupts.OnInterrupt(func() {
			agent.Stop()
		})
	}

	httpServer := &http.Server{
		Addr:         ":" + strconv.Itoa(servOpt.Port),
		Handler:      clientOpt.Handler,
		ReadTimeout:  servOpt.ReadTimeout,
		WriteTimeout: servOpt.WriteTimeout,
		IdleTimeout:  servOpt.IdleTimeout,
	}

	interrupts.ListenAndServe(httpServer, servOpt.GracePeriod)
}
