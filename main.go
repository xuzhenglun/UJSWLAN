package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/codegangsta/cli"
	"github.com/xuzhenglun/UJSWlan/core"
)

type Config struct {
	Username string
	Password string
}

func main() {
	app := cli.NewApp()
	app.Name = "UJS WAN LOGIN"
	app.Usage = "模拟WEB页面登录注销江大网络"
	app.Commands = []cli.Command{
		{
			Name:    "login",
			Aliases: []string{"i"},
			Usage:   "Exit after Login",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "username,u",
					Usage: "Student Id of UJS",
				},
				cli.StringFlag{
					Name:  "passwd,p",
					Usage: "Password of UJS",
				},
			},
			Action: func(c *cli.Context) {
				conn := core.NewConnect()
				conn.Setting(c.String("username"), c.String("passwd"))
				err := conn.Login()
				if err == core.LOGIN_SUCC {
					log.Println("登录成功")
				} else {
					log.Println("登录失败")
				}
			},
		}, {
			Name:    "logout",
			Aliases: []string{"o"},
			Usage:   "Logout",

			Action: func(c *cli.Context) {
				conn := core.NewConnect()
				if err := conn.Logout(); err == core.LOGOUT_SUCC {
					log.Println("登出成功")
				} else {
					log.Println("登出失败")
				}
			},
		},
		{
			Name:    "status",
			Aliases: []string{"s"},
			Usage:   "Show Connection Status",

			Action: func(c *cli.Context) {
				if err := core.Status(); err == nil {
					log.Println("Connection Good")
				} else {
					log.Println("Not Good")
				}
			},
		},
	}
	app.Action = func(c *cli.Context) {

		loadedConfig, err := ioutil.ReadFile("config.json")
		if err != nil {
			log.Println("JSON配置文件不存在，程序将创建实例文件5秒后退出...\n\t帮助输入-h或--help以查阅")

			if err := ioutil.WriteFile("config.json", []byte("{\n\t\"username\":\"yourId\",\n\t\"password\":\"yourpassword\"\n}"), 660); err != nil {
				log.Panic(err)
			}

			time.Sleep(5 * time.Second)
			os.Exit(0)
		}

		var config Config

		err = json.Unmarshal(loadedConfig, &config)
		if err != nil {
			log.Println("JSON配置文件有错误，程序将退出...\n\t帮助输入-h或--help以查阅")
			os.Exit((0))
		}

		conn := core.NewConnect()
		conn.Setting(config.Username, config.Password)

		go func() {
			for {
				if err = conn.Login(); err != core.LOGIN_SUCC {
					log.Println(err)
					log.Println("Try after 5 seconds")
					time.Sleep(5 * time.Second)
				} else {
					log.Println("恭喜你，登录成功啦\n\t如需注销请按下Crtl+C")
					break
				}
			}
		}()
		signalChan := make(chan os.Signal)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan
		log.Println("EXCTING...")
		if err = conn.Logout(); err != core.LOGOUT_SUCC {
			log.Println("LOGOUT FAILED")
			log.Println(err)
		} else {
			log.Println("退出成功")
		}
		time.Sleep(time.Second)
		os.Exit(0)
	}
	app.Run(os.Args)
}
