## api 定义

### 1. 创建下级账户   

请求协议：`http`

请求方式：`POST`

请求地址：`/account` 

请求参数：请求body必须包含username和pwd和用户类型
``` json
{
	"confirm_pwd" :"123456",    // 不能为空
	"pid": 5,    				// 不能为空
	"username": "liuyang06",   // 不能为空
	"nick_name": "nana123",   // 不能为空
	"pwd": "123456",     // 不能为空
	"role_id": 3,    // 不能为空
	"email": "123456789@qq.com",
	"privilege_id": 0,
	"contact": "",
	"state": "",
	"phone": "",
	"remark": "",
	"address": ""
}
```

返回参数：body中的session_id
``` json
{
	"success": true,
	"account_id": 1430   // 新创建的用户的id
}
```

### 2. 登录
请求协议：`http`

请求方式：`POST`

请求地址：`/account/login.do/account_name` 

请求参数：请求body中：username和pwd
``` json
{
	"username" : "account_name",
	"pwd" : "123456",
}
```

返回参数：body中：
``` json
{
	"success": true,
	"session_id": "c9f9173c-7cc8-44c3-81a8-7c72d9863f9a"
}
```

### 3. 退出
请求协议：`http`

请求方式：`POST`

请求地址：`/account/logout.do/account_name` 

请求参数：请求body中：username和pwd，请求头中添加返回的session_id
``` json
{
	"username" : "account_name",
}
```

返回参数：body中：
``` json
{
	"success": true,
	"msg": "delete session id is successful",
}
```

### 4. 账户信息以及下面的所有群组
请求协议：`http`

请求方式：`GET`

请求地址：`/account/:account_name`   (`:account_name`代表账户名称，比如访问：`172.16.0.74:8080/account/tiger`,其中tiger为账户名称)  

请求参数：**请求头中添加返回的session_id**  

返回参数：body中：  // 具体需要哪些待定
``` json
{

	"message" :" 获取用户信息成功",
	"account_info" : "账户的信息",
	"group_list" : "群组的信息",
	"device_list" ： "账户下所有的设备"
}
```

### 5. 修改账户信息

请求协议：`http`

请求方式：`POST`

请求地址：`/account/info/update`

请求参数：**请求头中添加返回的session_id**  

返回参数：body中：  (以下字段json格式，可以为空)
```
	Id       string `json:"id"`
	NickName string `json:"nick_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	Remark   string `json:"remark"`
	Contact  string `json:"contact"`
```

### 6. 修改账户密码
请求协议：`http`

请求方式：`POST`

请求地址：`/account/pwd/update`

请求参数：**请求头中添加返回的session_id**  

返回参数：body中：  // 具体需要哪些待定

``` json
{
	Id string `json:"id"`    // 账户id
	OldPwd     string `json:"old_pwd"`
	NewPwd     string `json:"new_pwd"`
	ConfirmPwd string `json:"confirm_pwd"`
}
```

### 7. 获取账户下级目录
请求协议：`http`

请求方式：`POST`

请求地址：`/account_class/:accountId`

请求参数：**请求头中添加返回的session_id**  

返回参数：body中：  (以下字段json格式)
```
		"result":    "success",
		"tree_data": resp,

		Id          int             `json:"id"`
		AccountName string          `json:"account_name"`
		Children    []*AccountClass `json:"children"`
```

### 8. 获取下级某个用户的所有设备
请求协议：`http`

请求方式：`GET`

请求地址：`/account_device/:accountId`

请求参数：**请求头中添加返回的session_id**  

返回参数：body中：  (以下字段json格式)
```
		"account_info": ai,
		"devices":      deviceAll,
```

### 9. 转移设备

请求协议：`http`

请求方式：`POST`

请求地址：`/account_device/:accountId`

请求参数：**请求头中添加返回的session_id**  

返回参数：body中：  (以下字段json格式)
```
		"result": "success",
		"msg":    "trans successful"
```


### 10. 导入设备

请求协议：`http`

请求方式：`POST`

请求地址：`/device/import/SurpeRoot`

请求参数：**请求头中添加返回的session_id**  

返回参数：body中：  (以下字段json格式)
```
		"msg":"import device successful."
```


### 11. 创建组
请求协议：`http`

请求方式：`POST`

请求地址：`/group/`

请求参数：**cookie中添加返回的session_id**  // 调试的时候暂时注释


``` json
{
  "device_infos": [
    8,
    10,
    11   // 至少有一个设备
  ],
  "group_info": {
    "group_name": "重庆组",    // 必须的
    "account_id": 6,         // 必须的
    "status": "1",
    "c_time": "2019-03-18 10:28:26"
  }
}
```

返回参数：  

``` json
		"result": "success",
		"msg":    "Create group successfully",
```

### 12. 更新组名  目前web只用改组名
/group/update

### 13. 删除组 
/group/delete

### 14. 更新组成员
/group/devices/update

