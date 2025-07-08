package middleware

const RateLimit = `
local now = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local maxCount = tonumber(ARGV[3])
local key = KEYS[1]

redis.call("ZREMRANGEBYSCORE",key,now-window,now)
local count = redis.call("ZCARD",key)

if tonumber(count) >= maxCount then
return 0
else
redis.call("ZADD",key,tonumber(now),tonumber(now))
return 1
end
`

const TokenBucket = `
    local now = tonumber(ARGV[1])
    local maxToken = tonumber(ARGV[2])
    local interval = tonumber(ARGV[3])
    local tokenKey = KEYS[1]
    local lastTimeKey = KEYS[2]
    local tokenCount = redis.call("GET",tokenKey)
    local lastTime = redis.call("GET",lastTimeKey)

    if not tokenCount then
		redis.call("SET",tokenKey,maxToken-1,"EX",300)
		redis.call("SET",lastTimeKey,now,"EX", 300)
		return 1
    end

    if lastTime then
		lastTime = tonumber(lastTime)
	else
		lastTime = now
    end

    local needAdd = math.floor((now-lastTime)/interval)
    if needAdd > 0 then
		local finalAdd = math.min(needAdd+tokenCount,maxToken)
		redis.call("SET",tokenKey,finalAdd,"EX",300)
		tokenCount = finalAdd
    end

    if tonumber(tokenCount) <= 0 then
        return 0
    else
		redis.call("SET",tokenKey,tokenCount-1,"EX",300)
		redis.call("SET",lastTimeKey,now,"EX", 300)
        return 1
    end

`
