package cmd

import (
	"context"
	"fmt"
	"redis-cli/client"

	"github.com/spf13/cobra"
)

var (
	// InfoCmd info命令
	InfoCmd = cobra.Command{
		Use: `info`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			c := client.Cli.Do(ctx, `info`)
			str, _ := c.Text()
			fmt.Print(str)
		},
	}
)

func init() {
	// ClientCmd.AddCommand(&InfoCmd)
}
