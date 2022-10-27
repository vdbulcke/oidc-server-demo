package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
	oidcserver "github.com/vdbulcke/oidc-server-demo/oidc-server"
	"github.com/vdbulcke/oidc-server-demo/oidc-server/logger"
	"go.uber.org/zap"
)

// args var
var configFilename string
var listenAddr string
var port int
var accessLog bool

// default
var DefaultListeningAddress = "127.0.0.1"
var DefaultListeningPost = 5557

func init() {
	// bind to root command
	rootCmd.AddCommand(startCmd)
	// add flags to sub command
	startCmd.Flags().StringVarP(&configFilename, "config", "c", "", "oidc server config file")
	startCmd.Flags().StringVarP(&listenAddr, "listen-addr", "", DefaultListeningAddress, "oidc server listening address")
	startCmd.Flags().IntVarP(&port, "port", "p", DefaultListeningPost, "oidc server call back port")
	startCmd.Flags().BoolVarP(&accessLog, "access-log", "", false, "enable access log")

	// required flags
	//nolint
	startCmd.MarkFlagRequired("config")

}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the oidc server",
	// Long: "",
	Run: runServer,
}

// runServer cobra server handler
func runServer(cmd *cobra.Command, args []string) {

	// Zap Logger
	logger := logger.GetZapLogger(Debug)
	//nolint
	defer logger.Sync()

	// parse config
	config, err := oidcserver.ParseConfig(configFilename)
	if err != nil {
		logger.Error("parsing config", zap.Error(err))
		os.Exit(1)
	}

	// override config
	config.ListenAddress = listenAddr
	config.ListenPort = port

	config.AccessLog = accessLog
	config.Debug = Debug

	// validate config
	if !oidcserver.ValidateConfig(config) {
		logger.Error("validating config", zap.Error(errors.New("Validation Error")))
		os.Exit(1)
	}

	// create a new OIDC server
	s, err := oidcserver.NewOIDCServer(logger, config)
	if err != nil {
		logger.Error("creatring OIDC server", zap.Error(err))
		os.Exit(1)
	}

	err = s.StartServer()
	if err != nil {
		logger.Error("starting server", zap.Error(err))
		os.Exit(1)
	}
}
