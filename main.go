package main

import (
	"Ticket12306/helper"
	"Ticket12306/local"
	"Ticket12306/remote"
	"github.com/gomodule/redigo/redis"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	localop   local.LocalGrab
	remoteop  remote.RedisKeys
	redisPool *redis.Pool
	mu        chan struct{} //用来同步，使用通信来共享内存而不是使用共享内存来通信

)

func init() {
	localop = local.LocalGrab{
		LocalTotal: 5000,
		LocalSold:  0,
	}
	remoteop = remote.RedisKeys{
		OrderHashKey:     "ticket_key",
		TotalTicketField: "ticket_total",
		SoldTicketField:  "ticket_sold",
	}
	redisPool = remote.NewPool()
	mu = make(chan struct{},1)
	mu <- struct{}{}
}
func main() {
	http.HandleFunc("/grabTicket", func(w http.ResponseWriter, r *http.Request) {
		conn := redisPool.Get()
		LogMsg := ""
		<-mu
		if localop.LocalGrabTicket() && remoteop.RemoteGrabTicket(conn) {
			helper.Resp(w, 1, "succeed", nil)
			LogMsg += "result:1,localSales:" + strconv.FormatInt(localop.LocalSold, 10)
		} else {

			helper.Resp(w, -1, "failed,ticket has been solded out", nil)
			LogMsg += "result:0,localSales:" + strconv.FormatInt(localop.LocalSold, 10)
		}
		mu <- struct{}{}
		writeLog(LogMsg, "./stat.log")
	})
	http.ListenAndServe(":5555", nil)
}
func writeLog(msg string, logPath string) {
	fd, _ := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer fd.Close()
	content := strings.Join([]string{msg, "\r\n"}, "")
	buf := []byte(content)
	fd.Write(buf)
}
