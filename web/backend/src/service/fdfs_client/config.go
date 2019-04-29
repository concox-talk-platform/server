package fdfs_client

import cfgWs "configs/web_server"

type config struct {
	trackerAddr []string
	maxConns    int
}

func newConfig() (*config, error) {
	config := &config{}
	config.trackerAddr = append(config.trackerAddr, cfgWs.TrackerServerAddr)
	config.maxConns = cfgWs.MaxConn

	return config, nil
}
