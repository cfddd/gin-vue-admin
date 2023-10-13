package service

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	plugGlobal "github.com/flipped-aurora/gin-vue-admin/server/plugin/register/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/register/model"
	userService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/mojocn/base64Captcha"
)

type RegisterService struct{}

func (e *RegisterService) PlugService(req model.Request) (res *system.SysUser, err error) {
	if err := utils.Verify(req, utils.LoginVerify); err != nil {
		return res, err
	}
	var (
		store = base64Captcha.DefaultMemStore
		user  system.SysUser
		us    *userService.UserService
	)
	if !store.Verify(req.CaptchaId, req.Captcha, true) {
		return res, errors.New("验证码错误")
	}

	u := &system.SysUser{Username: req.Username, Password: req.Password, NickName: req.Username, Phone: req.Phone, Email: req.Email, QQ: req.QQ}
	err = global.GVA_DB.Where("username = ?", u.Username).Preload("Authorities").Preload("Authority").First(&user).Error
	if err == nil {
		return res, errors.New("用户名已注册")
	}
	if user.Username != "" && user.Password != "" {
		return res, errors.New("用户名已注册")
	}

	var authorities []system.SysAuthority

	authorities = append(authorities, system.SysAuthority{
		AuthorityId: plugGlobal.GlobalConfig.AuthorityId,
	})
	u.Authorities = authorities
	u.AuthorityId = plugGlobal.GlobalConfig.AuthorityId

	if rest, err := us.Register(*u); err != nil {
		return &rest, errors.New("注册失败!")
	}
	if res, err = us.Login(u); err != nil {
		return res, errors.New("登陆失败!")
	}
	return res, nil
	// 前面的代码 拿不到正确的 user，所以需要再次查询一次

}
