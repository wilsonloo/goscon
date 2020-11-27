package main

import (
	"fmt"
	"log"
	"sync"
)

import (
	"github.com/go-redis/redis"
)

const (
	RedisAddr = ":6379"
)

// InitRedis 初始化 Redis
func InitRedis() *redis.Client {
	option := &redis.Options{
		Addr: RedisAddr,
	}

	client := redis.NewClient(option)

	return client
}

func main()  {
	var wg sync.WaitGroup

	for k := 0; k < 2; k++{
		wg.Add(1)
		go func(i int){
			client := InitRedis()
			pb := client.Subscribe("channel-wait-close")
			fmt.Println("listening channel to notify...")

			for {
				select {
				case msg := <- pb.Channel():
					if msg.Payload == "close"{
						wg.Done()
					}else{
						log.Println("receive channel message:", msg.Payload)
					}
				default:
				}
			}
		}(k)
	}

	wg.Wait()
	log.Println("test finished")
}
