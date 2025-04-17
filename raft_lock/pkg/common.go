package pkg

import "time"

const (
	readTimeout      = 5 * time.Second
	writeTimeout     = 5 * time.Second
	idleTimeout      = 60 * time.Second
	poolTimeOut      = 60 * time.Second
	poolSize         = 20
	DefaultSameSlot  = "{raft}:"
	groupDefaultName = "default"
	lockNum          = 1
	nodeNum          = 1

	RDBLock    = "Lock"
	RDBRenewal = "Renewal"
	RDBUnLock  = "UnLock"

	Success = 1
	Failed  = 2

	// enqueueCmd enqueues a given task message.
	//
	// Input:
	// KEYS[1] -> group_name
	// KEYS[2] -> lock_name
	// --
	// ARGV[1] -> lock expire time
	// ARGV[2] -> lock len
	// ARGV[3] -> time now

	// LockCmd Output:
	// Returns 1 if successfully enqueued
	// Returns 0
	LockCmd = `
		local n = tonumber(redis.call("HLEN", KEYS[1])) or 0
		if tonumber(n) < tonumber(ARGV[2]) then
			return redis.call("HSET", KEYS[1], KEYS[2], ARGV[1])
		else
			local all_items = redis.call("HGETALL", KEYS[1])
			local deleted = 0
			for i = 1, #all_items, 2 do
				local field = all_items[i]
				local value = tonumber(all_items[i+1])
				if value and value < tonumber(ARGV[3]) then
					redis.call("HDEL", KEYS[1], field)
					deleted = deleted + 1
				end
			end
			if (n - deleted) < tonumber(ARGV[2]) then
				return redis.call("HSET", KEYS[1], KEYS[2], ARGV[1])
			else
				return 0
			end
		end
	`

	// RenewalCmd Returns  enqueueCmd enqueues a given task message.
	// Input:
	// KEYS[1] -> group_name
	// KEYS[2] -> lock_name
	// --
	// ARGV[1] -> lock expire time
	// Output:
	// Returns 1 if successfully enqueued
	// Returns 0
	RenewalCmd = `
		local val =  redis.call("HGET", KEYS[1], KEYS[2])
		if val == false or val == nil  then
			return 0
		end
		redis.call("HSET", KEYS[1], KEYS[2], ARGV[1])
		return 1
	`
	// -- 打印传入的键和参数，用于调试
	//		redis.log(redis.LOG_NOTICE, "KEYS: " .. table.concat(KEYS, ", "))
	//		redis.log(redis.LOG_NOTICE, "ARGV: " .. table.concat(ARGV, ", "))
	//
	//		-- 获取哈希表中指定字段的值
	//		local val =  redis.call("HGET", KEYS[1], KEYS[2])
	//
	//		-- 打印获取到的值，用于调试
	//		-- redis.log(redis.LOG_NOTICE, "Value from HGET: " .. val)
	//		redis.log(redis.LOG_NOTICE, "Value from HGET: " .. tostring(val))
	//
	//		-- 判断值是否为 nil
	//		if val == false or val == nil  then
	//			-- 值为 nil，返回 0
	//			return 0
	//		end
	//
	//		-- 值不为 nil，更新哈希表中指定字段的值
	//		redis.call("HSET", KEYS[1], KEYS[2], ARGV[1])
	//
	//		-- 更新成功，返回 1
	//		return 1

	// DelMemberCmd enqueueCmd enqueues a given task message.
	//
	// Input:
	// KEYS[1] -> group_name
	// KEYS[2] -> lock_name
	// --
	// Output:
	// Returns 1 if successfully enqueued
	// Returns 0
	DelMemberCmd = `
		redis.call("HDEL", KEYS[1],KEYS[2])
   		return 1
	`
)
