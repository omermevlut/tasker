package scripts

// GetMoveExpiredScript ...
func GetMoveExpiredScript() string {
	return `local val = redis.call('zrangebyscore', KEYS[1], '-inf', ARGV[1])

if(next(val) ~= nil) then
    redis.call('zremrangebyrank', KEYS[1], 0, #val - 1)

    for i = 1, #val, 100 do
        redis.call('rpush', KEYS[2], unpack(val, i, math.min(i+99, #val)))
    end
end

return val
`
}

// GetPopFromActiveScript ..
func GetPopFromActiveScript() string {
	return `local job = redis.call('lpop', KEYS[1])

if(job ~= false) then
    local popped = cjson.decode(job)
    popped['executed_at'] = KEYS[3]
    popped = cjson.encode(popped)

    redis.call('zadd', KEYS[2], ARGV[1], popped)
end

return job`
}
