// ratelimit
package redis

const (
	rateLimitScript = `
local times = redis.call('incr',KEYS[1])
if times == 1 then
	redis.call('expire',KEYS[1],ARGV[1])
end
if times > tonumber(ARGV[2]) then
	return 0
end
return 1`
)

func RateLimit(key string) int64 {
	result := redisClient.Eval(rateLimitScript, []string{key}, []string{"1", "1"})
	rt, err := result.Result()
	if err != nil {
		return 0
	}
	return rt.(int64)
}
