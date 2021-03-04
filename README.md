# redis-cli by golang

[Chinese](./ZH.md)

As you know, if you want use `redis-cli` command in Windows, you can use `wsl`, or `docker`(install docker, and install redis container, than use `docker exec -it redis redis-cli` command), even can install `VMware` and install a `Linux` system.
All of these ways are not convenient than using `redis-cli` directly in Windows.

There is a respository by Microsoft, but it is discarded.

And the official `redis-cli` is writed by `c`, it is just used in Linux.

So I dicided rewrite a new `redis-cli` by `golang`, because `golang` is cross-platform's language.

This `redis-cli` is same as official `redis-cli` in the operate and the result.

You can download the executable file in `Releases`, and put it into `$PATH`, that't it.