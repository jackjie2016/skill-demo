参考文档：https://zhuanlan.zhihu.com/p/558395137

redis中存放总库存，每个机器上面都有存放自己的本地库存，演示的redis中库存是本地的10倍

``
func (spike *LocalSpike) InitRemoteData(conn redis.Conn, remoteSpike remoteSpike.RemoteSpikeKeys) {
defer conn.Close()

	//local ticket_total_nums = tonumber(redis.call('HGET', ticket_key, ticket_total_key))
	//local ticket_sold_nums = tonumber(redis.call('HGET', ticket_key, ticket_sold_key))
	conn.Do("HSet", remoteSpike.SpikeOrderHashKey, remoteSpike.TotalInventoryKey, spike.LocalInStock*10)
	conn.Do("HSet", remoteSpike.SpikeOrderHashKey, remoteSpike.QuantityOfOrderKey, spike.LocalSalesVolume*10)
}
``

 


压测
``
ab -n 20000 -c 100 http://127.0.0.1:3005/buy/skill
``

![test.png](res%2Ftest.png)


redis链接池用完之后再链接报错：
``
net/http.(*conn).serve.func1()
E:/go_work/system/go1.19.4/src/net/http/server.go:1850 +0xbf
panic({0x475a00, 0xc00006ef90})
E:/go_work/system/go1.19.4/src/runtime/panic.go:890 +0x262
demo/remoteSpike.(*RemoteSpikeKeys).NewPool.func1()
E:/go_work/skill/remoteSpike/RemoteSpikeKeys.go:22 +0x76
github.com/gomodule/redigo/redis.(*Pool).dial(0xc00011407c?, {0x53c7f0?, 0xc00001c0b8?})
E:/go_work/project/pkg/mod/github.com/gomodule/redigo@v1.8.9/redis/pool.go:397 +0x38
github.com/gomodule/redigo/redis.(*Pool).GetContext(0xc000000c80, {0x53c7f0, 0xc00001c0b8})
E:/go_work/project/pkg/mod/github.com/gomodule/redigo@v1.8.9/redis/pool.go:254 +0x53f
github.com/gomodule/redigo/redis.(*Pool).Get(...)
E:/go_work/project/pkg/mod/github.com/gomodule/redigo@v1.8.9/redis/pool.go:186
main.handleReq({0x53c3c8, 0xc000138380}, 0x2cf3ba?)
E:/go_work/skill/main.go:46 +0x56
net/http.HandlerFunc.ServeHTTP(0xc000123af0?, {0x53c3c8?, 0xc000138380?}, 0x0?)
E:/go_work/system/go1.19.4/src/net/http/server.go:2109 +0x2f
net/http.(*ServeMux).ServeHTTP(0x0?, {0x53c3c8, 0xc000138380}, 0xc000112300)
E:/go_work/system/go1.19.4/src/net/http/server.go:2487 +0x149
net/http.serverHandler.ServeHTTP({0xc00010aa20?}, {0x53c3c8, 0xc000138380}, 0xc000112300)
E:/go_work/system/go1.19.4/src/net/http/server.go:2947 +0x30c
net/http.(*conn).serve(0xc000000f00, {0x53c860, 0xc0000a1320})
E:/go_work/system/go1.19.4/src/net/http/server.go:1991 +0x607
created by net/http.(*Server).Serve
E:/go_work/system/go1.19.4/src/net/http/server.go:3102 +0x4db
2023/01/30 16:10:01 http: panic serving 127.0.0.1:49226: dial tcp :6379: connectex: Only one usage of each socket address (protocol/network address/port) is normally permitted.

``


redigo pool 链接的处理完请求要shi释放掉，要不然会异常
``
defer conn.Close()
``