package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tkitsunai/messenger/server"
	"github.com/tkitsunai/messenger/utils"
	"os"
)

var (
	MessengerCmd = &cobra.Command{
		Use:   "messenger",
		Short: "'messenger' is simple subject-based messaging pub/sub server.",
		Long:  `'messenger' is simple subject-based messaging pub/sub server.`,
		RunE:  start,
	}
)

func init() {
	MessengerCmd.PersistentFlags().StringP("server", "s", "", "start server")
}

func start(cmd *cobra.Command, args []string) error {
	tcpServer := server.NewTcpServer()
	errChan := make(chan error)
	go tcpServer.StartHandler(fmt.Sprintf(":%s", utils.GetConfig().Port), errChan)

	err := <-errChan
	if err != nil {
		return fmt.Errorf("failed start server - %s", err.Error())
	}
	return nil
}

func Execute() {
	err := MessengerCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
