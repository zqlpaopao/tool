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
