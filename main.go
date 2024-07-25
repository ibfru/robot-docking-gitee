package main

import (
	"community-robot-lib/framework"
	liboptions "community-robot-lib/options"
	"flag"
	"fmt"
	sdkadapter "ibfu/robot-docking-gitee/gitee-adapter"
	"time"

	_ "github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"
)

type options struct {
	service liboptions.ServiceOptions
	client  liboptions.ClientOptions
}

func (o *options) Validate() error {

	if err := o.service.Validate(); err != nil {
		return err
	}

	return o.client.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options

	o.client.AddFlags(fs)
	o.service.AddFlags(fs)

	_ = fs.Parse(args)
	return o
}

func main() {
	//opt := gatherOptions(flag.NewFlagSet(os.Args[0], flag.ExitOnError), os.Args[1:]...)
	opt := options{
		service: liboptions.ServiceOptions{
			Port:         7102,
			ConfigFile:   "D:\\Project\\github\\local\\config.yaml",
			GracePeriod:  300 * time.Second,
			ReadTimeout:  120 * time.Second,
			WriteTimeout: 120 * time.Second,
			IdleTimeout:  30 * time.Minute,
		},
		client: liboptions.ClientOptions{
			TokenPath:   "D:\\Project\\github\\local\\token",
			HandlerPath: "/dd",
		},
	}
	if err := opt.Validate(); err != nil {
		logrus.WithError(err).Fatal("Configuration invalid: " + err.Error())
		return
	}

	//secretAgent := new(secret.Agent)
	//if err := secretAgent.Start([]string{opt.client.TokenPath}); err != nil {
	//	logrus.WithError(err).Fatal("Error starting secret agent.")
	//}

	//defer secretAgent.Stop()

	//bot := newRobot(sdkadapter.GetClientInstance(secretAgent.GetSecret(opt.client.TokenPath)))

	//agent := config.NewConfigAgent(bot.NewConfig)
	//if err := agent.Start(opt.service.ConfigFile); err != nil {
	//	logrus.WithError(err).Errorf("start config:%s", opt.service.ConfigFile)
	//	return
	//}

	//defer interrupts.WaitForGracefulShutdown()
	//
	//interrupts.OnInterrupt(func() {
	//	agent.Stop()
	//})
	//
	//httpServer := &http.Server{
	//	Addr:         ":" + strconv.Itoa(opt.service.Port),
	//	Handler:      bot.setupRouter(),
	//	ReadTimeout:  opt.service.ReadTimeout,
	//	WriteTimeout: opt.service.WriteTimeout,
	//	IdleTimeout:  opt.service.IdleTimeout,
	//}

	bot := newRobot(sdkadapter.GetClientInstance([]byte("123")))
	opt.client.Handler = bot.setupRouter()
	fmt.Printf("\u001B[0;31;6m %s%d \u001B[0;30;6m \n", "++++++++++++++++++++++++",
		windows.GetCurrentProcessId())
	framework.Run(bot, opt.service, opt.client)

	fmt.Printf("\u001B[0;31;6m %s%d \u001B[0;30;6m \n", "=========================",
		windows.GetCurrentProcessId())
	//interrupts.ListenAndServe(httpServer, opt.service.GracePeriod)

}
