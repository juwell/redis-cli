package main

import (
	"context"
	"fmt"
	"os"
	"redis-cli/client"
	"redis-cli/cmd"
	"reflect"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/peterh/liner"
	"github.com/spf13/cobra"
)

func main() {
	fmt.Println(os.Getpid())

	// inputReader := bufio.NewReader(os.Stdin)
	// inputReader.ReadString('\n')

	root := &cobra.Command{
		Run:     doOnce,
		Version: client.Version(),
		// DisableFlagsInUseLine: false,
	}
	root.Flags().Bool(`help`, false, `help for this command`)
	root.Flags().BoolP(`version`, `v`, false, `Output version and exit`)
	root.SetVersionTemplate(`redis-cli {{printf "%s" .Version}}
`)
	root.SetUsageTemplate(usageTemplate())
	root.Flags().StringVarP(&client.Cfg.HostIP, `hostname`, `h`, `127.0.0.1`, `Server hostname`)
	root.Flags().StringVarP(&client.Cfg.HostPort, `port`, `p`, `6379`, `Server port`)
	root.Flags().StringVarP(&client.Cfg.HostSocket, `socket`, `s`, ``, `Server socket (overrides hostname and port)`)
	root.Flags().IntVarP(&client.Cfg.DBNum, `db`, `n`, 0, `Database number`)
	root.Flags().StringVar(&client.Cfg.UserName, `user`, ``, `Used to send ACL style 'AUTH username pass'. Needs -a`)
	root.Flags().StringVarP(&client.Cfg.PassWord, `pass`, `a`, ``, `Password to use when connecting to the server.
You can also use the REDISCLI_AUTH environment
variable to pass this password more safely`)
	root.Flags().BoolVarP(&client.Cfg.ClusterMode, `cluster`, `c`, false, `Cluster Manager command and arguments (see below).`)

	if err := root.Execute(); err != nil {
		fmt.Print(err.Error())
		return
	}

	// fmt.Println(`root cmd exist`)

	if len(client.Cfg.HostSocket) <= 0 {
		client.Cfg.HostSocket = client.Cfg.HostIP + ":" + client.Cfg.HostPort
	}

	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == `--help` || os.Args[i] == `--version` ||
			os.Args[i] == `-v` || os.Args[i] == `?` {
			return
		}
	}

	// 连接后开始等待输入命令, 执行
	// doing()
}

// 只执行一次
func doOnce(c *cobra.Command, args []string) {
	if len(args) > 0 {
		if args[0] == `help` || args[0] == `?` {
			simpleHelp()
		} else {
			// 启动连接
			if err := client.Cli.Connect(); err != nil {
				fmt.Printf(`Could not connect to Redis at %s: Connection refused`, client.Cfg.HostSocket)
			} else {
				working(args)

				client.Cli.Close()
			}

		}
		os.Exit(0)
	} else {
		doing()
	}
}

// 连续执行
func doing() {
	if err := client.Cli.Connect(); err != nil {
		fmt.Printf(`Could not connect to Redis at %s: Connection refused`, client.Cfg.HostSocket)
		os.Exit(1)
		return
	}
	// commandList := []string{
	// 	`set key value`, `hget key field`,
	// 	`hgetall key`,
	// }
	// inputReader := bufio.NewReader(os.Stdin)
	line := liner.NewLiner()
	line.SetCtrlCAborts(true)
	line.SetHitsCallback(hitsCallback)
	line.SetWordCompleter(completionCallback)
	prompt := ``

	for {
		if len(client.Cfg.HostSocket) > 0 {
			prompt = fmt.Sprintf(`%s> `, client.Cfg.HostSocket)
		} else {
			prompt = fmt.Sprintf(`%s:%s> `, client.Cfg.HostIP, client.Cfg.HostPort)
		}
		// strs, err := inputReader.ReadString('\n')
		str, err := line.Prompt(prompt)
		// strs = strs[:len(strs)-2]
		if err != nil {
			if err == liner.ErrPromptAborted {
				os.Exit(0)
			} else {
				fmt.Println(err.Error())
			}
		} else if len(str) <= 0 {
			// 没有输入, 直接换行
			continue
		} else {
			line.AppendHistory(str)
			cmds := strings.Split(str, ` `)
			if len(cmds) <= 0 {
				continue
			}
			// fmt.Println(`cmds:`, cmds)

			working(cmds)
		}
	}
}

