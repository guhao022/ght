package cmd

type HelpText struct {
	Help int
	Run  int
}

var Helptext = `
  ght - 为了更好的解决api测试问题（测试真蛋疼...）

    使用:

        ght command [argument] [argument] ......

    arguments 包括:

        run [command] [uri]	运行代理服务器

            run - 运行反向代理服务器
            Usage:
                run [ip/server]    设置需要代理的服务器，输出log信息，方便调试

    使用 "[command] help" 获取更多关于此命令的信息.
`

var Runhelp = `
    ght run - 运行反向代理服务器
        Usage:
            run -dev [uri]	开发模式，输出log信息
            run -stable	[uri]	稳定模式，不输出log信息
`
