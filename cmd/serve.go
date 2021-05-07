package cmd

import (
	"log"
	"net/http"
	"os"

	"github.com/fabianMendez/jarlib/handlers"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run an http server",
	Long:  `Run an http server which generates a self-contained jar of a dependency`,
	Run: func(cmd *cobra.Command, args []string) {
		http.HandleFunc("/", handlers.Generate())

		port := os.Getenv("PORT")
		if port == "" {
			port = "9090"
		}

		log.Printf("Listening on port %s\n", port)

		if err := http.ListenAndServe(":"+port, nil); err != nil {
			panic(err)
		}
	},
}
