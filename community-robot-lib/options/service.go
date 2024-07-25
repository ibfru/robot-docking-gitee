package options

import (
	"flag"
	"fmt"
	"time"
)

type ServiceOptions struct {
	Port         int
	ConfigFile   string
	GracePeriod  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func (o *ServiceOptions) Validate() error {
	if o.ConfigFile == "" {
		return fmt.Errorf("missing config-file")
	}

	return nil
}

func (o *ServiceOptions) AddFlags(fs *flag.FlagSet) {
	fs.IntVar(&o.Port, "port", 8888, "Port to listen on.")
	fs.StringVar(&o.ConfigFile, "config-file", "", "Path to config file.")
	fs.DurationVar(&o.GracePeriod, "grace-period", 180*time.Second, "On shutdown, try to handle remaining events for the specified duration.")
	fs.DurationVar(&o.ReadTimeout, "read-timeout", 180*time.Second, "the maximum duration for reading the entire request, including the body")
	fs.DurationVar(&o.WriteTimeout, "write-timeout", 180*time.Second, "the maximum duration before timing out writes of the response")
	fs.DurationVar(&o.IdleTimeout, "idle-timeout", 30*time.Minute, "the maximum amount of time to wait for the next request when keep-alives are enabled")
}
