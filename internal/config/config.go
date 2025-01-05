package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	SettingGracePeriodDay    int
	SettingInterestPerAnnum  float64
	SettingLateFeePercentage float64
	DataSource               string
}
