# redis-cli by golang

众所周知, Win下要使用`redis-cli`命令, 要么`wsl`, 要么`docker`, 甚至可以在`VMware`下安装一个linux系统, 再安装`redis-cli`, 这些都没有直接在命令行下执行`redis-cli.exe`来的方便.

而早期, 微软有维护一个项目, 但停留在3.0阶段, 已经废弃了.

而redis官方的则是用c写的, 而且只支持linux系统, 要在win上编译, 则需要改大量代码, 估计也是这个原因, 让微软放弃了.

于是乎, 利用`golang`的跨平台性, 重写了一个`redis-cli`, 操作方式, 以及结果显示也与官方的一致.

会在`Releases`中直接放出编译好的可执行程序, 直接下载放到环境变量`$PATH`下就可以使用了.

