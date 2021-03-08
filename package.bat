@echo off

@rmdir out /q /s
@mkdir out

@REM !!!! Modify Those Variable First !!!!
@set version=0.1.1
@set redisVersion=6.2.1
@REM !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

@echo ==Write Version==
@set versionFile=version.go
@echo package main > %versionFile%%
@echo const ( >> %versionFile%
@echo     Version      = `%version%` >> %versionFile%
@echo     RedisVersion = `%redisVersion%` >> %versionFile%
@echo ) >> %versionFile%

@REM Save Current "GOOS", Will Be Set Back After "go build"
@for /F %%i in ('go env GOOS') DO @set oldEnv=%%i

@echo ==Build Windows==
@go env -w GOOS=windows
@go build -o ./out/redis-cli.exe
@go env -w GOOS=%oldEnv%
@7z a ./out/redis-cli-v%version%-windows-amd64.7z ./out/redis-cli.exe -sdel

@echo ==Build Linux==
@go env -w GOOS=linux
@go build -o ./out/redis-cli
@go env -w GOOS=%oldEnv%
@7z a ./out/redis-cli-v%version%-linux-amd64.7z ./out/redis-cli -sdel

@echo ==Build MacOS==
@go env -w GOOS=darwin
@go build -o ./out/redis-cli
@go env -w GOOS=%oldEnv%
@7z a ./out/redis-cli-v%version%-darwin-amd64.7z ./out/redis-cli -sdel

@echo ==Finish==

@pause