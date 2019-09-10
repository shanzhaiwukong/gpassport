package gpassport_test

import (
	"fmt"
	"gpassport"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

type Demo struct {
	Name string    `json:"name"`
	Date time.Time `json:"time"`
}

func Test(t *testing.T) {
	cli := redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     "192.168.80.130:6379",
		DB:       0,
		Password: "123456",
	})
	p := gpassport.New(cli, "")
	p.AddWithEntity("a", "uid123456", &Demo{"张三", time.Now()}, time.Second*120)
	fmt.Println(p.GetUserID("a3"))
	_d := new(Demo)
	fmt.Println(p.GetUserIDAndEntity("a", _d))
	fmt.Println(_d)
	fmt.Println("存在a3==>", p.Exists("a3"))
	fmt.Println("存在a==>", p.Exists("a"))
}
