package main

import (
	localSpike2 "demo/localSpike"
	remoteSpike2 "demo/remoteSpike"
	"demo/util"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var redisPool *redis.Pool
var localSpike localSpike2.LocalSpike
var remoteSpike remoteSpike2.RemoteSpikeKeys
var done chan int

func init() {

	remoteSpike = remoteSpike2.RemoteSpikeKeys{
		SpikeOrderHashKey:  "ticket_hash_key",
		TotalInventoryKey:  "ticket_total_nums",
		QuantityOfOrderKey: "ticket_sold_nums",
	}
	redisPool = remoteSpike.NewPool()

	localSpike = localSpike2.NewLocalSpike()

	localSpike.InitRemoteData(redisPool.Get(), remoteSpike)

	done = make(chan int, 1)
	done <- 1
}

var err, total, through int

func main() {

	http.HandleFunc("/buy/skill", handleReq)
	http.ListenAndServe(":3005", nil)
}

func handleReq(w http.ResponseWriter, r *http.Request) {
	redisConn := redisPool.Get()

	LogMsg := ""
	<-done //全局读写锁
	total++
	if localSpike.LocalDeductionStock() && remoteSpike.RemoteDeductionStock(redisConn) {
		util.RespOk2(w, 1, "抢票成功")
		through++
		LogMsg = LogMsg + "result:1,localSales:" + strconv.FormatInt(localSpike.LocalSalesVolume, 10)
	} else {
		err++
		util.RespOk2(w, -1, "已售罄")
		LogMsg = LogMsg + "result:0,localSales:" + strconv.FormatInt(localSpike.LocalSalesVolume, 10)
	}
	done <- 1
	fmt.Println("err:", err, " ，total:", total, " ，through:", through, "localSpike:", localSpike)
	//将抢票状态写入到log中
	writeLog(LogMsg, "./stat.log")
}

func writeLog(msg string, logPath string) {
	fd, _ := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer fd.Close()
	content := strings.Join([]string{msg, "\r\n"}, "")
	buf := []byte(content)
	fd.Write(buf)
}
