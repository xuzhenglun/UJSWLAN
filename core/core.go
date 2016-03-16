package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Connect struct {
	Username string
	Password string
	RetRe    *regexp.Regexp
	ErrRe    *regexp.Regexp
}

type LoginErr struct {
	What string
}

func (this LoginErr) Error() string {
	return fmt.Sprintln("%v: %v", time.Now(), this.What)
}

var IP_NOT_ALLOW LoginErr = LoginErr{What: "本IP不允许Web方式登录"}
var ACCOUNT_NOT_ALLOW LoginErr = LoginErr{What: "本账号不允许Web方式登录"}
var PASSWD_NOT_ALLOW LoginErr = LoginErr{What: "本账号不允许修改密码"}
var PASSWD_ERR LoginErr = LoginErr{What: "帐号或密码有误，请重新输入"}
var HAVE_LOGIN LoginErr = LoginErr{What: "该账号正在别的的机器上使用，先注销它再说吧"}
var LIMIT_AREA_ERR LoginErr = LoginErr{What: "本账号只能在指定地址使用"}
var NOT_ENOUGH_MONEY LoginErr = LoginErr{What: "本账号费用超支或时长流量超过限制"}
var FORZE_ERR LoginErr = LoginErr{What: "本账号暂停使用"}
var FULL_ERR LoginErr = LoginErr{What: "System buffer full"}
var INUSE_ERR LoginErr = LoginErr{What: "本账号正在使用,不能修改"}
var NEW_NOT_MATCH LoginErr = LoginErr{What: "新密码与确认新密码不匹配,不能修改"}
var CHANGE_SUCC LoginErr = LoginErr{What: "密码修改成功"}
var LOGOUT_SUCC LoginErr = LoginErr{What: "注销成功"}
var LOGIN_SUCC LoginErr = LoginErr{What: "注销成功"}
var UNEXPECT_ERR LoginErr = LoginErr{What: "Unexpect Err"}

func RetMessage(msg, err string) error {
	m, e := strconv.Atoi(msg)
	if e != nil {
		log.Println(msg, err)
		return UNEXPECT_ERR
	}
	switch m {
	case 1:
		switch err {
		case "error0":
			return IP_NOT_ALLOW
		case "error1":
			return ACCOUNT_NOT_ALLOW
		case "error2":
			return PASSWD_NOT_ALLOW
		default:
			return PASSWD_ERR
		}
	case 2:
		return HAVE_LOGIN
	case 3:
		return LIMIT_AREA_ERR
	case 4:
		return NOT_ENOUGH_MONEY
	case 5:
		return FORZE_ERR
	case 6:
		return FULL_ERR
	case 8:
		return INUSE_ERR
	case 9:
		return NEW_NOT_MATCH
	case 10:
		return CHANGE_SUCC
	case 11:
		return LIMIT_AREA_ERR
	case 14:
		return LOGOUT_SUCC
	case 15:
		return LOGOUT_SUCC
	default:
		return UNEXPECT_ERR

	}
}

func NewConnect() *Connect {
	var c Connect
	c.RetRe, _ = regexp.Compile(`Msg=[\d]{2}`)
	c.ErrRe, _ = regexp.Compile(`msga='[\d]*`)
	return &c
}

func (this *Connect) Setting(Username, Password string) {
	this.Username = Username
	this.Password = Password
}

func (this Connect) Login() error {

	url := "http://192.168.100.83/0.htm"

	payload := strings.NewReader("DDDDD=" + this.Username + "&upass=" + this.Password + "&0MKKey=")

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		log.Println("MEW POST ERR")
		return UNEXPECT_ERR
	}

	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("postman-token", "943f2471-06ca-15a2-ad4d-cf71d104b05c")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println(err)
		log.Println("POST ERR")
		return UNEXPECT_ERR
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	succre, _ := regexp.Compile("You have successfully logged into our system.")
	if succre.Match(body) {
		return LOGIN_SUCC
	}

	msg := this.RetRe.FindString(string(body))
	e := this.ErrRe.FindString(string(body))
	return RetMessage(string([]byte(msg)[4:]), string([]byte(e)[6:]))
}

func (this Connect) Logout() error {

	url := "http://192.168.100.83/F.htm"

	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return UNEXPECT_ERR
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	msg := this.RetRe.FindString(string(body))
	e := this.ErrRe.FindString(string(body))

	return RetMessage(string([]byte(msg)[4:]), string([]byte(e)[6:]))
}

func Status() error {

	url := "http://baidu.com"

	_, err := http.Get(url)
	return err
}
