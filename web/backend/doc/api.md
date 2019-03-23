## api 定义

### 1. 注册   

请求协议：`http`

请求方式：`POST`

请求地址：`/account` 

请求参数：请求body必须包含username和pwd和用户类型
``` json
{
	"username" : "elephant",
	"pwd" : "123456",
	"role_type" : "1"
}
```

返回参数：body中的session_id
``` json
{
	"success": true,
	"session_id": "c9f9173c-7cc8-44c3-81a8-7c72d9863f9a"
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
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTI1NTYxMzgsImlhdCI6MTU1MjU1MjUzOCwic3ViIjoidHVydGxlMTIzNDU2In0.P18lLtyc3ht6YduUxVGq1dVXXlPwT8JPk03SRx4hQBw
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

请求参数：请求body中：username和pwd，请求头中添加登录返回的token，cookie中添加返回的session_id
``` json
{
	"username" : "account_name",
}
```

返回参数：body中：
``` json
{
	"success": true,
	"session_id": "c9f9173c-7cc8-44c3-81a8-7c72d9863f9a",
	"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTI1NTQ5NjAsImlhdCI6MTU1MjU1MTM2MCwic3ViIjoicGFuZGExMjM0NTYifQ.sqf7_08DlLjfqqAez5IZ2uUNFY8uHXr1pBBXam5Eh6Q"
}
```
``` go
message signOutReq {
	string sesssionId = 1
}

message signOutResp {
	bool success = 1
}
```

### 4. 账户信息以及下面的所有群组
请求协议：`http`

请求方式：`GET`

请求地址：`/account/:account_name`   (`:account_name`代表账户名称，比如访问：`172.16.0.74:8080/account/tiger`,其中tiger为账户名称)  

请求参数：**cookie中添加返回的session_id**  // 调试的时候暂时注释

返回参数：body中：  // 具体需要哪些待定
``` json
{

	"message" :" 获取用户信息成功",
	"account_info" : "账户的信息",
	"group_list" : "群组的信息",
	"device_list" ： "账户下所有的设备"
}
```
``` go 
message GetAccountInfoReq {
	 string userName
	 string sessionId
	 int    accountId
}

message GetAccountInfoResp {
	 int32 accountId          
	 int32 pid         
	 string userName    
	 string nickname    
	 string pwd         
	 string email       
	 string phone       
	 string remark      
	 string address     
	 int32 privilegeId 
	 int32 roleId      
	 string stat        
	 string llTime      
	 string cTime       
	 string changeTime 
	
	 string sessionId
}

message GetGroupsByAccountId {
	int   groupId
	string groupName
	int    accountId
	string status
	string cTime
}

message GetDevicesByAccountId {
	int    deviceId
	string iMei
	string aliasName
	int    accountId
	string status
	string activeStatus
	string bindStatus
	string cTime
}

message GetDeviceIdsByGroupId {
	int[] deviceIds
}
```
### 5. 修改账户信息

请求协议：`http`

请求方式：`POST`

请求地址：`/account/info/update`

请求参数：**cookie中添加返回的session_id**  // 调试的时候暂时注释

返回参数：body中：  // 具体需要哪些待定

``` go
{
	Id       string `json:"id"`
	NickName string `json:"nick_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	Remark   string `json:"remark"`
}
```

### 6. 修改账户密码
请求协议：`http`

请求方式：`POST`

请求地址：`/account/pwd/update`

请求参数：**cookie中添加返回的session_id**  // 调试的时候暂时注释

返回参数：body中：  // 具体需要哪些待定

``` json
{
	Id string `json:"id"`    // 账户id
	OldPwd     string `json:"old_pwd"`
	NewPwd     string `json:"new_pwd"`
	ConfirmPwd string `json:"confirm_pwd"`
}
```

### 7. 创建组
请求协议：`http`

请求方式：`POST`

请求地址：`/group/`

请求参数：**cookie中添加返回的session_id**  // 调试的时候暂时注释


``` json
{
  "device_ids": [
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

返回参数：  // 具体需要哪些待定

``` json
{"result":true,"type":"create group","message":"create group success"}
```


