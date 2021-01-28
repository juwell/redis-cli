package cmd

import (
	"context"
	"fmt"
	"redis-cli/client"
	"reflect"

	"github.com/spf13/cobra"
)

var (
	// GetCmd get命令
	GetCmd = cobra.Command{
		Use:     `get`,
		Aliases: []string{`mget`},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) <= 0 {
				errArgsWithNumber(cmd.Name())
				return
			}

			ctx := context.Background()
			c := client.Cli.MGet(ctx, args...)
			t, err := c.Result()
			// str, err := redigo.String()
			if err != nil {
				fmt.Print(err.Error())
			} else {
				str := ``
				for _, v := range t {
					switch reflect.TypeOf(v).Kind() {
					case reflect.String:
						str += v.(string)
						str += `
`
					default:
						str += fmt.Sprintf(`err type:%v`, reflect.TypeOf(v))
					}
				}
				fmt.Print(str)
			}
		},
	}
)

func init() {
	// ClientCmd.AddCommand(&GetCmd)
}
