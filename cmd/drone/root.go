package drone

import (
	"os"
	"superorbital/drone/utils"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const VERSION = "0.0.1"

// HTTP Client that makes it easier to mock HTTP requests
var httpClient utils.HttpClientInterface
var verbose bool

var rootCmd = &cobra.Command{
	Use:     "drone",
	Version: VERSION,
	Short:   "Drones as a service platform",
	Long: `
		RentADrone
		Drones as a service platform.
	`,
}

func init() {
	configureLogOptions()
	configureHttpClient()

	viper.SetEnvPrefix(utils.CONFIG_PREFIX)
	viper.BindEnv(utils.CONFIG_VALUE_ADDR)
	viper.BindEnv(utils.CONFIG_VALUE_TOKEN)
	viper.BindEnv(utils.CONFIG_VALUE_MAX_RETRIES)
	viper.SetDefault(utils.CONFIG_VALUE_MAX_RETRIES, 5)
}

func configureLogOptions() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

func configureHttpClient() {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = viper.GetInt(utils.CONFIG_VALUE_MAX_RETRIES)
	retryClient.HTTPClient.Timeout = utils.API_TIMEOUT * time.Second
	retryClient.Logger = utils.IntegratedLogger{}
	httpClient = retryClient.StandardClient()
}

// setLogOutput configures the log level for each command using the "verbose" flag
func setLogOutput() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

// Execute creates the CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error().Msg(err.Error())
	}
}
