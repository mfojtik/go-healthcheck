package main

import (
	"fmt"
	"log"
	"os"

	"github.com/openshift/go-healthcheck/api"
	"github.com/openshift/go-healthcheck/plugins"
	"github.com/spf13/cobra"
)

var repo = api.Repository{}

func registerPlugins() {
	repo.Register(plugins.HttpPlugin{})
	repo.Register(plugins.MongoPlugin{})
	repo.Register(plugins.TcpPlugin{})
}

func main() {

	req := &api.StatusRequest{}
	registerPlugins()

	mainCmd := &cobra.Command{
		Use:   "healthchk",
		Short: "A simple tool for health checking Docker containers",
		Run: func(c *cobra.Command, args []string) {
			c.Usage()
		},
	}

	mainCmd.PersistentFlags().BoolVar((&req.Verbose), "verbose", false, "Enable verbose output")
	mainCmd.PersistentFlags().StringVarP((&req.Port), "port", "p", "", "Network port to health check")
	mainCmd.PersistentFlags().StringVarP((&req.Socket), "socket", "", "unix:///var/run/docker.sock", "Docker socket to use")

	statusCmd := &cobra.Command{
		Use:   "status CONTAINER_ID",
		Short: "Return health status of the container",
		Run: func(c *cobra.Command, args []string) {

			if len(args) == 0 {
				c.Usage()
				return
			}

			req.SetArgs(args)

			if err := req.FindContainer(args[0]); err != nil {
				fmt.Println(err)
				return
			}

			req.InitializePlugins(&repo)

			ok := req.Execute()

			if !ok {
				if req.Verbose {
					log.Printf("Container %s health check failed.", req.Container.ID)
				}
				os.Exit(1)
			} else {
				if req.Verbose {
					log.Printf("Container %s health check succeeded.", req.Container.ID)
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

	statusCmd.PersistentFlags().StringVarP((&req.PluginList), "plugins", "P", "", "A comma separated list of plugins to use")

	mainCmd.AddCommand(statusCmd)
	mainCmd.AddCommand(listCmd)

	mainCmd.Execute()
}
