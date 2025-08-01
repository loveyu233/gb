package gb

import (
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

type item struct {
	Member    string      `json:"member"`
	Score     json.Number `json:"score"`
	Rank      int64       `json:"rank"`
	HashValue string      `json:"hash_value"`
}

type LuaRespData struct {
	Range  []item `json:"range"`
	Target item   `json:"target"`
}

// LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreAndGetHashValue 在LuaRedisZSetGetTargetKeyAndStartToEndRankByScore的基础上添加了获取item的hash_value
func LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreAndGetHashValue(zSetKey, hashKey string, start, end int64, targetMember string, descending bool) (*LuaRespData, error) {
	lua := `local key = KEYS[1]
			local hash_key = KEYS[2]  -- New parameter for hash key
			local start_pos = tonumber(ARGV[1])
			local end_pos = tonumber(ARGV[2])
			local target_member = ARGV[3]
			local descending = ARGV[4] == "true"
			
			local range_members, target_rank, target_score
			
			local exists = redis.call('exists', key)
			if exists == 0 then
				local result = {
					range = {},
					target = {
						member = target_member,
						score = 0,
						rank = -1,
						hash_value = nil
					}
				}
				return cjson.encode(result)
			end
			
			-- 根据排序方向选择不同的命令
			if descending then
				-- 从大到小排序
				range_members = redis.call('ZREVRANGE', key, start_pos, end_pos, 'WITHSCORES')
				target_rank = redis.call('ZREVRANK', key, target_member)
				target_score = redis.call('ZSCORE', key, target_member)
			else
				-- 从小到大排序
				range_members = redis.call('ZRANGE', key, start_pos, end_pos, 'WITHSCORES')
				target_rank = redis.call('ZRANK', key, target_member)
				target_score = redis.call('ZSCORE', key, target_member)
			end
			
			-- 处理空值情况
			local result = {
				range = {},
				target = {
					member = target_member,
					score = target_score and tonumber(target_score) or 0,
					rank = target_rank and tonumber(target_rank) or -1,
					hash_value = redis.call('HGET', hash_key, target_member) or nil
				}
			}
			
			-- 如果没有范围数据，直接返回空结果
			if #range_members == 0 then
				return cjson.encode(result)
			end
			
			-- 处理范围数据
			for i = 1, #range_members, 2 do
				local member = range_members[i]
				local score = range_members[i+1]
				local rank
				
				-- 根据排序方向获取排名
				if descending then
					rank = redis.call('ZREVRANK', key, member)
				else
					rank = redis.call('ZRANK', key, member)
				end
				
				table.insert(result.range, {
					member = member,
					score = tonumber(score) or 0,
					rank = rank and tonumber(rank) or -1,
					hash_value = redis.call('HGET', hash_key, member) or nil
				})
			end
			
			return cjson.encode(result)`

	// 将布尔值转换为字符串传递给 Lua
	descendingStr := "false"
	if descending {
		descendingStr = "true"
	}

	result, err := redis.NewScript(lua).Run(context.Background(), RedisClient.UniversalClient, []string{zSetKey, hashKey}, start, end, targetMember, descendingStr).Result()
	if err != nil {
		return nil, err
	}

	var luaRespData *LuaRespData
	err = json.Unmarshal([]byte(result.(string)), &luaRespData)
	if err != nil {
		if err.Error() == "json: cannot unmarshal object into Go struct field LuaRespData.range of type []gb.item" {
			return luaRespData, nil
		}
		return nil, err
	}

	return luaRespData, nil
}

func LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreAndGetHashValueDesc(zSetKey, hashKey string, start, end int64, targetMember string) (*LuaRespData, error) {
	return LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreAndGetHashValue(zSetKey, hashKey, start, end, targetMember, true)
}

func LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreAndGetHashValueAsc(zSetKey, hashKey string, start, end int64, targetMember string) (*LuaRespData, error) {
	return LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreAndGetHashValue(zSetKey, hashKey, start, end, targetMember, false)
}

// LuaRedisZSetGetTargetKeyAndStartToEndRankByScore 同时获取指定key的从start到end排名的数据和targetMember的排名和分数
// key: Redis ZSet 的键名
// start: 起始排名
// end: 结束排名
// targetMember: 目标成员
// descending: true表示从大到小排序(ZREVRANGE/ZREVRANK), false表示从小到大排序(ZRANGE/ZRANK)
func LuaRedisZSetGetTargetKeyAndStartToEndRankByScore(key string, start, end int64, targetMember string, descending bool) (*LuaRespData, error) {
	lua := `local key = KEYS[1]
			local start_pos = tonumber(ARGV[1])
			local end_pos = tonumber(ARGV[2])
			local target_member = ARGV[3]
			local descending = ARGV[4] == "true"
			
			local range_members, target_rank, target_score

			local exists = redis.call('exists', key)
			if exists == 0 then
				local result = {
					range = {},
					target = {
						member = target_member,
						score = 0,
						rank = -1
					}
				}
				return cjson.encode(result)
			end
			
			-- 根据排序方向选择不同的命令
			if descending then
				-- 从大到小排序
				range_members = redis.call('ZREVRANGE', key, start_pos, end_pos, 'WITHSCORES')
				target_rank = redis.call('ZREVRANK', key, target_member)
				target_score = redis.call('ZSCORE', key, target_member)
			else
				-- 从小到大排序
				range_members = redis.call('ZRANGE', key, start_pos, end_pos, 'WITHSCORES')
				target_rank = redis.call('ZRANK', key, target_member)
				target_score = redis.call('ZSCORE', key, target_member)
			end
			
			-- 处理空值情况
			local result = {
				range = {},
				target = {
					member = target_member,
					score = target_score and tonumber(target_score) or 0,
					rank = target_rank and tonumber(target_rank) or -1
				}
			}
			
			-- 如果没有范围数据，直接返回空结果
			if #range_members == 0 then
				return cjson.encode(result)
			end
			
			-- 处理范围数据
			for i = 1, #range_members, 2 do
				local member = range_members[i]
				local score = range_members[i+1]
				local rank
				
				-- 根据排序方向获取排名
				if descending then
					rank = redis.call('ZREVRANK', key, member)
				else
					rank = redis.call('ZRANK', key, member)
				end
				
				table.insert(result.range, {
					member = member,
					score = tonumber(score) or 0,
					rank = rank and tonumber(rank) or -1
				})
			end
			
			return cjson.encode(result)`

	// 将布尔值转换为字符串传递给 Lua
	descendingStr := "false"
	if descending {
		descendingStr = "true"
	}

	result, err := redis.NewScript(lua).Run(context.Background(), RedisClient.UniversalClient, []string{key}, start, end, targetMember, descendingStr).Result()
	if err != nil {
		return nil, err
	}

	var luaRespData *LuaRespData
	err = json.Unmarshal([]byte(result.(string)), &luaRespData)
	if err != nil {
		if err.Error() == "json: cannot unmarshal object into Go struct field LuaRespData.range of type []gb.item" {
			return luaRespData, nil
		}
		return nil, err
	}

	return luaRespData, nil
}

// LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreDesc 新增一个降序排列的便捷函数
func LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreDesc(key string, start, end int64, targetMember string) (*LuaRespData, error) {
	return LuaRedisZSetGetTargetKeyAndStartToEndRankByScore(key, start, end, targetMember, true)
}

// LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreAsc 新增一个升序排列的便捷函数
func LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreAsc(key string, start, end int64, targetMember string) (*LuaRespData, error) {
	return LuaRedisZSetGetTargetKeyAndStartToEndRankByScore(key, start, end, targetMember, false)
}

type MemberInfo struct {
	Member    string      `json:"member"`
	Score     json.Number `json:"score"`
	Rank      int64       `json:"rank"`
	Exists    bool        `json:"exists"`
	HashValue string      `json:"hash_value"`
}

// LuaRedisZSetGetMemberScoreAndRankAndGetHashValue 获取指定 member 的 score 和 rank 以及hash value
// key: Redis ZSet 的 key
// member: 要查询的成员
// descending: true 表示从大到小排序(ZREVRANK), false 表示从小到大排序(ZRANK)
func LuaRedisZSetGetMemberScoreAndRankAndGetHashValue(zSetKey, hashKey string, member string, descending bool) (*MemberInfo, error) {
	lua := `local key = KEYS[1]
				local hash_key = KEYS[2]  -- New parameter for hash key
				local member = ARGV[1]
				local descending = ARGV[2] == "true"
				
				-- 获取成员的分数
				local score = redis.call('ZSCORE', key, member)
				
				-- 如果成员不存在，返回空结果
				if not score then
					return cjson.encode({
						member = member,
						score = 0,
						rank = -1,
						exists = false,
						hash_value = nil
					})
				end
				
				-- 根据排序方向获取排名
				local rank
				if descending then
					-- 从大到小排序，使用 ZREVRANK (分数高的排名小)
					rank = redis.call('ZREVRANK', key, member)
				else
					-- 从小到大排序，使用 ZRANK (分数低的排名小)
					rank = redis.call('ZRANK', key, member)
				end

				return cjson.encode({
					member = member,
					score = tonumber(score),
					rank = tonumber(rank),
					exists = true,
					hash_value = redis.call('HGET', hash_key, member) or nil
				})`

	// 将布尔值转换为字符串传递给 Lua
	descendingStr := "false"
	if descending {
		descendingStr = "true"
	}

	result, err := redis.NewScript(lua).Run(context.Background(), RedisClient.UniversalClient, []string{zSetKey, hashKey}, member, descendingStr).Result()
	if err != nil {
		return nil, err
	}

	var memberInfo *MemberInfo
	err = json.Unmarshal([]byte(result.(string)), &memberInfo)
	if err != nil {
		return nil, err
	}

	return memberInfo, nil
}

// LuaRedisZSetGetMemberScoreAndRankAndGetHashValueDesc LuaRedisZSetGetMemberScoreAndRankAndGetHashValue的快捷降序
func LuaRedisZSetGetMemberScoreAndRankAndGetHashValueDesc(zSetKey, hashKey string, member string) (*MemberInfo, error) {
	return LuaRedisZSetGetMemberScoreAndRankAndGetHashValue(zSetKey, hashKey, member, true)
}

// LuaRedisZSetGetMemberScoreAndRankAndGetHashValueAsc 	LuaRedisZSetGetMemberScoreAndRankAndGetHashValue的快捷升序
func LuaRedisZSetGetMemberScoreAndRankAndGetHashValueAsc(zSetKey, hashKey string, member string) (*MemberInfo, error) {
	return LuaRedisZSetGetMemberScoreAndRankAndGetHashValue(zSetKey, hashKey, member, false)
}

// LuaRedisZSetGetMemberScoreAndRank 获取指定 member 的 score 和 rank
// key: Redis ZSet 的 key
// member: 要查询的成员
// descending: true 表示从大到小排序(ZREVRANK), false 表示从小到大排序(ZRANK)
func LuaRedisZSetGetMemberScoreAndRank(key string, member string, descending bool) (*MemberInfo, error) {
	lua := `local key = KEYS[1]
				local member = ARGV[1]
				local descending = ARGV[2] == "true"
				
				-- 获取成员的分数
				local score = redis.call('ZSCORE', key, member)
				
				-- 如果成员不存在，返回空结果
				if not score then
					return cjson.encode({
						member = member,
						score = 0,
						rank = -1,
						exists = false
					})
				end
				
				-- 根据排序方向获取排名
				local rank
				if descending then
					-- 从大到小排序，使用 ZREVRANK (分数高的排名小)
					rank = redis.call('ZREVRANK', key, member)
				else
					-- 从小到大排序，使用 ZRANK (分数低的排名小)
					rank = redis.call('ZRANK', key, member)
				end
				
				return cjson.encode({
					member = member,
					score = tonumber(score),
					rank = tonumber(rank),
					exists = true
				})`

	// 将布尔值转换为字符串传递给 Lua
	descendingStr := "false"
	if descending {
		descendingStr = "true"
	}

	result, err := redis.NewScript(lua).Run(context.Background(), RedisClient.UniversalClient, []string{key}, member, descendingStr).Result()
	if err != nil {
		return nil, err
	}

	var memberInfo *MemberInfo
	err = json.Unmarshal([]byte(result.(string)), &memberInfo)
	if err != nil {
		return nil, err
	}

	return memberInfo, nil
}

// LuaRedisZSetGetMemberScoreAndRankDesc 获取指定 member 的 score 和 rank,快捷降序
func LuaRedisZSetGetMemberScoreAndRankDesc(key string, member string) (*MemberInfo, error) {
	return LuaRedisZSetGetMemberScoreAndRank(key, member, true)
}

// LuaRedisZSetGetMemberScoreAndRankAsc 获取指定 member 的 score 和 rank,快捷升序
func LuaRedisZSetGetMemberScoreAndRankAsc(key string, member string) (*MemberInfo, error) {
	return LuaRedisZSetGetMemberScoreAndRank(key, member, false)
}

// LuaRedisZSetGetMultipleMembersScoreAndRankAndHashValues 批量获取多个 member 的 score 和 rank 以及对应hashKey的value
func LuaRedisZSetGetMultipleMembersScoreAndRankAndHashValues(zSetKey, hashKey string, members []string, descending bool) ([]*MemberInfo, error) {
	if len(members) == 0 {
		return []*MemberInfo{}, nil
	}

	lua := `local key = KEYS[1]
				local hash_key = KEYS[2]  -- New parameter for hash key
				local descending = ARGV[1] == "true"
				local members = {}
				
				-- 从 ARGV[2] 开始是 member 列表
				for i = 2, #ARGV do
					table.insert(members, ARGV[i])
				end
				
				local results = {}
				
				for _, member in ipairs(members) do
					-- 获取成员的分数
					local score = redis.call('ZSCORE', key, member)
					
					if not score then
						-- 成员不存在
						table.insert(results, {
							member = member,
							score = 0,
							rank = -1,
							exists = false,
							hash_value = nil
						})
					else
						-- 根据排序方向获取排名
						local rank
						if descending then
							rank = redis.call('ZREVRANK', key, member)
						else
							rank = redis.call('ZRANK', key, member)
						end
						
						table.insert(results, {
							member = member,
							score = tonumber(score),
							rank = tonumber(rank),
							exists = true,
							hash_value = redis.call('HGET', hash_key, member) or nil
						})
					end
				end
				
				return cjson.encode(results)`

	// 构建参数列表
	args := make([]interface{}, len(members)+1)
	args[0] = "true"
	if !descending {
		args[0] = "false"
	}
	for i, member := range members {
		args[i+1] = member
	}

	result, err := redis.NewScript(lua).Run(context.Background(), RedisClient.UniversalClient, []string{zSetKey, hashKey}, args...).Result()
	if err != nil {
		return nil, err
	}

	var memberInfos []*MemberInfo
	err = json.Unmarshal([]byte(result.(string)), &memberInfos)
	if err != nil {
		return nil, err
	}

	return memberInfos, nil
}

// LuaRedisZSetGetMultipleMembersScoreAndRankAndHashValuesDesc LuaRedisZSetGetMultipleMembersScoreAndRankAndHashValues的desc版本
func LuaRedisZSetGetMultipleMembersScoreAndRankAndHashValuesDesc(zSetKey, hashKey string, members []string) ([]*MemberInfo, error) {
	return LuaRedisZSetGetMultipleMembersScoreAndRankAndHashValues(zSetKey, hashKey, members, true)
}

// LuaRedisZSetGetMultipleMembersScoreAndRankAndHashValuesAsc LuaRedisZSetGetMultipleMembersScoreAndRankAndHashValues的asc版本
func LuaRedisZSetGetMultipleMembersScoreAndRankAndHashValuesAsc(zSetKey, hashKey string, members []string) ([]*MemberInfo, error) {
	return LuaRedisZSetGetMultipleMembersScoreAndRankAndHashValues(zSetKey, hashKey, members, false)
}

// LuaRedisZSetGetMultipleMembersScoreAndRank 批量获取多个 member 的 score 和 rank
// key: Redis ZSet 的 key
// members: 要查询的成员列表
// descending: true 表示从大到小排序, false 表示从小到大排序
func LuaRedisZSetGetMultipleMembersScoreAndRank(key string, members []string, descending bool) ([]*MemberInfo, error) {
	if len(members) == 0 {
		return []*MemberInfo{}, nil
	}

	lua := `local key = KEYS[1]
				local descending = ARGV[1] == "true"
				local members = {}
				
				-- 从 ARGV[2] 开始是 member 列表
				for i = 2, #ARGV do
					table.insert(members, ARGV[i])
				end
				
				local results = {}
				
				for _, member in ipairs(members) do
					-- 获取成员的分数
					local score = redis.call('ZSCORE', key, member)
					
					if not score then
						-- 成员不存在
						table.insert(results, {
							member = member,
							score = 0,
							rank = -1,
							exists = false
						})
					else
						-- 根据排序方向获取排名
						local rank
						if descending then
							rank = redis.call('ZREVRANK', key, member)
						else
							rank = redis.call('ZRANK', key, member)
						end
						
						table.insert(results, {
							member = member,
							score = tonumber(score),
							rank = tonumber(rank),
							exists = true
						})
					end
				end
				
				return cjson.encode(results)`

	// 构建参数列表
	args := make([]interface{}, len(members)+1)
	args[0] = "true"
	if !descending {
		args[0] = "false"
	}
	for i, member := range members {
		args[i+1] = member
	}

	result, err := redis.NewScript(lua).Run(context.Background(), RedisClient.UniversalClient, []string{key}, args...).Result()
	if err != nil {
		return nil, err
	}

	var memberInfos []*MemberInfo
	err = json.Unmarshal([]byte(result.(string)), &memberInfos)
	if err != nil {
		return nil, err
	}

	return memberInfos, nil
}

// LuaRedisZSetGetMultipleMembersScoreAndRankDesc 批量获取多个 member 的 score 和 rank,快捷降序
func LuaRedisZSetGetMultipleMembersScoreAndRankDesc(key string, members []string) ([]*MemberInfo, error) {
	return LuaRedisZSetGetMultipleMembersScoreAndRank(key, members, true)
}

// LuaRedisZSetGetMultipleMembersScoreAndRankAsc 批量获取多个 member 的 score 和 rank,快捷升序
func LuaRedisZSetGetMultipleMembersScoreAndRankAsc(key string, members []string) ([]*MemberInfo, error) {
	return LuaRedisZSetGetMultipleMembersScoreAndRank(key, members, false)
}

// 1. 分布式锁相关

// LuaRedisDistributedLock 获取分布式锁
// key: 锁的键名, value: 锁的值(通常是唯一ID), expireSeconds: 过期时间(秒)
// 返回: true表示获取锁成功, false表示获取锁失败
func LuaRedisDistributedLock(key, value string, expireSeconds int64) (bool, error) {
	lua := `if redis.call('SET', KEYS[1], ARGV[1], 'NX', 'EX', ARGV[2]) then
				return 1
			else
				return 0
			end`

	result, err := redis.NewScript(lua).Run(context.Background(), RedisClient.UniversalClient, []string{key}, value, expireSeconds).Result()
	if err != nil {
		return false, err
	}
	return result.(int64) == 1, nil
}

// LuaRedisDistributedUnlock 释放分布式锁
// key: 锁的键名, value: 锁的值
// 返回: true表示释放成功, false表示释放失败(锁不存在或值不匹配)
func LuaRedisDistributedUnlock(key, value string) (bool, error) {
	lua := `if redis.call('GET', KEYS[1]) == ARGV[1] then
				return redis.call('DEL', KEYS[1])
			else
				return 0
			end`

	result, err := redis.NewScript(lua).Run(context.Background(), RedisClient.UniversalClient, []string{key}, value).Result()
	if err != nil {
		return false, err
	}
	return result.(int64) == 1, nil
}

// 2. 限流相关

// LuaRedisRateLimit 滑动窗口限流
// key: 限流键名, window: 时间窗口(秒), limit: 限制次数
// 返回: 剩余可用次数, 如果返回-1表示超过限制
func LuaRedisRateLimit(key string, window, limit int64) (int64, error) {
	lua := `local key = KEYS[1]
			local window = tonumber(ARGV[1])
			local limit = tonumber(ARGV[2])
			local current_time = redis.call('TIME')[1]
			
			-- 清理过期数据
			redis.call('ZREMRANGEBYSCORE', key, 0, current_time - window)
			
			-- 获取当前窗口内的请求数
			local current_requests = redis.call('ZCARD', key)
			
			if current_requests < limit then
				-- 添加当前请求
				redis.call('ZADD', key, current_time, current_time)
				redis.call('EXPIRE', key, window)
				return limit - current_requests - 1
			else
				return -1
			end`

	result, err := redis.NewScript(lua).Run(context.Background(), RedisClient.UniversalClient, []string{key}, window, limit).Result()
	if err != nil {
		return -1, err
	}
	return result.(int64), nil
}

// 3. 计数器相关

type CounterResult struct {
	CurrentValue int64 `json:"current_value"`
	IsSuccess    bool  `json:"is_success"`
}

// LuaRedisIncrWithLimit 带限制的计数器递增
// key: 计数器键名, increment: 递增值, maxValue: 最大值, expireSeconds: 过期时间
func LuaRedisIncrWithLimit(key string, increment, maxValue, expireSeconds int64) (*CounterResult, error) {
	lua := `local key = KEYS[1]
			local increment = tonumber(ARGV[1])
			local max_value = tonumber(ARGV[2])
			local expire_seconds = tonumber(ARGV[3])
			
			local current = redis.call('GET', key)
			if not current then
				current = 0
			else
				current = tonumber(current)
			end
			
			if current + increment <= max_value then
				local new_value = redis.call('INCRBY', key, increment)
				redis.call('EXPIRE', key, expire_seconds)
				return cjson.encode({
					current_value = new_value,
					is_success = true
				})
			else
				return cjson.encode({
					current_value = current,
					is_success = false
				})
			end`

	result, err := redis.NewScript(lua).Run(context.Background(), RedisClient.UniversalClient, []string{key}, increment, maxValue, expireSeconds).Result()
	if err != nil {
		return nil, err
	}

	var counterResult *CounterResult
	err = json.Unmarshal([]byte(result.(string)), &counterResult)
	return counterResult, err
}

// 4. 队列相关

// LuaRedisQueuePushWithLimit 带限制的队列推送
// key: 队列键名, value: 要推送的值, maxLength: 队列最大长度
// 返回: 队列当前长度, 如果返回-1表示队列已满
func LuaRedisQueuePushWithLimit(key, value string, maxLength int64) (int64, error) {
	lua := `local key = KEYS[1]
			local value = ARGV[1]
			local max_length = tonumber(ARGV[2])
			
			local current_length = redis.call('LLEN', key)
			
			if current_length < max_length then
				return redis.call('LPUSH', key, value)
			else
				return -1
			end`

	result, err := redis.NewScript(lua).Run(context.Background(), RedisClient.UniversalClient, []string{key}, value, maxLength).Result()
	if err != nil {
		return -1, err
	}
	return result.(int64), nil
}

// 5. 缓存相关

// LuaRedisSetWithVersion 带版本控制的设置
// key: 缓存键名, value: 值, version: 版本号, expireSeconds: 过期时间
// 返回: true表示设置成功, false表示版本冲突
func LuaRedisSetWithVersion(key, value string, version, expireSeconds int64) (bool, error) {
	lua := `local key = KEYS[1]
			local value = ARGV[1]
			local version = tonumber(ARGV[2])
			local expire_seconds = tonumber(ARGV[3])
			
			local version_key = key .. ':version'
			local current_version = redis.call('GET', version_key)
			
			if not current_version or tonumber(current_version) < version then
				redis.call('SET', key, value, 'EX', expire_seconds)
				redis.call('SET', version_key, version, 'EX', expire_seconds)
				return 1
			else
				return 0
			end`

	result, err := redis.NewScript(lua).Run(context.Background(), RedisClient.UniversalClient, []string{key}, value, version, expireSeconds).Result()
	if err != nil {
		return false, err
	}
	return result.(int64) == 1, nil
}

// 6. 库存扣减

type StockResult struct {
	Success        bool  `json:"success"`
	RemainingStock int64 `json:"remaining_stock"`
}

// LuaRedisDecrStock 原子扣减库存
// key: 库存键名, quantity: 扣减数量
func LuaRedisDecrStock(key string, quantity int64) (*StockResult, error) {
	lua := `local key = KEYS[1]
			local quantity = tonumber(ARGV[1])
			
			local stock = redis.call('GET', key)
			if not stock then
				return cjson.encode({
					success = false,
					remaining_stock = 0
				})
			end
			
			stock = tonumber(stock)
			if stock >= quantity then
				local remaining = redis.call('DECRBY', key, quantity)
				return cjson.encode({
					success = true,
					remaining_stock = remaining
				})
			else
				return cjson.encode({
					success = false,
					remaining_stock = stock
				})
			end`

	result, err := redis.NewScript(lua).Run(context.Background(), RedisClient.UniversalClient, []string{key}, quantity).Result()
	if err != nil {
		return nil, err
	}

	var stockResult *StockResult
	err = json.Unmarshal([]byte(result.(string)), &stockResult)
	return stockResult, err
}

// 7. HyperLogLog 去重计数

// LuaRedisHLLAddAndCount HyperLogLog 添加元素并返回估计数量
// key: HLL键名, elements: 要添加的元素列表
func LuaRedisHLLAddAndCount(key string, elements []string) (int64, error) {
	lua := `local key = KEYS[1]
			local elements = {}
			
			for i = 1, #ARGV do
				table.insert(elements, ARGV[i])
			end
			
			if #elements > 0 then
				redis.call('PFADD', key, unpack(elements))
			end
			
			return redis.call('PFCOUNT', key)`

	args := make([]interface{}, len(elements))
	for i, element := range elements {
		args[i] = element
	}

	result, err := redis.NewScript(lua).Run(context.Background(), RedisClient.UniversalClient, []string{key}, args...).Result()
	if err != nil {
		return 0, err
	}
	return result.(int64), nil
}

// 8. 排行榜相关

type LeaderboardMember struct {
	Member string      `json:"member"`
	Score  json.Number `json:"score"`
	Rank   int64       `json:"rank"`
}

// LuaRedisLeaderboardIncr 排行榜分数递增并返回排名信息
// key: 排行榜键名, member: 成员名, increment: 递增分数
func LuaRedisLeaderboardIncr(key, member string, increment float64) (*LeaderboardMember, error) {
	lua := `local key = KEYS[1]
			local member = ARGV[1]
			local increment = tonumber(ARGV[2])
			
			local new_score = redis.call('ZINCRBY', key, increment, member)
			local rank = redis.call('ZREVRANK', key, member)
			
			return cjson.encode({
				member = member,
				score = tonumber(new_score),
				rank = tonumber(rank)
			})`

	result, err := redis.NewScript(lua).Run(context.Background(), RedisClient.UniversalClient, []string{key}, member, increment).Result()
	if err != nil {
		return nil, err
	}

	var leaderboardMember *LeaderboardMember
	err = json.Unmarshal([]byte(result.(string)), &leaderboardMember)
	return leaderboardMember, err
}

// 9. 延迟队列

type DelayedMessage struct {
	ID      string `json:"id"`
	Payload string `json:"payload"`
	Score   int64  `json:"score"`
}

// LuaRedisDelayQueuePop 从延迟队列中弹出到期的消息
// key: 延迟队列键名, currentTime: 当前时间戳, limit: 最多弹出数量
func LuaRedisDelayQueuePop(key string, currentTime int64, limit int64) ([]*DelayedMessage, error) {
	lua := `local key = KEYS[1]
			local current_time = tonumber(ARGV[1])
			local limit = tonumber(ARGV[2])
			
			-- 获取到期的消息
			local messages = redis.call('ZRANGEBYSCORE', key, 0, current_time, 'WITHSCORES', 'LIMIT', 0, limit)
			
			if #messages == 0 then
				return cjson.encode({})
			end
			
			-- 构建结果
			local results = {}
			local members_to_remove = {}
			
			for i = 1, #messages, 2 do
				local payload = messages[i]
				local score = messages[i+1]
				
				table.insert(results, {
					id = payload,
					payload = payload,
					score = tonumber(score)
				})
				table.insert(members_to_remove, payload)
			end
			
			-- 删除已处理的消息
			if #members_to_remove > 0 then
				redis.call('ZREM', key, unpack(members_to_remove))
			end
			
			return cjson.encode(results)`

	result, err := redis.NewScript(lua).Run(context.Background(), RedisClient.UniversalClient, []string{key}, currentTime, limit).Result()
	if err != nil {
		return nil, err
	}

	var messages []*DelayedMessage
	err = json.Unmarshal([]byte(result.(string)), &messages)
	return messages, err
}

// 10. 布隆过滤器模拟 (使用多个 Hash)

// LuaRedisBloomAdd 布隆过滤器添加元素
// key: 布隆过滤器键名, element: 要添加的元素
func LuaRedisBloomAdd(key, element string) error {
	lua := `local key = KEYS[1]
			local element = ARGV[1]
			
			-- 使用多个哈希函数模拟布隆过滤器
			local hash1 = redis.sha1hex(element .. '1') % 1000000
			local hash2 = redis.sha1hex(element .. '2') % 1000000
			local hash3 = redis.sha1hex(element .. '3') % 1000000
			
			redis.call('SETBIT', key, hash1, 1)
			redis.call('SETBIT', key, hash2, 1)
			redis.call('SETBIT', key, hash3, 1)
			
			return 'OK'`

	_, err := redis.NewScript(lua).Run(context.Background(), RedisClient.UniversalClient, []string{key}, element).Result()
	return err
}

// LuaRedisBloomExists 检查布隆过滤器中是否存在元素
// key: 布隆过滤器键名, element: 要检查的元素
func LuaRedisBloomExists(key, element string) (bool, error) {
	lua := `local key = KEYS[1]
			local element = ARGV[1]
			
			-- 使用相同的哈希函数
			local hash1 = redis.sha1hex(element .. '1') % 1000000
			local hash2 = redis.sha1hex(element .. '2') % 1000000
			local hash3 = redis.sha1hex(element .. '3') % 1000000
			
			local bit1 = redis.call('GETBIT', key, hash1)
			local bit2 = redis.call('GETBIT', key, hash2)
			local bit3 = redis.call('GETBIT', key, hash3)
			
			if bit1 == 1 and bit2 == 1 and bit3 == 1 then
				return 1
			else
				return 0
			end`

	result, err := redis.NewScript(lua).Run(context.Background(), RedisClient.UniversalClient, []string{key}, element).Result()
	if err != nil {
		return false, err
	}
	return result.(int64) == 1, nil
}
