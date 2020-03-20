package remote

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

const LuaScript = `
        local ticket_key = KEYS[1]
        local ticket_total_key = ARGV[1]
        local ticket_sold_key = ARGV[2]
        local ticket_total_nums = tonumber(redis.call('HGET', ticket_key, ticket_total_key))
        local ticket_sold_nums = tonumber(redis.call('HGET', ticket_key, ticket_sold_key))
		-- 查看是否还有余票,增加订单数量,返回结果值
        if(ticket_total_nums >= ticket_sold_nums) then
            return redis.call('HINCRBY', ticket_key, ticket_sold_key, 1)
        end
        return 0
`

type RedisKeys struct {
	OrderHashKey     string
	TotalTicketField string
	SoldTicketField  string
}

func NewPool() *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", ":6379")
			if err != nil {
				fmt.Println("redis.Dial err=", err)
				panic("conn failed...")
			}
			return conn, err
		},
		MaxIdle:     8000,
		MaxActive:   10000,
		IdleTimeout: 6 * time.Second,
	}
}

//Redis端扣库存
func (keys *RedisKeys) RemoteGrabTicket(conn redis.Conn) bool {
	lua := redis.NewScript(1, LuaScript)
	res, err := redis.Int(lua.Do(conn, keys.OrderHashKey, keys.TotalTicketField, keys.SoldTicketField))
	if err != nil {
		fmt.Println("redis.Dial err=", err)
		panic("conn failed...")
		return false
	}
	return res != 0
}
