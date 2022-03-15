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
