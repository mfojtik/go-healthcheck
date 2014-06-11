package main

import (
	"fmt"
	"log"
	"os"

	"github.com/openshift/go-healthcheck/api"
	"github.com/openshift/go-healthcheck/plugins"
	"github.com/spf13/cobra"
)

func main() {

	var enabledPlugins string

	repo := api.Repository{}

	repo.Add(plugins.HttpPlugin{})
	repo.Add(plugins.MongoPlugin{})
	repo.Add(plugins.TcpPlugin{})
	repo.Add(plugins.FilePlugin{})

	req := &api.StatusRequest{}

	mainCmd := &cobra.Command{
		Use:   "healthchk",
		Short: "A simple tool for health checking Docker containers",
		Run: func(c *cobra.Command, args []string) {
			c.Usage()
		},
	}

	mainCmd.PersistentFlags().BoolVar((&req.Verbose), "verbose", false, "Enable verbose output")
	mainCmd.PersistentFlags().StringVar((&req.Socket), "socket", "unix:///var/run/docker.sock", "Docker socket to use")

	statusCmd := &cobra.Command{
		Use:   "status CONTAINER_ID",
		Short: "Return health status of the container",
		Run: func(c *cobra.Command, args []string) {

			if len(args) == 0 {
				c.Usage()
				return
			}

			req.SetPlugins(api.ParsePlugins(enabledPlugins, repo))

			req.SetArgs(args)

			if err := req.FindContainer(args[0]); err != nil {
				fmt.Println(err)
				return
			}

			ok := req.Execute()

			if !ok {
				if req.Verbose {
					log.Printf("Container %s health check FAILED", req.Container.ID)
				}
				os.Exit(1)
			} else {
				if req.Verbose {
					log.Printf("Container %s is OK", req.Container.ID)
				}
			}
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available plugins",
		Run: func(c *cobra.Command, args []string) {
			plugins := repo.List()
			for i := 0; i < len(plugins); i++ {
				fmt.Printf("%s\n", plugins[i].Name())
			}
		},
	}

	statusCmd.PersistentFlags().StringVarP((&enabledPlugins), "plugins", "P", "", "A comma separated list of plugins to use")
	statusCmd.PersistentFlags().StringVarP((&req.Port), "port", "p", "", "Network port to health check")

	mainCmd.AddCommand(statusCmd)
	mainCmd.AddCommand(listCmd)

	mainCmd.Execute()
}
