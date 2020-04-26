// Package server exposes proxying operations
package server

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// SetupDefaults can be used to ensure fallback values as available
// if the user doesn't configure them via flags or environment variables
func SetupDefaults(proxyCmd *cobra.Command) {
	proxyCmd.PersistentFlags().Duration("headerTimeout", time.Second*30, "Time to get the first byte")
	proxyCmd.PersistentFlags().Duration("readTimeout", time.Second*30, "Time to get the first byte from body")
	proxyCmd.PersistentFlags().Duration("writeTimeout", time.Minute*2, "Time to send the first byte")
	proxyCmd.PersistentFlags().String("bind", "0.0.0.0:8080", "Address to listen for incoming requests")
	proxyCmd.PersistentFlags().String("tls", ".", "Folder to look for tls.key and tls.crt files")
	proxyCmd.PersistentFlags().String("tls-crt", "tls.crt", "Name of tls certificate")
	proxyCmd.PersistentFlags().String("tls-key", "tls.key", "Name of tls key")

	viper.SetDefault("proxy.server.headerTimeout", (time.Second * 30).String())
	viper.SetDefault("proxy.server.readTimeout", (time.Minute).String())
	viper.SetDefault("proxy.server.writeTimeout", (time.Minute * 2).String())
	viper.SetDefault("proxy.server.bind", "0.0.0.0:8080")
	viper.SetDefault("proxy.server.tls", ".")
	viper.SetDefault("proxy.server.tls.key", "tls.key")
	viper.SetDefault("proxy.server.tls.crt", "tls.crt")

	viper.BindPFlag("proxy.server.headerTimeout", proxyCmd.PersistentFlags().Lookup("headerTimeout"))
	viper.BindPFlag("proxy.server.readTimeout", proxyCmd.PersistentFlags().Lookup("readTimeout"))
	viper.BindPFlag("proxy.server.writeTimeout", proxyCmd.PersistentFlags().Lookup("writeTimeout"))
	viper.BindPFlag("proxy.server.bind", proxyCmd.PersistentFlags().Lookup("bind"))
	viper.BindPFlag("proxy.server.tls", proxyCmd.PersistentFlags().Lookup("tls"))
	viper.BindPFlag("proxy.server.tls.key", proxyCmd.PersistentFlags().Lookup("tls-key"))
	viper.BindPFlag("proxy.server.tls.crt", proxyCmd.PersistentFlags().Lookup("tls-crt"))
}

// New configures an http.Server using chain as the handler
func New(chain http.Handler) *http.Server {
	server := &http.Server{
		Handler:           chain,
		ReadHeaderTimeout: viper.GetDuration("proxy.server.headerTimeout"),
		ReadTimeout:       viper.GetDuration("proxy.server.readTimeout"),
		WriteTimeout:      viper.GetDuration("proxy.server.writeTimeout"),
		Addr:              viper.GetString("proxy.server.bind"),
	}

	return server
}

// ListenAndServeTLS reads the configuration from viper and starts a HTTPS server
func ListenAndServeTLS(s *http.Server) error {
	return s.ListenAndServeTLS(
		filepath.Join(viper.GetString("proxy.server.tls"), viper.GetString("proxy.server.tls.crt")),
		filepath.Join(viper.GetString("proxy.server.tls"), viper.GetString("proxy.server.tls.key")))
}
