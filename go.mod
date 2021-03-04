module redis-cli

go 1.13

require (
	github.com/go-redis/redis/v8 v8.4.10
	github.com/peterh/liner v1.2.1
	github.com/spf13/cobra v1.1.1
)

replace github.com/peterh/liner => github.com/juwell/liner v0.0.0-20210304023022-6050fc0afd03
