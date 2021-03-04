@echo off

@rmdir out /q /s
@mkdir out
@set version=0.0.1

@echo ==Build Windows==
@go env -w GOOS=windows
@go build -o ./out/redis-cli.exe
@7z a ./out/redis-cli-v%version%-windows-amd64.7z ./out/redis-cli.exe -sdel

@echo ==Build Linux==
@go env -w GOOS=linux
@go build -o ./out/redis-cli
@7z a ./out/redis-cli-v%version%-linux-amd64.7z ./out/redis-cli -sdel

@echo ==Build MacOS==
@go env -w GOOS=darwin
@go build -o ./out/redis-cli
@7z a ./out/redis-cli-v%version%-darwin-amd64.7z ./out/redis-cli -sdel

@echo ==Finish==
@pause