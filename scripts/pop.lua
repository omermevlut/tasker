local job = redis.call('lpop', KEYS[1])

if(job ~= false) then
    local popped = cjson.decode(job)
    popped['executed_at'] = KEYS[3]
    popped = cjson.encode(popped)

    redis.call('zadd', KEYS[2], ARGV[1], popped)
end

return job