// 执行一次命令
func working(cmds []string) {
	if len(cmds) <= 0 {
		// 没有命令, 啥都不做
		return
	}

	// 如果不需要发消息, 则会直接退出
	if !analysisCmd(cmds) {
		return
	}

	t := make([]interface{}, 0, len(cmds))
	for _, s := range cmds {
		// fmt.Println(`s:`, s)
		if len(s) <= 0 {
			continue
		}
		t = append(t, s)
	}
	str := ``

	// fmt.Println(`here`)
	switch t[0] {
	case "monitor":
		cmd.MonitorCmd.SetArgs(cmds)
		cmd.MonitorCmd.Execute()
	default:
		// 直接发到服务器, 然后打印返回信息
	DOCOMMAND:
		ctx := context.Background()
		re := client.Cli.Do(ctx, t...)
		val, err := re.Result()
		switch {
		case err == redis.Nil:
			fmt.Println(`(nil)`)

		case err != nil:
			errStr := err.Error()
			if client.Cfg.ClusterMode > 0 &&
				(len(errStr) > 5 && errStr[:5] == `MOVED`) ||
				(len(errStr) > 3 && errStr[:3] == `ASK`) {
				// 集群模式的重定向
				slot := 0
				host := ``
				info := ``
				fmt.Sscanf(errStr, `%s %d %s`, &info, &slot, &host)
				temp := strings.Split(host, `:`)
				if len(temp) == 2 {
					client.Cfg.HostIP = temp[0]
					client.Cfg.HostPort = temp[1]
				}
				client.Cfg.HostSocket = host

				if err := client.Cli.Redirection(); err != nil {
					fmt.Println(err)
					os.Exit(1)
					return
				}

				// 加上重定向提示
				str = fmt.Sprintf(`-> Redirected to slot [%d] located at %s
	`, slot, host)
				goto DOCOMMAND
			} else {
				fmt.Println(`(error) ` + err.Error())
			}

		case err == nil:

			// fmt.Println(`here4`, fmt.Sprintf(`type:%v, kind:%v`, reflect.TypeOf(val), reflect.TypeOf(val).Kind()))
			switch reflect.TypeOf(val).Kind() {
			case reflect.String:
				str += `"` + val.(string) + `"`
			case reflect.Int64:
				str += fmt.Sprintf(`(integer) %v`, val.(int64))
			case reflect.Slice:
				temp := val.([]interface{})
				if len(temp) <= 0 {
					str += `(empty array)`
				}
				for i, s := range temp {
					if i > 0 {
						str += `
	`
					}
					switch reflect.TypeOf(s).Kind() {
					case reflect.String:
						str += fmt.Sprintf(`%d) "%s"`, i+1, s.(string))
					case reflect.Int64:
						str += fmt.Sprintf(`%d) (integer) %v`, i+1, s.(int64))
					default:
						str += fmt.Sprintf(`%d)err type:%v, kind:%v`, i+1, reflect.TypeOf(s), reflect.TypeOf(s).Kind())
					}
				}
			default:
				str += fmt.Sprintf(`err type:%v, kind:%v`, reflect.TypeOf(val), reflect.TypeOf(val).Kind())
			}
			fmt.Println(str)
		}
	}
}

func analysisCmd(cmds []string) bool {
	if len(cmds) <= 0 {
		return false
	}

	switch cmds[0] {
	case `help`:
		fallthrough
	case `?`:
		if len(cmds) == 1 {
			simpleHelp()
		} else if len(cmds) >= 2 {
			// 查找命令帮助
			c := strings.Join(cmds[1:], ` `)
			h, ok := cmd.CommandHelps.Find(c)
			if !ok {
				fmt.Println()
			} else {
				fmt.Println(commandHelpTemplate(h))
			}
		} else {
			fmt.Println()
		}
		return false
	case `version`:
		fmt.Printf(versionTemplate())
		return false
	}
	return true
}

func hitsCallback(line string) (string, int, bool) {
	if len(line) <= 0 {
		return ``, liner.ColorCodeGray, false
	}

	lines := strings.Split(line, ` `)
	if len(lines) <= 0 {
		lines = []string{line}
	} else {
		for i := 0; i < len(lines); i++ {
			if len(lines[i]) <= 0 {
				lines = append(lines[:i], lines[i+1:]...)
				i--
			}
		}
	}
	parms := make([]string, 0, 10)
	control := false
	conStr := ``
	index := 0
	temp := ``
	for i := len(lines); i >= 0; i-- {
		// fmt.Println(`[test]`, lines[:i])
		if h, ok := cmd.CommandHelps.Find(strings.Join(lines[:i], ` `)); ok {
			// 整理参数
			index = 0
			conStr = ``
			control = false
			temp = ``
			for t := 0; t < len(h.Params); t++ {
				switch h.Params[t] {
				case '[':
					control = true
					conStr += string(h.Params[t])
				case ']':
					control = false
					conStr += string(h.Params[t])
					index = t + 1
				case ' ':
					if control {
						conStr += string(h.Params[t])
					} else if index == t {
						index++
					} else {
						parms = append(parms, temp)
						index = t + 1
						temp = ``
					}
				default:
					if control {
						conStr += string(h.Params[t])
					} else {
						temp += string(h.Params[t])
					}
				}
			}

			if len(temp) > 0 {
				parms = append(parms, temp)
			}

			outStr := conStr
			if len(lines)-1 < len(parms) {
				outStr = strings.Join(parms[len(lines)-1:], ` `) + ` ` + conStr
			}

			if line[len(line)-1] != ' ' {
				return ` ` + outStr, liner.ColorCodeGray, false
			}
			return outStr, liner.ColorCodeGray, false
		}
	}

	// if h, ok := cmd.CommandHelps.Find(line); ok {
	// 	if line[len(line)-1] != ' ' {
	// 		return ` ` + h.Params, liner.ColorCodeGray, false
	// 	}
	// 	return h.Params, liner.ColorCodeGray, false
	// }
	return ``, liner.ColorCodeGray, false
}

func completionCallback(line string, pos int) (head string, completions []string, tail string) {
	if strings.HasPrefix(strings.ToUpper(line), `HELP `) {
		head = line[:5]
		line = line[5:]
		pos -= 5
	}

	tail = line[pos:]
	i := strings.LastIndex(line[:pos], ` `)
	if i > 0 && i < len(line) {
		head += line[:i]
	}
	if pos < len(line) {
		tail = line[pos+1:]
	}
	tempStr := line[i+1 : pos]

	completions = make([]string, 0, 10)
	for key := range cmd.CommandHelps {
		if strings.HasPrefix(key, strings.ToUpper(tempStr)) {
			completions = append(completions, key)
		}
	}

	return
}
