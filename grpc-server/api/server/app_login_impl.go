/*
@Time : 2019/4/12 19:29 
@Author : yanKoo
@File : talk_cloud_app_login_impl
@Software: GoLand
@Description:
*/
package server

import (
	"context"
	"database/sql"
	pb "server/grpc-server/api/talk_cloud"
	"server/grpc-server/log"
	"time"

	tu "server/grpc-server/dao/user"
	tuc "server/grpc-server/dao/user_cache"
	"sync"
)

type TalkCloudServiceImpl struct{}

// 用户登录
func (tcs *TalkCloudServiceImpl) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRsp, error) {
	log.Log.Debugf("%s enter login with time: %d", req.Name, time.Now().UnixNano())
	//　验证用户名是否存在以及密码是否正确，然后就生成一个uuid session, 把sessionid放进metadata返回给客户端，
	//  然后之后的每一次连接都需要客户端加入这个metadata，使用拦截器，对用户进行鉴权
	if req.Name == "" || req.Passwd == "" {
		return &pb.LoginRsp{Res: &pb.Result{Code: 422, Msg: "用户名或密码不能为空"}}, nil
	}

	res, err := tu.SelectUserByKey(req.Name)
	if err != nil && err != sql.ErrNoRows {
		log.Log.Printf("App login >>> err != nil && err != sql.ErrNoRows <<< error : %s", err)
		loginRsp := &pb.LoginRsp{
			Res: &pb.Result{
				Code: 500,
				Msg:  "User Login Process failed. Please try again later"},
		}
		return loginRsp, nil
	}

	if err == sql.ErrNoRows {
		log.Log.Printf("App login error : %s", err)
		loginRsp := &pb.LoginRsp{
			Res: &pb.Result{
				Code: 500,
				Msg:  "User is not exist error. Please try again later"},
		}
		return loginRsp, nil
	}

	if res.PassWord != req.Passwd {
		log.Log.Printf("App login error : %s", err)
		loginRsp := &pb.LoginRsp{
			Res: &pb.Result{
				Code: 500,
				Msg:  "User Login pwd error. Please try again later"},
		}
		return loginRsp, nil
	}
	userInfo := &pb.Member{
		Id:          int32(res.Id),
		IMei:        res.IMei,
		UserName:    res.UserName,
		NickName:    res.NickName,
		UserType:    int32(res.UserType),
		LockGroupId: int32(res.LockGroupId),
		Online:      tuc.USER_ONLINE, // 登录就在线
	}

	var (
		errMap   = &sync.Map{}
		wg       sync.WaitGroup
		fList    = make(chan *pb.FriendsRsp, 1)
		gList    = make(chan *pb.GroupListRsp, 1)
		existErr bool
	)
	defer func() {
		close(fList)
		close(gList)
	}()

	// 1. 处理登录session
	//go processSession(req, errMap, &wg)

	// 2. 将用户信息添加进redis
	addUserInfoToCache(userInfo, &wg)

	// 3. 获取好友列表
	getFriendList(int32(res.Id), fList, errMap, &wg)

	// 4. 群组列表
	getGroupList(int32(res.Id), gList, errMap, &wg)

	//wg.Wait()

	//遍历该map，参数是个函数，该函数参的两个参数是遍历获得的key和value，返回一个bool值，当返回false时，遍历立刻结束。
	errMap.Range(func(k, v interface{}) bool {
		err = v.(error)
		if err != nil {
			log.Log.Println(k, " gen error: ", err)
			existErr = true
			return false
		}
		return true
	})
	if existErr {
		return &pb.LoginRsp{Res: &pb.Result{Code: 500, Msg: "process error, please try again"}}, nil
	}

	loginRep := &pb.LoginRsp{
		UserInfo:   userInfo,
		FriendList: (<-fList).FriendList,
		GroupList:  (<-gList).GroupList,
		Res:        &pb.Result{Code: 200, Msg: req.Name + " login successful"},
	}

	log.Log.Debugf("login done with time : %d", time.Now().UnixNano())
	return loginRep, nil
}
