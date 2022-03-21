package pkg

//GetStrSetStr string 操作，获取key是否存在，不存在设置key及过期时间
//1代表成功，2代表失败
//没有对设置key的过程结果做处理
const GetStrSetStr = `
	local redisKeys = redis.call('GET',KEYS[1]);
	local res = {"0",'ERR'};
	if  not redisKeys  
	then 
		redis.call('SET',KEYS[1],ARGV[1]);
		redis.call('EXPIRE', KEYS[1], ARGV[2]);
		res[1]= "1";res[2]=ARGV[1];
		return res;
	else 
		res[1]= "2";res[2]=redisKeys;
		return res;
	end
	return res;
`
//IncrNum num 增加操作
//KEYS[1] 是key
//ARGV[1] 是总数
//ARGV[2] 要增加的数值
//ARGV[3] 过期时间，默认是-1
//1 是成功 返回增加后的值
//2 是失败 返回现在的值或者错误
const IncrNum = `
	local key = redis.call('GET',KEYS[1])
	local res = {'2','0'}
	local expire = tonumber(ARGV[3]) 
	local num1 = tonumber(ARGV[1]) - tonumber(ARGV[2])
	local num = -1 
	if key then  num = tonumber(key) end
	
	if (key and num > -1 and  num <= num1)
	then
		redis.call('incrBy',KEYS[1],tonumber(ARGV[2])) 
		res[1]= '1' res[2] = tostring(num + tonumber(ARGV[2]))
	else
		res[2] = key
	end

	if not key
	then 
		redis.call('SET',KEYS[1],ARGV[2]) 
		res[1]= '1' res[2] = ARGV[2]
	end
	
	if expire > 0 then  redis.call('EXPIRE', KEYS[1], expire) end
	return res;
`

//HSetANdExpire hash set 操作，设置成员并设置过期时间
//KEYS[1] 设置的hash key
//ARGV[1] 要设置的成员的个数
//ARGV[2] 过期时间
//1 key 2 1000 field1 value1 field2 value2
const HSetANdExpire = `
	local fieldIndex=3
	local valueIndex=4
	local key=KEYS[1]
	local fieldCount=ARGV[1]
	local expired=ARGV[2]
	for i=1,fieldCount,1 do
	  redis.pcall('HSET',key,ARGV[fieldIndex],ARGV[valueIndex])
	  fieldIndex=fieldIndex+2
	  valueIndex=valueIndex+2
	end
	redis.pcall('EXPIRE',key,expired)
`

//SetExpireByTTl 设置key的expire，如果存在，则设置为剩余时间加上当前时间
//KEYS[1] key
//ARGV[1] 要设置过着延长的时间,如果存在，且当前是-1，则设置为传入的
//ARGV[2] value
const SetExpireByTTl = `
	local redisKeys = redis.call('GET',KEYS[1]); 
	local ttl = 0; 
	if not redisKeys 
	then 
		redis.call('SET',KEYS[1],ARGV[1]);
		redis.call('EXPIRE', KEYS[1], ARGV[2]);
	else
		ttl = redis.call('ttl',KEYS[1]); 
		if (ttl > 0  or ttl == -1)
		then 
			ttl = ttl + ARGV[2];
		else 
			ttl = ARGV[2]  
		end
		redis.call('EXPIRE', KEYS[1], ttl);
	end
`
