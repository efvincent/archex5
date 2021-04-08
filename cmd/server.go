package cmd

import (
	"github.com/efvincent/archex5/API"
	"github.com/spf13/cobra"
)

// holds the port & host config parameters (see init())
var port string
var host string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the API",
	Long: `Starts the HTTP API on the specificed port (defaults to 8080).
		
Note the server blocks the process. Press CTRL-C to stop the server running`,
	Run: func(cmd *cobra.Command, args []string) {
		API.Run(host, port)
	},
}

func init() {
	serverCmd.Flags().StringVar(&host, "host", "localhost", "The HTTP Host for the API.")
	serverCmd.Flags().StringVar(&port, "port", "8080", "The HTTP Port for the API.")
	rootCmd.AddCommand(serverCmd)
}
