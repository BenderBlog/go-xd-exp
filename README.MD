# go-xd-exp

仅仅是一个中转学校物理实验信息的服务器，本意上是绕过 Dart 对 Header 检测的一个临时方式。

使用了 go-fiber 库作为服务器框架，go-resty 作为 HTTP 请求工具，另外有 GBK 转 UTF-8 的需求。

这个程序将会在某个学校内部的服务器上运行，且仅仅作为一种临时措施。等到后面版本，我会尝试使用 Dart 的原生方式实现本地访问物理实验服务器。