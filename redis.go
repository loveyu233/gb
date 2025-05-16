package gb

import (
	"context"
	"errors"
	r "github.com/go-redis/redis/v8"
	"strings"
	"sync"
	"time"
)

type OptionsFunc func(*r.UniversalOptions)

func WithDB(db int) OptionsFunc {
	return func(o *r.UniversalOptions) {
		o.DB = db
	}
}

func WithUsername(username string) OptionsFunc {
	return func(o *r.UniversalOptions) {
		o.Username = username
	}
}

func WithPassword(password string) OptionsFunc {
	return func(o *r.UniversalOptions) {
		o.Password = password
	}
}

func WithDialTimeout(timeout time.Duration) OptionsFunc {
	return func(o *r.UniversalOptions) {
		o.DialTimeout = timeout
	}
}

func WithReadTimeout(timeout time.Duration) OptionsFunc {
	return func(o *r.UniversalOptions) {
		o.ReadTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) OptionsFunc {
	return func(o *r.UniversalOptions) {
		o.WriteTimeout = timeout
	}
}

func NewRedisClient(addr string, opts ...OptionsFunc) (*Client, error) {
	client := &Client{
		Ctx: context.Background(),
	}

	// 加入互斥锁，防止同时创建多个客户端
	client.Mutex.Lock()
	defer client.Mutex.Unlock()

	// 创建通用客户端，支持单节点、哨兵模式、Cluster集群模式
	uo := &r.UniversalOptions{
		Addrs:         strings.Split(addr, ","),
		RouteRandomly: true,
	}

	// 注册相关配置
	for _, registerOptFunc := range opts {
		registerOptFunc(uo)
	}
	client.Cli = r.NewUniversalClient(uo)
	return client, nil
}

type Client struct {
	Cli   r.UniversalClient // 通用客户端
	Ctx   context.Context
	Mutex sync.RWMutex
}

var (
	R *Client
)

func (rc *Client) Get(key string) string {
	var mes string
	var cmd *r.StringCmd

	cmd = rc.Cli.Get(rc.Ctx, key)

	if err := cmd.Err(); err != nil {
		mes = ""
	} else {
		mes = cmd.Val()
	}
	return mes
}

func (rc *Client) TTL(key string) time.Duration {
	return rc.Cli.TTL(rc.Ctx, key).Val()
}

func (rc *Client) GetRaw(key string) (bts []byte, err error) {
	bts, err = rc.Cli.Get(rc.Ctx, key).Bytes()

	if err != nil && err != r.Nil {
		return []byte{}, err
	}
	return bts, nil
}

func (rc *Client) MGet(keys ...string) ([]string, error) {
	var sliceCmd *r.SliceCmd
	sliceCmd = rc.Cli.MGet(rc.Ctx, keys...)

	if err := sliceCmd.Err(); err != nil && !errors.Is(err, r.Nil) {
		return []string{}, err
	}
	tmp := sliceCmd.Val()
	strSlice := make([]string, 0, len(tmp))
	for _, v := range tmp {
		if v != nil {
			strSlice = append(strSlice, v.(string))
		} else {
			strSlice = append(strSlice, "")
		}
	}
	return strSlice, nil
}

func (rc *Client) MGets(keys ...string) (ret []interface{}, err error) {
	ret, err = rc.Cli.MGet(rc.Ctx, keys...).Result()

	if err != nil && !errors.Is(err, r.Nil) {
		return []interface{}{}, err
	}
	return ret, nil
}

func (rc *Client) Set(key string, value interface{}, expire time.Duration) (bool, error) {
	var err error
	err = rc.Cli.Set(rc.Ctx, key, value, expire).Err()

	if err != nil {
		return false, err
	}
	return true, nil
}

// HGetAll 从redis获取hash的所有键值对
func (rc *Client) HGetAll(key string) map[string]string {
	var hash map[string]string
	var stringMapCmd *r.StringStringMapCmd
	stringMapCmd = rc.Cli.HGetAll(rc.Ctx, key)

	if err := stringMapCmd.Err(); err != nil && !errors.Is(err, r.Nil) {
		hash = make(map[string]string)
	} else {
		hash = stringMapCmd.Val()
	}

	return hash
}

// HGet 从redis获取hash单个值
func (rc *Client) HGet(key string, fields string) (string, error) {
	var stringCmd *r.StringCmd
	stringCmd = rc.Cli.HGet(rc.Ctx, key, fields)

	err := stringCmd.Err()
	if err != nil && !errors.Is(err, r.Nil) {
		return "", err
	}
	if errors.Is(err, r.Nil) {
		return "", nil
	}
	return stringCmd.Val(), nil
}

// HMGet 批量获取hash值
func (rc *Client) HMGet(key string, fileds ...string) []string {
	var sliceCmd *r.SliceCmd
	sliceCmd = rc.Cli.HMGet(rc.Ctx, key, fileds...)

	if err := sliceCmd.Err(); err != nil && !errors.Is(err, r.Nil) {
		return []string{}
	}
	tmp := sliceCmd.Val()
	strSlice := make([]string, 0, len(tmp))
	for _, v := range tmp {
		if v != nil {
			strSlice = append(strSlice, v.(string))
		} else {
			strSlice = append(strSlice, "")
		}
	}
	return strSlice
}

// HMGetMap 批量获取hash值，返回map
func (rc *Client) HMGetMap(key string, fields ...string) map[string]string {
	if len(fields) == 0 {
		return make(map[string]string)
	}

	var sliceCmd *r.SliceCmd
	sliceCmd = rc.Cli.HMGet(rc.Ctx, key, fields...)

	if err := sliceCmd.Err(); err != nil && !errors.Is(err, r.Nil) {
		return make(map[string]string)
	}

	tmp := sliceCmd.Val()
	hashRet := make(map[string]string, len(tmp))

	var tmpTagID string

	for k, v := range tmp {
		tmpTagID = fields[k]
		if v != nil {
			hashRet[tmpTagID] = v.(string)
		} else {
			hashRet[tmpTagID] = ""
		}
	}
	return hashRet
}

// HMSet 设置redis的hash
func (rc *Client) HMSet(key string, hash map[string]string, expire time.Duration) bool {
	if len(hash) > 0 {
		var err error
		err = rc.Cli.HMSet(rc.Ctx, key, hash).Err()

		if err != nil {
			return false
		}
		rc.Cli.Expire(rc.Ctx, key, expire)
		return true
	}
	return false
}

// HSet hset
func (rc *Client) HSet(key string, field string, value interface{}) bool {
	var err error
	err = rc.Cli.HSet(rc.Ctx, key, field, value).Err()
	if err != nil {
		return false
	}
	return true
}

// HDel ...
func (rc *Client) HDel(key string, field ...string) bool {
	var intCmd *r.IntCmd
	intCmd = rc.Cli.HDel(rc.Ctx, key, field...)

	if err := intCmd.Err(); err != nil {
		return false
	}

	return true
}

// SetWithErr ...
func (rc *Client) SetWithErr(key string, value interface{}, expire time.Duration) error {

	return rc.Cli.Set(rc.Ctx, key, value, expire).Err()
}

// SetNx 设置redis的string 如果键已存在
func (rc *Client) SetNx(key string, value interface{}, expiration time.Duration) bool {
	var res bool
	var err error
	res, err = rc.Cli.SetNX(rc.Ctx, key, value, expiration).Result()

	if err != nil {
		return false
	}

	return res
}

// SetNxWithErr 设置redis的string 如果键已存在
func (rc *Client) SetNxWithErr(key string, value interface{}, expiration time.Duration) (bool, error) {
	return rc.Cli.SetNX(rc.Ctx, key, value, expiration).Result()
}

// Incr redis自增
func (rc *Client) Incr(key string) bool {
	var err error
	err = rc.Cli.Incr(rc.Ctx, key).Err()

	if err != nil {
		return false
	}
	return true
}

// IncrWithErr ...
func (rc *Client) IncrWithErr(key string) (int64, error) {
	return rc.Cli.Incr(rc.Ctx, key).Result()
}

// IncrBy 将 key 所储存的值加上增量 increment 。
func (rc *Client) IncrBy(key string, increment int64) (int64, error) {
	var intCmd *r.IntCmd
	intCmd = rc.Cli.IncrBy(rc.Ctx, key, increment)

	if err := intCmd.Err(); err != nil {
		return 0, err
	}
	return intCmd.Val(), nil
}

// Decr redis自减
func (rc *Client) Decr(key string) bool {
	var err error
	err = rc.Cli.Decr(rc.Ctx, key).Err()

	if err != nil {
		return false
	}
	return true
}

// Type ...
func (rc *Client) Type(key string) (string, error) {
	var statusCmd *r.StatusCmd
	statusCmd = rc.Cli.Type(rc.Ctx, key)

	if err := statusCmd.Err(); err != nil {
		return "", err
	}
	return statusCmd.Val(), nil
}

// ZRevRange 倒序获取有序集合的部分数据
func (rc *Client) ZRevRange(key string, start, stop int64) ([]string, error) {
	var stringSliceCmd *r.StringSliceCmd

	stringSliceCmd = rc.Cli.ZRevRange(rc.Ctx, key, start, stop)
	if err := stringSliceCmd.Err(); err != nil && err != r.Nil {
		return []string{}, err
	}
	return stringSliceCmd.Val(), nil
}

// ZRevRangeWithScores ...
func (rc *Client) ZRevRangeWithScores(key string, start, stop int64) ([]r.Z, error) {
	var zSliceCmd *r.ZSliceCmd
	zSliceCmd = rc.Cli.ZRevRangeWithScores(rc.Ctx, key, start, stop)

	if err := zSliceCmd.Err(); err != nil && err != r.Nil {
		return []r.Z{}, err
	}
	return zSliceCmd.Val(), nil
}

// ZRange ...
func (rc *Client) ZRange(key string, start, stop int64) ([]string, error) {
	var stringSliceCmd *r.StringSliceCmd
	stringSliceCmd = rc.Cli.ZRange(rc.Ctx, key, start, stop)

	if err := stringSliceCmd.Err(); err != nil && err != r.Nil {
		return []string{}, err
	}
	return stringSliceCmd.Val(), nil
}

// ZRevRank ...
func (rc *Client) ZRevRank(key string, member string) (int64, error) {
	var intCmd *r.IntCmd

	intCmd = rc.Cli.ZRevRank(rc.Ctx, key, member)

	if err := intCmd.Err(); err != nil && err != r.Nil {
		return 0, err
	}
	return intCmd.Val(), nil
}

// ZRevRangeByScore ...
func (rc *Client) ZRevRangeByScore(key string, opt *r.ZRangeBy) (res []string, err error) {

	res, err = rc.Cli.ZRevRangeByScore(rc.Ctx, key, opt).Result()

	if err != nil && !errors.Is(err, r.Nil) {
		return []string{}, err
	}

	return res, nil
}

// ZRevRangeByScoreWithScores ...
func (rc *Client) ZRevRangeByScoreWithScores(key string, opt *r.ZRangeBy) (res []r.Z, err error) {
	res, err = rc.Cli.ZRevRangeByScoreWithScores(rc.Ctx, key, opt).Result()

	if err != nil && !errors.Is(err, r.Nil) {
		return []r.Z{}, err
	}

	return res, nil
}

// ZCard 获取有序集合的基数
func (rc *Client) nZCard(key string) (int64, error) {
	var intCmd *r.IntCmd
	intCmd = rc.Cli.ZCard(rc.Ctx, key)

	if err := intCmd.Err(); err != nil {
		return 0, err
	}
	return intCmd.Val(), nil
}

// ZScore 获取有序集合成员 member 的 score 值
func (rc *Client) ZScore(key string, member string) (float64, error) {
	var floatCmd *r.FloatCmd
	floatCmd = rc.Cli.ZScore(rc.Ctx, key, member)

	err := floatCmd.Err()
	if err != nil && !errors.Is(err, r.Nil) {
		return 0, err
	}
	return floatCmd.Val(), err
}

// ZAdd 将一个或多个 member 元素及其 score 值加入到有序集 key 当中
func (rc *Client) ZAdd(key string, members ...*r.Z) (int64, error) {
	var intCmd *r.IntCmd
	intCmd = rc.Cli.ZAdd(rc.Ctx, key, members...)

	if err := intCmd.Err(); err != nil && !errors.Is(err, r.Nil) {
		return 0, err
	}
	return intCmd.Val(), nil
}

// ZCount 返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。
func (rc *Client) ZCount(key string, min, max string) (int64, error) {
	var intCmd *r.IntCmd
	intCmd = rc.Cli.ZCount(rc.Ctx, key, min, max)

	if err := intCmd.Err(); err != nil && !errors.Is(err, r.Nil) {
		return 0, err
	}

	return intCmd.Val(), nil
}

// Del redis删除
func (rc *Client) Del(key string) int64 {
	var res int64
	var err error

	res, err = rc.Cli.Del(rc.Ctx, key).Result()

	if err != nil {
		return 0
	}
	return res
}

// DelWithErr ...
func (rc *Client) DelWithErr(key string) (int64, error) {
	return rc.Cli.Del(rc.Ctx, key).Result()
}

// HIncrBy 哈希field自增
func (rc *Client) HIncrBy(key string, field string, incr int) {
	rc.Cli.HIncrBy(rc.Ctx, key, field, int64(incr))
}

// Exists 键是否存在
func (rc *Client) Exists(key string) bool {
	var res int64
	var err error
	res, err = rc.Cli.Exists(rc.Ctx, key).Result()

	if err != nil {
		return false
	}
	return res == 1
}

// ExistsWithErr ...
func (rc *Client) ExistsWithErr(key string) (bool, error) {
	var res int64
	var err error
	res, err = rc.Cli.Exists(rc.Ctx, key).Result()

	if err != nil {
		return false, nil
	}
	return res == 1, nil
}

// LPush 将一个或多个值 value 插入到列表 key 的表头
func (rc *Client) LPush(key string, values ...interface{}) (int64, error) {
	var intCmd *r.IntCmd

	intCmd = rc.Cli.LPush(rc.Ctx, key, values...)

	if err := intCmd.Err(); err != nil {
		return 0, err
	}

	return intCmd.Val(), nil
}

// RPush 将一个或多个值 value 插入到列表 key 的表尾(最右边)。
func (rc *Client) RPush(key string, values ...interface{}) (int64, error) {
	var intCmd *r.IntCmd
	intCmd = rc.Cli.RPush(rc.Ctx, key, values...)

	if err := intCmd.Err(); err != nil {
		return 0, err
	}

	return intCmd.Val(), nil
}

// RPop 移除并返回列表 key 的尾元素。
func (rc *Client) RPop(key string) (string, error) {
	var stringCmd *r.StringCmd
	stringCmd = rc.Cli.RPop(rc.Ctx, key)

	if err := stringCmd.Err(); err != nil {
		return "", err
	}

	return stringCmd.Val(), nil
}

// LRange 获取列表指定范围内的元素
func (rc *Client) LRange(key string, start, stop int64) (res []string, err error) {
	res, err = rc.Cli.LRange(rc.Ctx, key, start, stop).Result()

	if err != nil {
		return []string{}, err
	}

	return res, nil
}

// LLen ...
func (rc *Client) LLen(key string) int64 {
	intCmd := rc.Cli.LLen(rc.Ctx, key)

	if err := intCmd.Err(); err != nil {
		return 0
	}

	return intCmd.Val()
}

// LLenWithErr ...
func (rc *Client) LLenWithErr(key string) (int64, error) {

	return rc.Cli.LLen(rc.Ctx, key).Result()
}

// LRem ...
func (rc *Client) LRem(key string, count int64, value interface{}) int64 {
	var intCmd *r.IntCmd
	intCmd = rc.Cli.LRem(rc.Ctx, key, count, value)

	if err := intCmd.Err(); err != nil {
		return 0
	}

	return intCmd.Val()
}

// LIndex ...
func (rc *Client) LIndex(key string, idx int64) (string, error) {

	return rc.Cli.LIndex(rc.Ctx, key, idx).Result()
}

// LTrim ...
func (rc *Client) LTrim(key string, start, stop int64) (string, error) {

	return rc.Cli.LTrim(rc.Ctx, key, start, stop).Result()
}

// ZRemRangeByRank 移除有序集合中给定的排名区间的所有成员
func (rc *Client) ZRemRangeByRank(key string, start, stop int64) (res int64, err error) {
	res, err = rc.Cli.ZRemRangeByRank(rc.Ctx, key, start, stop).Result()

	if err != nil {
		return 0, err
	}

	return res, nil
}

// Expire 设置过期时间
func (rc *Client) Expire(key string, expiration time.Duration) (res bool, err error) {
	res, err = rc.Cli.Expire(rc.Ctx, key, expiration).Result()

	if err != nil {
		return false, err
	}

	return res, err
}

// ZRem 从zset中移除变量
func (rc *Client) ZRem(key string, members ...interface{}) (res int64, err error) {
	res, err = rc.Cli.ZRem(rc.Ctx, key, members...).Result()

	if err != nil {
		return 0, err
	}
	return res, nil
}

// SAdd 向set中添加成员
func (rc *Client) SAdd(key string, member ...interface{}) (int64, error) {
	var intCmd *r.IntCmd
	intCmd = rc.Cli.SAdd(rc.Ctx, key, member...)

	if err := intCmd.Err(); err != nil {
		return 0, err
	}
	return intCmd.Val(), nil
}

// SMembers 返回set的全部成员
func (rc *Client) SMembers(key string) ([]string, error) {
	var stringSliceCmd *r.StringSliceCmd
	stringSliceCmd = rc.Cli.SMembers(rc.Ctx, key)

	if err := stringSliceCmd.Err(); err != nil {
		return []string{}, err
	}
	return stringSliceCmd.Val(), nil
}

// SIsMember ...
func (rc *Client) SIsMember(key string, member interface{}) (bool, error) {
	var boolCmd *r.BoolCmd
	boolCmd = rc.Cli.SIsMember(rc.Ctx, key, member)

	if err := boolCmd.Err(); err != nil {
		return false, err
	}
	return boolCmd.Val(), nil
}

// HKeys 获取hash的所有域
func (rc *Client) HKeys(key string) []string {
	var stringSliceCmd *r.StringSliceCmd
	stringSliceCmd = rc.Cli.HKeys(rc.Ctx, key)

	if err := stringSliceCmd.Err(); err != nil && err != r.Nil {
		return []string{}
	}
	return stringSliceCmd.Val()
}

// HLen 获取hash的长度
func (rc *Client) HLen(key string) int64 {
	var intCmd *r.IntCmd
	intCmd = rc.Cli.HLen(rc.Ctx, key)

	if err := intCmd.Err(); err != nil && err != r.Nil {
		return 0
	}
	return intCmd.Val()
}

// GeoAdd 写入地理位置
func (rc *Client) GeoAdd(key string, location *r.GeoLocation) (res int64, err error) {
	res, err = rc.Cli.GeoAdd(rc.Ctx, key, location).Result()

	if err != nil {
		return 0, err
	}

	return res, nil
}

// GeoRadius 根据经纬度查询列表
func (rc *Client) GeoRadius(key string, longitude, latitude float64, query *r.GeoRadiusQuery) (res []r.GeoLocation, err error) {
	res, err = rc.Cli.GeoRadius(rc.Ctx, key, longitude, latitude, query).Result()

	if err != nil {
		return []r.GeoLocation{}, err
	}

	return res, nil
}

// Eval 执行lua脚本
func (rc *Client) Eval(script string, keys []string, args ...interface{}) (res interface{}, err error) {
	res, err = rc.Cli.Eval(rc.Ctx, script, keys, args...).Result()

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (rc *Client) PopSetMembers(key string) (res []string, err error) {
	luaScript := `
local set_members = redis.call('SMEMBERS', KEYS[1])
redis.call('DEL', KEYS[1])
return set_members
`
	srcMembers, err := rc.Eval(luaScript, []string{key})
	if err != nil {
		return
	}

	members, ok := srcMembers.([]interface{})
	if !ok {
		return res, errors.New("failed to convert members to []interface{}")
	}

	for _, member := range members {
		res = append(res, member.(string))
	}
	return
}

func (rc *Client) Close() (err error) {
	err = rc.Cli.Close()

	return
}

// DelByPrefix 删除指定前缀的所有键，忽略错误
func (rc *Client) DelByPrefix(prefix string) int64 {
	var cursor uint64
	var count int64

	for {
		keys, nextCursor, _ := rc.Cli.Scan(rc.Ctx, cursor, prefix+"*", 100).Result()
		cursor = nextCursor

		if len(keys) > 0 {
			n, _ := rc.Cli.Del(rc.Ctx, keys...).Result()
			count += n
		}

		if cursor == 0 {
			break
		}
	}

	return count
}
