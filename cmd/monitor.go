package cmd

import (
	"fmt"
	"redis-cli/client"

	"github.com/spf13/cobra"
)

var (
	// MonitorCmd monitor命令
	MonitorCmd = cobra.Command{
		Use: `monitor`,
		Run: func(cmd *cobra.Command, args []string) {
			c := client.NewMonitorClient(*client.Cli.Options())
			if err := c.Connect(); err != nil {
				fmt.Println(err.Error())
			} else {
				c.Do(func(d []byte) {
					if len(d) <= 0 {
						fmt.Println(`[E] len(d) <= 0`)
						return
					}

					// switch d[0] {
					// case '-':
					// case '+':
					// case ':':
					// case ',':
					// case '_':
					// case '$':
					// case '*':
					// case '%':
					// case '~':
					// case '#':
					// case '=':
					// case '>':
					// }

					fmt.Println(string(d[1:]))
				})
			}
			c.Close()
		},
	}
)

// func monitorProcess(ctx context.Context, cmd redis.Cmder) error {
// 	conn := client.Cli.Conn(ctx)

// }
