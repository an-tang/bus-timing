package configuration

import (
	"github.com/spf13/viper"
)

var Config Configs

type Configs struct {
	// PostgreSQL   PostgreSQL `mapstructure:"postgresql"`
	// Redis        Redis      `mapstructure:"redis"`
	Server       Server      `mapstructure:"server"`
	UWaveConfig  UWaveConfig `mapstructure:"uwave"`
	SecretKeyJWT string      `mapstructure:"secret_key_jwt"`
}

type Server struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	WriteTimeout int    `mapstructure:"write_timeout"`
	IdleTimeout  int    `mapstructure:"idle_timeout"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
}

// type PostgreSQL struct {
// 	Host         string `mapstructure:"host"`
// 	Port         string `mapstructure:"port"`
// 	SSLMode      string `mapstructure:"ssl_mode"`
// 	UserName     string `mapstructure:"username"`
// 	Password     string `mapstructure:"password"`
// 	Database     string `mapstructure:"database"`
// 	MaxOpenConns int    `mapstructure:"max_open_conns"`
// 	MaxIdleConns int    `mapstructure:"max_idle_conns"`
// }

type UWaveConfig struct {
	Endpoint string `mapstructure:"endpoint"`
}

type Redis struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

func LoadConfig(path string) (err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&Config)
	return
}
