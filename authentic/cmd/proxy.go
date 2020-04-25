package cmd

import (
	"github.com/spf13/cobra"
)

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Proxy requests to the given target if the required cookies are present",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)

	proxyCmd.PersistentFlags().String("bind", "0.0.0.0:8080", "Address to listen for incoming requests")
	proxyCmd.PersistentFlags().String("cookieName", "authenticated", "Cookie which should be sent by the server")
	proxyCmd.PersistentFlags().String("target", "http://localhost:8081/", "URL to send requests after they are authenticated")
}
