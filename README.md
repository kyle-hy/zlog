# zap 日志库的异步文件日志Sink
## Sink
* 自定义输出协议的Sink
  * 打印的日志先放入channel缓存
  * 后台起一个协程读取channel写入文件
  * 使用bufio，channel缓存有多条日志则合并，len(chan)为0则对bufio直接Flush。
  * 使用lumberjack滚动日志

* TODO:后台写文件的协程，可使用runtime.SetFinalizer优化，更优雅。

## 使用方式
需要执行InitLog函数初始化
``` go
zaplog.InitLog(zaplog.BufioSize(1024*8), zaplog.WithFields(map[string]interface{}{"app": "dddd"}))
```

## 测试结果

MacBook Pro (13-inch, M1, 2020) 
起1000个协程，总共写1000w条日志
* 有bufio：1.099µs/p
* 无bufio：3.089µs/p
