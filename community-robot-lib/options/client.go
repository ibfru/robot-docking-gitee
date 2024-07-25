package options

import (
	"flag"
	"net/http"
)

// ClientOptions holds options for interacting with Client.
type ClientOptions struct {
	TokenPath       string
	TokenGenerator  func() []byte
	RepoCacheDir    string
	CacheRepoOnPV   bool
	HandlerPath     string
	CacheEndpoint   string
	CacheMaxRetries int
	Handler         http.Handler
}

// NewClientOptions creates a ClientOptions with default values.
func NewClientOptions() *ClientOptions {
	return &ClientOptions{}
}

// AddFlags injects Client options into the given FlagSet.
func (o *ClientOptions) AddFlags(fs *flag.FlagSet) {
	o.addFlags("/etc/Client/oauth", fs)
}

// AddFlagsWithoutDefaultClientTokenPath injects Client options into the given
// Flagset without setting a default for for the ClientTokenPath, allowing to
// use an anonymous Gitee client
func (o *ClientOptions) AddFlagsWithoutDefaultClientTokenPath(fs *flag.FlagSet) {
	o.addFlags("", fs)
}

func (o *ClientOptions) addFlags(defaultClientTokenPath string, fs *flag.FlagSet) {
	fs.StringVar(
		&o.TokenPath,
		"Client-token-path",
		defaultClientTokenPath,
		"Path to the file containing the Client OAuth secret.",
	)
}

// Validate validates Client options.
func (o *ClientOptions) Validate() error {
	return nil
}
