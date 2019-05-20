package cmd

import (
	"os"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/go-kit/kit/log"
	"github.com/ghostsquad/go-timejumper"
	"github.com/oklog/run"
)

var rootCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
				  love by spf13 and friends in Go.
				  Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
	  // Do Stuff Here
	},
  }
  
  func Execute() {
	if err := rootCmd.Execute(); err != nil {
	  fmt.Println(err)
	  os.Exit(1)
	}
  }