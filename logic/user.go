package logic

import (
	"NothingBlog/dao/mysql"
	"NothingBlog/dao/redis"
	"NothingBlog/models"
	"NothingBlog/package/jwt"
	"NothingBlog/package/snowflake"
	"NothingBlog/package/verifycode"
	"crypto/md5"
	"encoding/hex"
	"errors"

	"go.uber.org/zap"
)

const secret = "xl_ngblog" // md5加密时使用的密钥

var ErrUserPassword = errors.New("用户密码不正确")
var ErrUserNotExisted = errors.New("用户不存在")
var ErrVerifyCode = errors.New("验证码错误")

// 用户注册的逻辑层处理函数 - 返回值error
func UserSignup(u *models.SignUpParams) (err error) {
	// 检查用户注册是否合法 - 是否存在
	nu := &models.User{
		UserName: u.UserName,
		Email:    u.Email,
	}
	err = mysql.QueryUserByName(nu)

	if err != nil {
		zap.L().Debug("用户已经存在", zap.Error(err))
		return
	}

	// 检查验证码是否有效
	if ok := verifycode.CheckVerifyCode(u.VerifyCode.Id, u.VerifyCode.Code); !ok {
		zap.L().Debug("验证码输入错误")
		return ErrVerifyCode
	}

	// 使用雪花算法生成用户id
	nu.UserId = snowflake.GetNextId().Int64()

	// 用户密码加密
	nu.Password = EncryptContent(u.Password)

	// 持久化存储 - mysql
	return mysql.InsertUser(nu)
}

// 用户登录逻辑
func UserLogin(u *models.LoginParams) (token string, err error) {
	usr := &models.User{
		UserName: u.UserName,
	}
	// 从数据库读取对应的用户 - 返回nil表示用户不存在
	err = mysql.QueryUserByName(usr)
	if err == nil {
		return "", ErrUserNotExisted
	}

	if err == mysql.ErrQueryFailed {
		return "", err
	}

	// 检查密码是否匹配
	if EncryptContent(u.Password) != usr.Password {
		return "", ErrUserPassword
	}

	// 检查验证码是否匹配
	ok := verifycode.CheckVerifyCode(u.VerifyCode.Id, u.VerifyCode.Code)
	if !ok {
		return "", ErrVerifyCode
	}

	// 生成token - 这里生成的是access token，jwt token存在的问题是服务端在签发后无法控制token的有效性
	// 只能因为到期而失效，所以为了尽可能缓解这个问题，一般将这个token的有效期设置的很短，然后再引入一个不含任何自定义
	// 信息的refresh token，这个token专用于在access token失效后进行重置使用 - 本代码没有实现
	token, err = jwt.GenerateJwtToken(usr.UserId)
	if err != nil {
		return "", err
	}

	// 将token存储在redis中，实现单一用户单点登录
	redis.InsertLoginInfo(usr.UserId, token)

	// 验证成功
	return token, nil
}

// 使用md5算法对内容加密
func EncryptContent(str string) string {
	ecy := md5.New()
	ecy.Write([]byte(secret))
	return hex.EncodeToString(ecy.Sum([]byte(str)))
}
