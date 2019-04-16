package server_login

import "log"

type Runner struct {
	Controller             ControlChan
	Error                  ControlChan
	Data                   DataChan
	DataType               int
	Srv                    Pbsrv
	dataSize               int
	longLived              bool
	Dispatcher             dfn
	FirstLoginDataExecutor fn
	ImMsgDataDataExecutor  fn
	KeepAliveExecutor      fn
}

func NewRunner(size int, longlived bool, srv Pbsrv, d dfn, flde fn, imde fn, ke fn) *Runner {
	return &Runner{
		Controller:             make(chan string, 2),
		Error:                  make(chan string, 2),
		Data:                   make(chan interface{}, size),
		DataType:               0,
		Srv:                    srv,
		longLived:              longlived,
		dataSize:               size,
		Dispatcher:             d,
		FirstLoginDataExecutor: flde,
		ImMsgDataDataExecutor:  imde,
		KeepAliveExecutor:      ke,
	}
}

func (r *Runner) startDispatch() {
	defer func() {
		if !r.longLived {
			close(r.Controller)
			close(r.Data)
			close(r.Error)
		}
	}()

	for {
		select {
		// 1 是返回登录数据， 2是推送ImMsgData， 3 是keepalive
		case c := <-r.Controller:
			if c == READY_TO_DISPATCH {
				err := r.Dispatcher(r.Data, &r.DataType, r.Srv)
				log.Println("dispatcher done")
				if err != nil {
					r.Error <- CLOSE
				} else {
					log.Println("else done", r.DataType)
					switch r.DataType {
					case 1:
						log.Println("first login data get")
						r.Controller <- READY_TO_FIRST_LOGIN_DATA_EXECUTE
					case 2:
						r.Controller <- READY_TO_IM_MSG_DATA_EXECUTE
					case 3:
						r.Controller <- READY_TO_kEEP_ALIVE_EXECUTE
					}
				}
			}

			if c == READY_TO_FIRST_LOGIN_DATA_EXECUTE {
				err := r.FirstLoginDataExecutor(r.Data, r.Srv)
				if err != nil {
					r.Error <- CLOSE
				} else {
					r.Controller <- READY_TO_DISPATCH
				}
			}

			if c == READY_TO_IM_MSG_DATA_EXECUTE {
				err := r.ImMsgDataDataExecutor(r.Data, r.Srv)
				if err != nil {
					r.Error <- CLOSE
				} else {
					r.Controller <- READY_TO_DISPATCH
				}
			}

			if c == READY_TO_kEEP_ALIVE_EXECUTE {
				err := r.KeepAliveExecutor(r.Data, r.Srv)
				if err != nil {
					r.Error <- CLOSE
				} else {
					r.Controller <- READY_TO_DISPATCH
				}
			}
		case e := <-r.Error:
			if e == CLOSE {
				return
			}
		default:
		}
	}
}

func (r *Runner) StartAll() {
	r.Controller <- READY_TO_DISPATCH
	r.startDispatch()
}
