package cmd

import (
	"log"
	"os"

	"github.com/fabianMendez/jarlib/core"
	"github.com/spf13/cobra"
)

var (
	dependency  string
	output      string
	javaVersion string
)

func init() {
	rootCmd.Flags().StringVarP(&dependency, "dependency", "d", "", "Dependency to generate")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "Output file")
	rootCmd.Flags().StringVarP(&javaVersion, "java-version", "j", "1.8", "Java version to use")

	rootCmd.MarkFlagRequired("dependency")
	rootCmd.MarkFlagRequired("output")
}

var rootCmd = &cobra.Command{
	Use:   "jarlib",
	Short: "Jarlib is a simple tool to generate a self-contained jar of a dependency",
	Long:  `A simple tool to generate a self-contained jar of a dependency`,
	RunE: func(cmd *cobra.Command, args []string) error {
		out, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return err
		}

		err = core.Generate(dependency, javaVersion, out)
		if err != nil {
			return err
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
