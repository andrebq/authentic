package cmd

import (
	"log"
	"net/http"
	"net/url"

	"github.com/andrebq/authentic/auth"
	"github.com/andrebq/authentic/internal/firebase"
	"github.com/andrebq/authentic/internal/session"
	"github.com/andrebq/authentic/internal/tcache"
	"github.com/andrebq/authentic/proxy"
	"github.com/andrebq/authentic/server"
	"github.com/spf13/cobra"
)

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Proxy requests to the given target if the required cookies are present",
	Run: func(cmd *cobra.Command, args []string) {
		target, err := url.Parse(cmd.Flag("target").Value.String())
		if err != nil {
			log.Fatal(err)
		}

		catalog, err := firebase.Users()
		if err != nil {
			log.Fatal(err)
		}

		redis, err := tcache.NewRedis("localhost:6379")
		if err != nil {
			log.Fatal(err)
		}

		s, err := session.New(redis, nil)
		if err != nil {
			panic(err)
		}

		authServer := auth.New("/auth/", s, catalog)

		proxyServer := proxy.NewReverse(
			cmd.Flag("cookieName").Value.String(),
			cmd.Flag("realm").Value.String(),
			s,
			target)

		mux := http.NewServeMux()
		mux.Handle("/auth/", authServer)
		mux.Handle("/", proxyServer)

		srv := server.New(mux)

		log.Printf("Starting authentic proxy: %v", srv.Addr)
		err = server.ListenAndServeTLS(srv)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)

	proxyCmd.PersistentFlags().String("cookieName", "authenticated", "Cookie which should be sent by the server")
	proxyCmd.PersistentFlags().String("realm", "Secure", "Realm to use for WWW-Authenticate")
	proxyCmd.PersistentFlags().String("target", "http://localhost:8081/", "URL to send requests after they are authenticated")

	server.SetupDefaults(proxyCmd)
}
