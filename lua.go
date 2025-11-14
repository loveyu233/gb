package gb

import (
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
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

// LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreAndGetHashValue 方法用于处理LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreAndGetHashValue相关逻辑。
func (r *RedisConfig) LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreAndGetHashValue(zSetKey, hashKey string, start, end int64, targetMember string, descending bool) (*LuaRespData, error) {
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

	result, err := redis.NewScript(lua).Run(context.Background(), r.UniversalClient, []string{zSetKey, hashKey}, start, end, targetMember, descendingStr).Result()
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

// LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreAndGetHashValueDesc 方法用于处理LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreAndGetHashValueDesc相关逻辑。
func (r *RedisConfig) LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreAndGetHashValueDesc(zSetKey, hashKey string, start, end int64, targetMember string) (*LuaRespData, error) {
	return r.LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreAndGetHashValue(zSetKey, hashKey, start, end, targetMember, true)
}

// LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreAndGetHashValueAsc 方法用于处理LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreAndGetHashValueAsc相关逻辑。
func (r *RedisConfig) LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreAndGetHashValueAsc(zSetKey, hashKey string, start, end int64, targetMember string) (*LuaRespData, error) {
	return r.LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreAndGetHashValue(zSetKey, hashKey, start, end, targetMember, false)
}

// LuaRedisZSetGetTargetKeyAndStartToEndRankByScore 方法用于处理LuaRedisZSetGetTargetKeyAndStartToEndRankByScore相关逻辑。
func (r *RedisConfig) LuaRedisZSetGetTargetKeyAndStartToEndRankByScore(key string, start, end int64, targetMember string, descending bool) (*LuaRespData, error) {
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

	result, err := redis.NewScript(lua).Run(context.Background(), r.UniversalClient, []string{key}, start, end, targetMember, descendingStr).Result()
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

// LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreDesc 方法用于处理LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreDesc相关逻辑。
func (r *RedisConfig) LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreDesc(key string, start, end int64, targetMember string) (*LuaRespData, error) {
	return r.LuaRedisZSetGetTargetKeyAndStartToEndRankByScore(key, start, end, targetMember, true)
}

// LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreAsc 方法用于处理LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreAsc相关逻辑。
func (r *RedisConfig) LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreAsc(key string, start, end int64, targetMember string) (*LuaRespData, error) {
	return r.LuaRedisZSetGetTargetKeyAndStartToEndRankByScore(key, start, end, targetMember, false)
}

type MemberInfo struct {
	Member    string      `json:"member"`
	Score     json.Number `json:"score"`
	Rank      int64       `json:"rank"`
	Exists    bool        `json:"exists"`
	HashValue string      `json:"hash_value"`
}

// LuaRedisZSetGetMemberScoreAndRankAndGetHashValue 方法用于处理LuaRedisZSetGetMemberScoreAndRankAndGetHashValue相关逻辑。
func (r *RedisConfig) LuaRedisZSetGetMemberScoreAndRankAndGetHashValue(zSetKey, hashKey string, member string, descending bool) (*MemberInfo, error) {
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

	result, err := redis.NewScript(lua).Run(context.Background(), r.UniversalClient, []string{zSetKey, hashKey}, member, descendingStr).Result()
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

// LuaRedisZSetGetMemberScoreAndRankAndGetHashValueDesc 方法用于处理LuaRedisZSetGetMemberScoreAndRankAndGetHashValueDesc相关逻辑。
func (r *RedisConfig) LuaRedisZSetGetMemberScoreAndRankAndGetHashValueDesc(zSetKey, hashKey string, member string) (*MemberInfo, error) {
	return r.LuaRedisZSetGetMemberScoreAndRankAndGetHashValue(zSetKey, hashKey, member, true)
}

// LuaRedisZSetGetMemberScoreAndRankAndGetHashValueAsc 方法用于处理LuaRedisZSetGetMemberScoreAndRankAndGetHashValueAsc相关逻辑。
func (r *RedisConfig) LuaRedisZSetGetMemberScoreAndRankAndGetHashValueAsc(zSetKey, hashKey string, member string) (*MemberInfo, error) {
	return r.LuaRedisZSetGetMemberScoreAndRankAndGetHashValue(zSetKey, hashKey, member, false)
}

// LuaRedisZSetGetMemberScoreAndRank 方法用于处理LuaRedisZSetGetMemberScoreAndRank相关逻辑。
func (r *RedisConfig) LuaRedisZSetGetMemberScoreAndRank(key string, member string, descending bool) (*MemberInfo, error) {
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

	result, err := redis.NewScript(lua).Run(context.Background(), r.UniversalClient, []string{key}, member, descendingStr).Result()
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

// LuaRedisZSetGetMemberScoreAndRankDesc 方法用于处理LuaRedisZSetGetMemberScoreAndRankDesc相关逻辑。
func (r *RedisConfig) LuaRedisZSetGetMemberScoreAndRankDesc(key string, member string) (*MemberInfo, error) {
	return r.LuaRedisZSetGetMemberScoreAndRank(key, member, true)
}

// LuaRedisZSetGetMemberScoreAndRankAsc 方法用于处理LuaRedisZSetGetMemberScoreAndRankAsc相关逻辑。
func (r *RedisConfig) LuaRedisZSetGetMemberScoreAndRankAsc(key string, member string) (*MemberInfo, error) {
	return r.LuaRedisZSetGetMemberScoreAndRank(key, member, false)
}

// LuaRedisZSetGetMultipleMembersScoreAndRankAndHashValues 方法用于处理LuaRedisZSetGetMultipleMembersScoreAndRankAndHashValues相关逻辑。
func (r *RedisConfig) LuaRedisZSetGetMultipleMembersScoreAndRankAndHashValues(zSetKey, hashKey string, members []string, descending bool) ([]*MemberInfo, error) {
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

	result, err := redis.NewScript(lua).Run(context.Background(), r.UniversalClient, []string{zSetKey, hashKey}, args...).Result()
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

// LuaRedisZSetGetMultipleMembersScoreAndRankAndHashValuesDesc 方法用于处理LuaRedisZSetGetMultipleMembersScoreAndRankAndHashValuesDesc相关逻辑。
func (r *RedisConfig) LuaRedisZSetGetMultipleMembersScoreAndRankAndHashValuesDesc(zSetKey, hashKey string, members []string) ([]*MemberInfo, error) {
	return r.LuaRedisZSetGetMultipleMembersScoreAndRankAndHashValues(zSetKey, hashKey, members, true)
}

// LuaRedisZSetGetMultipleMembersScoreAndRankAndHashValuesAsc 方法用于处理LuaRedisZSetGetMultipleMembersScoreAndRankAndHashValuesAsc相关逻辑。
func (r *RedisConfig) LuaRedisZSetGetMultipleMembersScoreAndRankAndHashValuesAsc(zSetKey, hashKey string, members []string) ([]*MemberInfo, error) {
	return r.LuaRedisZSetGetMultipleMembersScoreAndRankAndHashValues(zSetKey, hashKey, members, false)
}

// LuaRedisZSetGetMultipleMembersScoreAndRank 方法用于处理LuaRedisZSetGetMultipleMembersScoreAndRank相关逻辑。
func (r *RedisConfig) LuaRedisZSetGetMultipleMembersScoreAndRank(key string, members []string, descending bool) ([]*MemberInfo, error) {
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

	result, err := redis.NewScript(lua).Run(context.Background(), r.UniversalClient, []string{key}, args...).Result()
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

// LuaRedisZSetGetMultipleMembersScoreAndRankDesc 方法用于处理LuaRedisZSetGetMultipleMembersScoreAndRankDesc相关逻辑。
func (r *RedisConfig) LuaRedisZSetGetMultipleMembersScoreAndRankDesc(key string, members []string) ([]*MemberInfo, error) {
	return r.LuaRedisZSetGetMultipleMembersScoreAndRank(key, members, true)
}

// LuaRedisZSetGetMultipleMembersScoreAndRankAsc 方法用于处理LuaRedisZSetGetMultipleMembersScoreAndRankAsc相关逻辑。
func (r *RedisConfig) LuaRedisZSetGetMultipleMembersScoreAndRankAsc(key string, members []string) ([]*MemberInfo, error) {
	return r.LuaRedisZSetGetMultipleMembersScoreAndRank(key, members, false)
}

// 1. 分布式锁相关

// LuaRedisDistributedLock 方法用于处理LuaRedisDistributedLock相关逻辑。
func (r *RedisConfig) LuaRedisDistributedLock(key, value string, expireSeconds int64) (bool, error) {
	lua := `if redis.call('SET', KEYS[1], ARGV[1], 'NX', 'EX', ARGV[2]) then
				return 1
			else
				return 0
			end`

	result, err := redis.NewScript(lua).Run(context.Background(), r.UniversalClient, []string{key}, value, expireSeconds).Result()
	if err != nil {
		return false, err
	}
	return result.(int64) == 1, nil
}

// LuaRedisDistributedUnlock 方法用于处理LuaRedisDistributedUnlock相关逻辑。
func (r *RedisConfig) LuaRedisDistributedUnlock(key, value string) (bool, error) {
	lua := `if redis.call('GET', KEYS[1]) == ARGV[1] then
				return redis.call('DEL', KEYS[1])
			else
				return 0
			end`

	result, err := redis.NewScript(lua).Run(context.Background(), r.UniversalClient, []string{key}, value).Result()
	if err != nil {
		return false, err
	}
	return result.(int64) == 1, nil
}

// 2. 限流相关

// LuaRedisRateLimit 方法用于处理LuaRedisRateLimit相关逻辑。
func (r *RedisConfig) LuaRedisRateLimit(key string, window, limit int64) (int64, error) {
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

	result, err := redis.NewScript(lua).Run(context.Background(), r.UniversalClient, []string{key}, window, limit).Result()
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

// LuaRedisIncrWithLimit 方法用于处理LuaRedisIncrWithLimit相关逻辑。
func (r *RedisConfig) LuaRedisIncrWithLimit(key string, increment, maxValue, expireSeconds int64) (*CounterResult, error) {
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

	result, err := redis.NewScript(lua).Run(context.Background(), r.UniversalClient, []string{key}, increment, maxValue, expireSeconds).Result()
	if err != nil {
		return nil, err
	}

	var counterResult *CounterResult
	err = json.Unmarshal([]byte(result.(string)), &counterResult)
	return counterResult, err
}

// 4. 队列相关

// LuaRedisQueuePushWithLimit 方法用于处理LuaRedisQueuePushWithLimit相关逻辑。
func (r *RedisConfig) LuaRedisQueuePushWithLimit(key, value string, maxLength int64) (int64, error) {
	lua := `local key = KEYS[1]
			local value = ARGV[1]
			local max_length = tonumber(ARGV[2])
			
			local current_length = redis.call('LLEN', key)
			
			if current_length < max_length then
				return redis.call('LPUSH', key, value)
			else
				return -1
			end`

	result, err := redis.NewScript(lua).Run(context.Background(), r.UniversalClient, []string{key}, value, maxLength).Result()
	if err != nil {
		return -1, err
	}
	return result.(int64), nil
}

// 5. 缓存相关

// LuaRedisSetWithVersion 方法用于处理LuaRedisSetWithVersion相关逻辑。
func (r *RedisConfig) LuaRedisSetWithVersion(key, value string, version, expireSeconds int64) (bool, error) {
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

	result, err := redis.NewScript(lua).Run(context.Background(), r.UniversalClient, []string{key}, value, version, expireSeconds).Result()
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

// LuaRedisDecrStock 方法用于处理LuaRedisDecrStock相关逻辑。
func (r *RedisConfig) LuaRedisDecrStock(key string, quantity int64) (*StockResult, error) {
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

	result, err := redis.NewScript(lua).Run(context.Background(), r.UniversalClient, []string{key}, quantity).Result()
	if err != nil {
		return nil, err
	}

	var stockResult *StockResult
	err = json.Unmarshal([]byte(result.(string)), &stockResult)
	return stockResult, err
}

// 7. HyperLogLog 去重计数

// LuaRedisHLLAddAndCount 方法用于处理LuaRedisHLLAddAndCount相关逻辑。
func (r *RedisConfig) LuaRedisHLLAddAndCount(key string, elements []string) (int64, error) {
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

	result, err := redis.NewScript(lua).Run(context.Background(), r.UniversalClient, []string{key}, args...).Result()
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

// LuaRedisLeaderboardIncr 方法用于处理LuaRedisLeaderboardIncr相关逻辑。
func (r *RedisConfig) LuaRedisLeaderboardIncr(key, member string, increment float64) (*LeaderboardMember, error) {
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

	result, err := redis.NewScript(lua).Run(context.Background(), r.UniversalClient, []string{key}, member, increment).Result()
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

// LuaRedisDelayQueuePop 方法用于处理LuaRedisDelayQueuePop相关逻辑。
func (r *RedisConfig) LuaRedisDelayQueuePop(key string, currentTime int64, limit int64) ([]*DelayedMessage, error) {
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

	result, err := redis.NewScript(lua).Run(context.Background(), r.UniversalClient, []string{key}, currentTime, limit).Result()
	if err != nil {
		return nil, err
	}

	var messages []*DelayedMessage
	err = json.Unmarshal([]byte(result.(string)), &messages)
	return messages, err
}

// 10. 布隆过滤器模拟 (使用多个 Hash)

// LuaRedisBloomAdd 方法用于处理LuaRedisBloomAdd相关逻辑。
func (r *RedisConfig) LuaRedisBloomAdd(key, element string) error {
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

	_, err := redis.NewScript(lua).Run(context.Background(), r.UniversalClient, []string{key}, element).Result()
	return err
}

// LuaRedisBloomExists 方法用于处理LuaRedisBloomExists相关逻辑。
func (r *RedisConfig) LuaRedisBloomExists(key, element string) (bool, error) {
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

	result, err := redis.NewScript(lua).Run(context.Background(), r.UniversalClient, []string{key}, element).Result()
	if err != nil {
		return false, err
	}
	return result.(int64) == 1, nil
}

type luaRedisIDConfig struct {
	key         string // 健名
	startNumber int64  // 起始值
	iNCRValue   int64  // 每次自增的值
}

type WithLuaRedisIDConfigOption func(*luaRedisIDConfig)

// WithLuaRedisIDConfigKeyName 函数用于处理WithLuaRedisIDConfigKeyName相关逻辑。
func WithLuaRedisIDConfigKeyName(key string) WithLuaRedisIDConfigOption {
	return func(config *luaRedisIDConfig) {
		config.key = key
	}
}

// WithLuaRedisIDConfigStartNumber 函数用于处理WithLuaRedisIDConfigStartNumber相关逻辑。
func WithLuaRedisIDConfigStartNumber(startNumber int64) WithLuaRedisIDConfigOption {
	return func(config *luaRedisIDConfig) {
		config.startNumber = startNumber
	}
}

// WithLuaRedisIDConfigINCRValue 函数用于处理WithLuaRedisIDConfigINCRValue相关逻辑。
func WithLuaRedisIDConfigINCRValue(INCRValue int64) WithLuaRedisIDConfigOption {
	return func(config *luaRedisIDConfig) {
		config.iNCRValue = INCRValue
	}
}

// LuaRedisID 方法用于处理LuaRedisID相关逻辑。
func (r *RedisConfig) LuaRedisID(opts ...WithLuaRedisIDConfigOption) (int64, error) {
	idConfig := &luaRedisIDConfig{
		key:         "global-id",
		startNumber: 10000,
		iNCRValue:   1,
	}
	for i := range opts {
		opts[i](idConfig)
	}
	var script = `
		local current = tonumber(redis.call('GET', KEYS[1])) or 0
		local target = tonumber(ARGV[1])
		local iNCRValue = tonumber(ARGV[2])
		if current < target then
			redis.call('SET', KEYS[1], target)
		end
		return redis.call('incrby', KEYS[1], iNCRValue)
	`
	result, err := redis.NewScript(script).Run(context.Background(), r.UniversalClient, []string{idConfig.key}, idConfig.startNumber, idConfig.iNCRValue).Result()
	if err != nil {
		return 0, err
	}
	return cast.ToInt64(result), nil
}
