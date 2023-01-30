package remoteSpike

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type RemoteSpikeKeys struct {
	SpikeOrderHashKey  string //redis中秒杀订单hash结构key
	TotalInventoryKey  string //hash结构中总订单库存key
	QuantityOfOrderKey string //hash结构中已有订单数量key

}

func (r *RemoteSpikeKeys) NewPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   10000,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

// SpikeOrderHashKey:  "ticket_hash_key",
// TotalInventoryKey:  "ticket_total_nums",
// QuantityOfOrderKey: "ticket_sold_nums",
const LuaScript = `
local ticket_key = KEYS[1] 
local ticket_total_key = ARGV[1]  
local ticket_sold_key = ARGV[2]
local ticket_total_nums = tonumber(redis.call('HGET', ticket_key, ticket_total_key))      
local ticket_sold_nums = tonumber(redis.call('HGET', ticket_key, ticket_sold_key))   
       if(ticket_total_nums >= ticket_sold_nums) then
            return redis.call('HINCRBY', ticket_key, ticket_sold_key, 1) end
        return 0`

const lua = `redis.call('SET', KEYS[1], ARGV[1]);redis.call('EXPIRE', KEYS[1], ARGV[2]); return 1`
const lua2 = `local key = KEYS[1];
local increment = ARGV[1];
local stock = ARGV[2];
if 1 == redis.call('exists',key) then 
    if tonumber(stock) < (tonumber(increment) + tonumber(redis.call('get',key))) then 
        return 0
    else
        return redis.call('incrby',key,increment)
    end
else
    return redis.call('incrby',key,increment)
end
`

// 远端统一扣库存
func (RemoteSpikeKeys *RemoteSpikeKeys) RemoteDeductionStock(conn redis.Conn) bool {
	lua := redis.NewScript(1, LuaScript)

	result, err := redis.Int(lua.Do(conn, RemoteSpikeKeys.SpikeOrderHashKey, RemoteSpikeKeys.TotalInventoryKey, RemoteSpikeKeys.QuantityOfOrderKey))
	if err != nil {
		fmt.Println("redis 异常", err)
		return false
	}

	defer conn.Close()
	return result != 0
}
