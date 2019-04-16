package server_login

import pb "api/talk_cloud"

const (
	READY_TO_DISPATCH                 = "d"
	READY_TO_FIRST_LOGIN_DATA_EXECUTE = "flde"
	READY_TO_IM_MSG_DATA_EXECUTE      = "imde"
	READY_TO_kEEP_ALIVE_EXECUTE       = "kae"
	CLOSE                             = "c"
)

type ControlChan chan string

type DataChan chan interface{}

type Pbsrv pb.TalkCloud_DataPublishServer

type fn func(dc DataChan, srv pb.TalkCloud_DataPublishServer) error

type dfn func(dc DataChan, dt *int, srv pb.TalkCloud_DataPublishServer) error
