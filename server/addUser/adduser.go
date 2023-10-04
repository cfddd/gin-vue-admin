package main

import (
	"bufio"
	"fmt"
	"gorm.io/driver/mysql"
	"os"
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"

	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gofrs/uuid/v5"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	r  systemReq.Register
	db *gorm.DB
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Register
//@description: 用户注册
//@param: u model.SysUser
//@return: userInter system.SysUser, err error

func Register(u system.SysUser) (err error) {
	var user system.SysUser
	if !errors.Is(db.Where("username = ?", u.Username).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return err
	}
	// 否则 附加uuid 密码hash加密 注册
	u.Password = utils.BcryptHash(u.Password)
	u.UUID = uuid.Must(uuid.NewV4())
	err = db.Create(&u).Error
	if err != nil {
		println(err)
	}
	return err
}

func cin(scanner *bufio.Scanner) {
	// 从文件中读取数据
	if scanner.Scan() {
		r.Username = scanner.Text()
	}

	if scanner.Scan() {
		r.NickName = scanner.Text()
	}

	r.HeaderImg = "https://qmplusimg.henrongyi.top/1576554439myAvatar.png"

	r.AuthorityId = 1

	r.AuthorityIds = append(r.AuthorityIds, 1)

	if scanner.Scan() {
		r.Phone = scanner.Text()
	}

	if scanner.Scan() {
		r.Email = scanner.Text()
	}

	if scanner.Scan() {
		r.QQ = scanner.Text()
	}
}

func caseT(scanner *bufio.Scanner) {
	cin(scanner)

	var authorities []system.SysAuthority
	for _, v := range r.AuthorityIds {
		authorities = append(authorities, system.SysAuthority{
			AuthorityId: v,
		})
	}

	user := &system.SysUser{
		Username:       r.Username,
		NickName:       r.NickName,
		Password:       r.Password,
		HeaderImg:      r.HeaderImg,
		AuthorityId:    r.AuthorityId,
		Authorities:    authorities,
		Enable:         r.Enable,
		Phone:          r.Phone,
		Email:          r.Email,
		QQ:             r.QQ,
		DACountInMouth: 0,
	}
	err := Register(*user)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	dsn := "root:123456@tcp(177.7.0.13:3306)/gva?charset=utf8mb4&parseTime=True&loc=Local"
	// 177.7.0.13:3306是再docker网络中的mysql ip+端口
	// localhost:3306是本机mysql默认ip+端口，按需修改
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	fmt.Println("请输入插入用户的数量")
	file, err := os.Open("usersList")
	if err != nil {
		panic("Failed to open usersList file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// 设置分隔符为空格
	scanner.Split(bufio.ScanWords)

	scanner.Scan()
	T, err := strconv.Atoi(scanner.Text())
	for T > 0 {
		T--
		caseT(scanner)
	}
}
/*
批量导入用户，userList文件格式如下
T行格式如下
Username NickName Phone Email QQ
---
运行命令
./go_build_adduser_go_linux
*/