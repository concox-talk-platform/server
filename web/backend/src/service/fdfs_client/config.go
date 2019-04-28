package fdfs_client

import "configs"

type config struct {
	trackerAddr []string
	maxConns    int
}

func newConfig() (*config, error) {
	config := &config{}
	config.trackerAddr = append(config.trackerAddr, configs.TrackerServerAddr)
	config.maxConns = configs.MaxConn

	return config, nil
}
