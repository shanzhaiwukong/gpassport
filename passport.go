package gpassport

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

// Passport 结构体
type Passport struct {
	lock sync.Mutex
	cli  *redis.Client
	pre  string
}

// New 创建
func New(cli *redis.Client, pre string) *Passport {
	p := &Passport{
		cli: cli,
		pre: pre,
	}
	return p
}

// Add 添加
func (p *Passport) Add(token, userID string, expire time.Duration) *Passport {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.cli.Set(p.pre+token, userID, expire)
	p.cli.Set(p.pre+userID, token, expire)
	return p
}

// AddWithEntity 添加并带上实体
func (p *Passport) AddWithEntity(token, userID string, entity interface{}, expire time.Duration) *Passport {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.cli.Set(p.pre+token, userID, expire)
	p.cli.Set(p.pre+userID, token, expire)
	bs, _ := json.Marshal(entity)
	p.cli.Set(fmt.Sprintf("%s%s$%s", p.pre, token, userID), string(bs), expire)
	return p
}

// RemoveByToken 通过token移除
func (p *Passport) RemoveByToken(token string) *Passport {
	p.lock.Lock()
	defer p.lock.Unlock()
	if userID, err := p.cli.Get(p.pre + token).Result(); userID != "" && err == nil {
		p.cli.Del(p.pre + userID)
		p.cli.Del(fmt.Sprintf("%s%s$%s", p.pre, token, userID))
	}
	p.cli.Del(p.pre + token)
	return p
}

// RemoveByUserID 通过userID移除
func (p *Passport) RemoveByUserID(userID string) *Passport {
	p.lock.Lock()
	defer p.lock.Unlock()
	if token, err := p.cli.Get(p.pre + userID).Result(); token != "" && err == nil {
		p.cli.Del(p.pre + token)
		p.cli.Del(fmt.Sprintf("%s%s$%s", p.pre, token, userID))
	}
	p.cli.Del(p.pre + userID)
	return p
}

// UpdateByToken 通过Token更新有效时间
func (p *Passport) UpdateByToken(token string, expire time.Duration) *Passport {
	p.lock.Lock()
	defer p.lock.Unlock()
	if userID, err := p.cli.Get(p.pre + token).Result(); userID != "" && err == nil {
		p.cli.Expire(p.pre+userID, expire)
		p.cli.Expire(fmt.Sprintf("%s%s$%s", p.pre, token, userID), expire)
		p.cli.Expire(p.pre+token, expire)
	}
	return p
}

// UpdateByUserID 通过userID更新有效时间
func (p *Passport) UpdateByUserID(userID string, expire time.Duration) *Passport {
	p.lock.Lock()
	defer p.lock.Unlock()
	if token, err := p.cli.Get(p.pre + userID).Result(); token != "" && err == nil {
		p.cli.Expire(p.pre+token, expire)
		p.cli.Expire(fmt.Sprintf("%s%s$%s", p.pre, token, userID), expire)
		p.cli.Expire(p.pre+userID, expire)
	}
	return p
}

// GetUserID 获取用户ID
func (p *Passport) GetUserID(token string) (string, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.cli.Get(p.pre + token).Result()
}

// GetUserIDAndEntity 获取用户ID及实体
func (p *Passport) GetUserIDAndEntity(token string, entity interface{}) (userID string, err error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if userID, err = p.cli.Get(p.pre + token).Result(); userID == "" || err != nil {
		return
	}
	if v, _err := p.cli.Get(fmt.Sprintf("%s%s$%s", p.pre, token, userID)).Result(); v != "" && _err == nil {
		err = json.Unmarshal([]byte(v), entity)
	} else {
		err = _err
	}
	return
}

// Exists token是否存在
func (p *Passport) Exists(token string) bool {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.cli.Exists(p.pre+token).Val() == 1
}
