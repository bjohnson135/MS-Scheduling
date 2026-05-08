package main

import (
	"embed"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"v2.staffjoy.com/account"
	"v2.staffjoy.com/apidocs"
	"v2.staffjoy.com/environments"
	"v2.staffjoy.com/healthcheck"
)

//go:embed swagger/account.swagger.json
var swaggerFS embed.FS

const (
	// ServiceName identifies this app in logs
	ServiceName = "account"
)

var (
	logger *logrus.Entry
	config environments.Config
)

// Setup environment, logger, etc
func init() {
	// Set the ENV environment variable to control dev/stage/prod behavior
	var err error
	config, err = environments.GetConfig(os.Getenv(environments.EnvVar))
	if err != nil {
		panic("Unable to determine configuration")
	}
	logger = config.GetLogger(ServiceName)
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := http.NewServeMux()

	swaggerJSON, err := swaggerFS.ReadFile("swagger/account.swagger.json")
	if err != nil {
		panic(err)
	}

	mux.HandleFunc("/swagger.json", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		res.Write(swaggerJSON)
	})
	mux.HandleFunc(healthcheck.HEALTHPATH, healthcheck.Handler)
	apidocs.Serve(mux, logger)

	// Custom runtime option to emit empty fields (like false bools)
	gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}))
	opts := []grpc.DialOption{grpc.WithInsecure()}
	errEndPoint := RegisterAccountServiceHandlerFromEndpoint(ctx, gwmux, account.Endpoint, opts)
	if errEndPoint != nil {
		return errEndPoint
	}
	mux.Handle("/", gwmux)

	return http.ListenAndServe(":80", mux)
}

func main() {
	logger.Debugf("Initialized accountapi environment %s", config.Name)

	if err := run(); err != nil {
		logger.Fatal(err)
	}
}
