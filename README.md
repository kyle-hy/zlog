# zap 开源日志库的异步文件日志Sink
* 自定义输出方式的Sink
  * 打印的日志先放入channel缓存
  * 后台起一个协程读取channel写入文件
  * 使用bufio，channel缓存有多条日志则合并，len(chan)为0则对bufio直接Flush。
  * 使用lumberjack滚动日志

* 后台写文件的协程，可使用runtime.SetFinalizer优化，更优雅。
 
