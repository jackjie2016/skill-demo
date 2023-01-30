package localSpike

import (
	"demo/remoteSpike"
	"github.com/gomodule/redigo/redis"
)

type LocalSpike struct {
	LocalInStock     int64
	LocalSalesVolume int64
}

func NewLocalSpike() LocalSpike {
	return LocalSpike{
		LocalInStock:     200,
		LocalSalesVolume: 0,
	}
}

// InitRemoteData 初始化远端数据，是本地数据的10倍
func (spike *LocalSpike) InitRemoteData(conn redis.Conn, remoteSpike remoteSpike.RemoteSpikeKeys) {
	defer conn.Close()

	//local ticket_total_nums = tonumber(redis.call('HGET', ticket_key, ticket_total_key))
	//local ticket_sold_nums = tonumber(redis.call('HGET', ticket_key, ticket_sold_key))
	conn.Do("HSet", remoteSpike.SpikeOrderHashKey, remoteSpike.TotalInventoryKey, spike.LocalInStock*10)
	conn.Do("HSet", remoteSpike.SpikeOrderHashKey, remoteSpike.QuantityOfOrderKey, spike.LocalSalesVolume*10)
}
func (spike *LocalSpike) LocalDeductionStock() bool {
	spike.LocalSalesVolume = spike.LocalSalesVolume + 1
	return spike.LocalSalesVolume <= spike.LocalInStock
}
