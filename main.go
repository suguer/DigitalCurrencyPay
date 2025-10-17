package main

import (
	"DigitalCurrency/internal/config"
	"DigitalCurrency/internal/consumer"
	"DigitalCurrency/internal/crontab"
	"DigitalCurrency/internal/logger"
	"DigitalCurrency/internal/model"
	"DigitalCurrency/internal/model/dao"
	"DigitalCurrency/internal/router"
	"DigitalCurrency/internal/runner"
	"DigitalCurrency/internal/service/user"
	"DigitalCurrency/internal/util"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

var (
	configFile string = "/etc/config.yaml"
	ctx        context.Context
)

func init() {
	ctx = context.Background()
}

func main() {
	command := ""
	if len(os.Args) < 2 {
		command = "help"
	} else {
		command = os.Args[1]
	}
	switch command {
	case "1":
		run()
	case "2":
		account()
	case "3":
		password()
	case "4":
		secret()
	case "5":
		defaultMessage()
	case "help":
		fmt.Println(
			`
DCpay 管理命令:
1 - 运行程序
2 - 设置用户名
3 - 设置密码
4 - 设置Secret
5 - 查看默认信息
`)
	default:
		fmt.Printf("未知命令: %v\n", command)
		return
	}
	return

}

func init() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取工作目录失败: %v", err)
	}
	_, err = config.Load(filepath.Join(wd, configFile))
	logger.InitLogger(config.Conf.Storage)
	dao.InitDatabase(config.Conf.Database)
	dao.InitRedis(&config.Conf.Redis)
	dao.InitCache(ctx)
	model.Migration()
}
func run() {

	consumer.InitConsumers(ctx, config.Conf.Queue.Driver)
	runner.InitRunner(ctx)
	loc, _ := time.LoadLocation("Asia/Shanghai")
	time.Local = loc
	c := cron.New(cron.WithLocation(loc))
	crontab.InitCrontab(c, ctx)
	go c.Run()
	r := gin.Default()
	router.Register(r)
	r.Run(":8080")
}

func account() {
	var input string
	if len(os.Args) > 2 {
		input = os.Args[2]
	} else {
		fmt.Print("请输入用户名: ")
		fmt.Scanln(&input)
	}
	err := user.Update(1, map[string]any{
		"username": input,
	})
	if err != nil {
		fmt.Printf("更新用户名失败: %v\n", err)
		return
	}
}

func password() {
	var input string
	if len(os.Args) > 2 {
		input = os.Args[2]
	} else {
		fmt.Print("请输入密码: ")
		fmt.Scanln(&input)
	}
	err := user.Update(1, map[string]any{
		"password": input,
	})
	if err != nil {
		fmt.Printf("更新密码失败: %v\n", err)
		return
	}
}

func secret() {
	var input string
	if len(os.Args) > 2 {
		input = os.Args[2]
	} else {
		fmt.Print("请输入Secret: ")
		fmt.Scanln(&input)
	}
	err := user.Update(1, map[string]any{
		"secret": input,
	})
	if err != nil {
		fmt.Printf("更新Secret失败: %v\n", err)
		return
	}
}
func defaultMessage() {
	path := filepath.Join(config.Conf.Storage.Path, "defaultMessage.txt")
	defaultMessage, err := util.ReadFile(path)
	if err != nil {
		fmt.Printf("读取默认信息失败: %v\n", err)
		return
	}
	fmt.Printf("默认信息\n%v\n", defaultMessage)
}
