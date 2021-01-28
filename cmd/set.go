package cmd

import (
	"context"
	"fmt"
	"redis-cli/client"

	"github.com/spf13/cobra"
)

var (
	// SetCmd set命令
	SetCmd = cobra.Command{
		Use:     `set`,
		Aliases: []string{`mset`},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) <= 0 {
				errArgsWithNumber(cmd.Name())
				return
			}

			ctx := context.Background()
			t := make([]interface{}, len(args))
			for i, s := range args {
				t[i] = s
			}
			c := client.Cli.MSetNX(ctx, t...)
			_, err := c.Result()
			// str, err := redigo.String(c.Result())
			if err != nil {
				fmt.Print(err.Error())
			} else {
				fmt.Print(`OK`)
			}
		},
	}
)

func init() {
	// ClientCmd.AddCommand(&SetCmd)
}
