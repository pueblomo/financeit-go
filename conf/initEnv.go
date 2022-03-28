package conf

import "github.com/spf13/viper"

func InitEnv() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	CheckErrorFatal(err)
}

func GetDBUser() string {
	return viper.GetString("DBUSER")
}

func GetDBPw() string {
	return viper.GetString("DBPASSWORD")
}