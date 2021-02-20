package main

import (
	"fmt"
	"redis-cli/client"
	"redis-cli/cmd"
)

func simpleHelp() {
	fmt.Printf(`redis-cli %s with redis %s
To get help about Redis commands type:
    "help @<group>" to get a list of commands in <group>
    "help <command>" for help on <command>
    "help <tab>" to get a list of possible help topics
    "quit" to exit

To set redis-cli preferences:
    ":set hints" enable online hints
    ":set nohints" disable online hints
Set your preferences in ~/.redisclirc`, client.Version(), client.RedisVersion())
}

func usage() {
	fmt.Printf(`redis-cli %s with redis %s

Usage: redis-cli [OPTIONS] [cmd [arg [arg ...]]]
  -h <hostname>      Server hostname (default: 127.0.0.1).
  -p <port>          Server port (default: 6379).
  -s <socket>        Server socket (overrides hostname and port).
  -a <password>      Password to use when connecting to the server.
					 You can also use the REDISCLI_AUTH environment
					 variable to pass this password more safely
					 (if both are used, this argument takes precedence).
  --user <username>  Used to send ACL style 'AUTH username pass'. Needs -a.
  --pass <password>  Alias of -a for consistency with the new --user option.
  --askpass          Force user to input password with mask from STDIN.
					 If this argument is used, '-a' and REDISCLI_AUTH
					 environment variable will be ignored.
  -u <uri>           Server URI.
  -r <repeat>        Execute specified command N times.
  -i <interval>      When -r is used, waits <interval> seconds per command.
					 It is possible to specify sub-second times like -i 0.1.
  -n <db>            Database number.
  -3                 Start session in RESP3 protocol mode.
  -x                 Read last argument from STDIN.
  -d <delimiter>     Delimiter between response bulks for raw formatting (default: \n).
  -D <delimiter>     Delimiter between responses for raw formatting (default: \n).
  -c                 Enable cluster mode (follow -ASK and -MOVED redirections).
  --tls              Establish a secure TLS connection.
  --sni <host>       Server name indication for TLS.
  --cacert <file>    CA Certificate file to verify with.
  --cacertdir <dir>  Directory where trusted CA certificates are stored.
					 If neither cacert nor cacertdir are specified, the default
					 system-wide trusted root certs configuration will apply.
  --cert <file>      Client certificate to authenticate with.
  --key <file>       Private key file to authenticate with.
  --raw              Use raw formatting for replies (default when STDOUT is
					 not a tty).
  --no-raw           Force formatted output even when STDOUT is not a tty.
  --csv              Output in CSV format.
  --stat             Print rolling stats about server: mem, clients, ...
  --latency          Enter a special mode continuously sampling latency.
					 If you use this mode in an interactive session it runs
					 forever displaying real-time stats. Otherwise if --raw or
					 --csv is specified, or if you redirect the output to a non
					 TTY, it samples the latency for 1 second (you can use
					 -i to change the interval), then produces a single output
					 and exits.
  --latency-history  Like --latency but tracking latency changes over time.
					 Default time interval is 15 sec. Change it using -i.
  --latency-dist     Shows latency as a spectrum, requires xterm 256 colors.
					 Default time interval is 1 sec. Change it using -i.
  --lru-test <keys>  Simulate a cache workload with an 80-20 distribution.
  --replica          Simulate a replica showing commands received from the master.
  --rdb <filename>   Transfer an RDB dump from remote server to local file.
  --pipe             Transfer raw Redis protocol from stdin to server.
  --pipe-timeout <n> In --pipe mode, abort with error if after sending all data.
					 no reply is received within <n> seconds.
					 Default timeout: 30. Use 0 to wait forever.
  --bigkeys          Sample Redis keys looking for keys with many elements (complexity).
  --memkeys          Sample Redis keys looking for keys consuming a lot of memory.
  --memkeys-samples <n> Sample Redis keys looking for keys consuming a lot of memory.
					 And define number of key elements to sample
  --hotkeys          Sample Redis keys looking for hot keys.
					 only works when maxmemory-policy is *lfu.
  --scan             List all keys using the SCAN command.
  --pattern <pat>    Keys pattern when using the --scan, --bigkeys or --hotkeys
					 options (default: *).
  --intrinsic-latency <sec> Run a test to measure intrinsic system latency.
					 The test will run for the specified amount of seconds.
  --eval <file>      Send an EVAL command using the Lua script at <file>.
  --ldb              Used with --eval enable the Redis Lua debugger.
  --ldb-sync-mode    Like --ldb but uses the synchronous Lua debugger, in
					 this mode the server is blocked and script changes are
					 not rolled back from the server memory.
  --cluster <command> [args...] [opts...]
					 Cluster Manager command and arguments (see below).
  --verbose          Verbose mode.
  --no-auth-warning  Don't show warning message when using password on command
					 line interface.
  --help             Output this help and exit.
  --version          Output version and exit.
Cluster Manager Commands:
  Use --cluster help to list all available cluster manager commands.
Examples:
  cat /etc/passwd | redis-cli -x set mypasswd
  redis-cli get mypasswd
  redis-cli -r 100 lpush mylist x
  redis-cli -r 100 -i 1 info | grep used_memory_human:
  redis-cli --eval myscript.lua key1 key2 , arg1 arg2 arg3
  redis-cli --scan --pattern '*:12345*'
  (Note: when using --eval the comma separates KEYS[] from ARGV[] items)
When no command is given, redis-cli starts in interactive mode.
Type "help" in interactive mode for information on available commands
and settings.`, client.Version(), client.RedisVersion())
}

func usageTemplate() string {
	return `Usage: redis-cli [OPTIONS] [cmd [arg [arg ...]]]
 
Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}

Cluster Manager Commands:
  Use --cluster help to list all available cluster manager commands.

Examples:
  cat /etc/passwd | redis-cli -x set mypasswd
  redis-cli get mypasswd
  redis-cli -r 100 lpush mylist x
  redis-cli -r 100 -i 1 info | grep used_memory_human:
  redis-cli --eval myscript.lua key1 key2 , arg1 arg2 arg3
  redis-cli --scan --pattern '*:12345*'
  (Note: when using --eval the comma separates KEYS[] from ARGV[] items)

When no command is given, redis-cli starts in interactive mode.
Type "help" in interactive mode for information on available commands
and settings.

`
}

func commandHelpTemplate(h cmd.CommandHelp) string {
	out := fmt.Sprintf("\r\n  \x1b[1m%s\x1b[0m \x1b[90m%s\x1b[0m\r\n", h.Name, h.Params)
	out += fmt.Sprintf("  \x1b[33msummary:\x1b[0m %s\r\n", h.Summary)
	out += fmt.Sprintf("  \x1b[33msince:\x1b[0m %s\r\n", h.Since)
	if len(h.Group) > 0 {
		out += fmt.Sprintf("  \x1b[33mgroup:\x1b[0m %s\r\n", h.Group)
	}
	return out
}
