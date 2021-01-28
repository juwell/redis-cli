package cmd

import "fmt"

func errArgsWithNumber(cmd string) {
	fmt.Printf(`(error) ERR wrong number of arguments for '%s' command
`, cmd)
}
