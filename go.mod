module redis-cli

go 1.13

require (
	github.com/go-redis/redis/v8 v8.4.10
	// github.com/gomodule/redigo v1.8.3 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/go-redis/redis/extra/redisotel v0.0.0
)

replace github.com/go-redis/redis/v8 => ../go-redis

replace github.com/go-redis/redis/extra/redisotel => ../go-redis/extra/redisotel
