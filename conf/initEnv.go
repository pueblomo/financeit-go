package conf

import "github.com/spf13/viper"

func InitEnv() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	CheckError(err)
	viper.AutomaticEnv()
}

func GetDBPw() string {
	return viper.GetString("DBPASSWORD")
}

func GetDBUrl() string {
	return viper.GetString("DBURL")
}