// auth_phone_test
package auth

import (
	"WolaiWebservice/config"
	authControllers "WolaiWebservice/controllers/auth"
	"WolaiWebservice/models"
	myRedis "WolaiWebservice/redis"
	authService "WolaiWebservice/service/auth"
	"WolaiWebservice/utils/encrypt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/redis.v3"
)

const (
	TEST_PHONE = "18570305176"
	TEST_PWD   = "12345678"
)

var (
	redisClient *redis.Client

	RedisFailErr error
)

func init() {
	models.Initialize()
	myRedis.Initialize()
	redisClient = redis.NewClient(
		&redis.Options{
			Addr:     config.Env.Redis.Host + config.Env.Redis.Port,
			Password: config.Env.Redis.Password,
			DB:       config.Env.Redis.Db,
			PoolSize: config.Env.Redis.PoolSize,
		})
	_, err := redisClient.Ping().Result()
	RedisFailErr = err
}

func TestAuthPhoneRegister(t *testing.T) {
	err := authService.SendSMSCode(TEST_PHONE, myRedis.SC_REGISTER_RAND_CODE)
	if err != nil {
		t.Fatalf(err.Error())
	}
	randCode, err := redisClient.HGet(myRedis.SC_REGISTER_RAND_CODE+TEST_PHONE, "randCode").Result()
	if err != nil {
		t.Fatalf(err.Error())
	}
	status, err, authInfo := authControllers.AuthPhoneRegister(TEST_PHONE, randCode, TEST_PWD)
	if status == 1001 {
		t.Fatalf(err.Error())
	} else if status == 2 {
		t.Fatalf(err.Error())
	}
	user, err := models.ReadUser(authInfo.Id)
	if err != nil {
		t.Fatalf("用户不存在")
	}
	pwd := encrypt.EncryptPassword(TEST_PWD, *user.Salt)
	if pwd != *user.Password {
		t.Errorf("密码不匹配")
	}
}

func TestAuthPhonePasswordLogin(t *testing.T) {
	status, err, authInfo := authControllers.AuthPhonePasswordLogin(TEST_PHONE, TEST_PWD)
	if status != 0 {
		t.Fatalf(err.Error())
	}
	if *authInfo.Phone != TEST_PHONE {
		t.Errorf("登陆信息有误")
	}
}

func TestAuthPhoneRandCodeLogin(t *testing.T) {
	err := authService.SendSMSCode(TEST_PHONE, myRedis.SC_LOGIN_RAND_CODE)
	if err != nil {
		t.Fatalf(err.Error())
	}
	randCode, err := redisClient.HGet(myRedis.SC_LOGIN_RAND_CODE+TEST_PHONE, "randCode").Result()
	if err != nil {
		t.Fatalf(err.Error())
	}
	status, err, authInfo := authControllers.AuthPhoneRandCodeLogin(TEST_PHONE, randCode, true)
	if status != 0 {
		t.Fatalf(err.Error())
	}
	if *authInfo.Phone != TEST_PHONE {
		t.Errorf("登陆信息有误")
	}
}

func TestForgotPassword(t *testing.T) {

}
