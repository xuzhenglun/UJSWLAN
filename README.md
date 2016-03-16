# UJSWLAN
无聊分析了下江大校园网的WEB登录协议，简陋到爆炸=-=

---
##HTTP协议
### 登录协议：
```Http
POST /0.htm HTTP/1.1
Host: 192.168.100.83
Cache-Control: no-cache
Content-Type: application/x-www-form-urlencoded

DDDDD={此处是用户名}&upass={此处是密码}&0MKK
```

###登出协议：
```Http
GET /F.htm HTTP/1.1
Host: 192.168.100.83
Cache-Control: no-cache
```

## cURL：


###登入：
```cUEL
curl --request POST \
  --url http://192.168.100.83/0.htm \
  --header 'cache-control: no-cache' \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data 'DDDDD={此处是用户名}&upass={此处是密码}&0MKKey='
```

###登出：
```cURL
curl --request GET \
  --url http://192.168.100.83/F.htm \
  --header 'cache-control: no-cache'
  ```
  
---
##薄弱点：
1. WIFI无加密，验证没有HTTPS，即账户和密码通过明文信息向全世界公开广播
2. 登录验证貌似是通过MAC绑定，监控TTL保证不能私接路由器（没有仔细验证）。
3. 注销断网没有任何验证，没有Cookie，仅仅验证MAC。
综合2、3，可以伪造MAC让任何人断网，后通过无线抓包再重连的时候窃取任何人的账户密码。*关键这个过程毫无技术难度！！！*

##槽点：
- 登录状态返回的状态判断码简直醉人，上一段以供瞻仰：
```Javascript
 <SCRIPT language=javascript>
            <!--     		
Msg=14;time='542       ';flow='1825385   ';fsele=1;fee='0         ';xsele=0;xip='000.000.000.000.';mac='00-00-00-00-00-00';va=00;vb=00;vc=00;vd=0000;ve=0000;vf=0000;
ipm="c0a86453";ss1="000d482e544a";ss2="0000";ss3="0a030b3c";ss4="000000000000";msga='';/* can not modify !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!*/
pp="<p></p>";
flow0=flow%1024;flow1=flow-flow0;flow0=flow0*1000;flow0=flow0-flow0%1024;fee1=fee-fee%100;
flow3='.';
if(flow0/1024<10)flow3='.00';
else{if(flow0/1024<100)flow3='.0';}
UT="已使用时间 Used time : "+time +" Min"+pp;
UF="已使用流量 Used flux : "+flow1/1024+flow3+flow0/1024+" MByte"+pp;
if(fsele==1)UM="本账号余额  Balance : "+"RMB"+fee1/10000;
else UM="";
function DispTFM(){
	switch(Msg){
		case 0:
		case 1:
			if((Msg==1)&&(msga!="")){
			switch(msga){
			case 'error0':
				document.write("本IP不允许Web方式登录<br>The IP does not allow Web-log");
				break;
			case 'error1':
				document.write("本账号不允许Web方式登录<br>The account does not allow Web-log");
				break;
			case 'error2':
				document.write("本账号不允许修改密码<br>This account does not allow change password");
				break;				
			default:
				document.write("帐号或密码有误，请重新输入");
				break;}
			}
			else document.write("账号或密码不对，请重新输入<br>Ivalid account or password, please login again");
		break;      		
	case 2:
		document.write("该账号正在IP为："+xip+"的机器上使用，<br><br>请点击<a href='a11.htm'>继续</a>断开它的连接并重新输入用户名和密码登陆本机。");
		break;      		
	case 3:
		document.write("本账号只能在指定地址使用<br>This account can be used on the appointed address only."+pp+xip);
		break;
	case 4:
		document.write("本账号费用超支或时长流量超过限制<br>This account overspent or over time limit");
		break;
	case 5:
		document.write("本账号暂停使用<br>This account has been suspended");
		break;
	case 6:
		document.write("System buffer full");
		break;      		
	case 8:
		document.write("本账号正在使用,不能修改<br>This account is in use. Unable to change");
		break;
	case 9:
		document.write("新密码与确认新密码不匹配,不能修改<br>New password and confirmation do not match. Unable to change");
		break;
	case 10:
		document.write("密码修改成功<br>Password Changed Successfully");
		break;
	case 11:
		document.write("本账号只能在指定地址使用<br>This account can be used on the appointed address only :"+pp+mac);
		break;      	
		case 7:
	      	document.write(UT+UF+UM);
		break;
	case 14:
		document.write("注销成功 Logout successfully"+pp+UT+UF+UM);
		break;
	case 15:
		document.write("登录成功 Login successfully"+pp+UT+UF+UM);
		break;
		}}		
// -->

        </SCRIPT>
```
  
