package conf

import "github.com/spf13/viper"

// InitConfigYml init config that use yaml file.
func InitConfigYml() (*viper.Viper, error) {
	return initConfig("yaml")
}

// initConfig return new viper config and error if any.
func initConfig(t string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("app")
	v.SetConfigType(t)
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	v.WatchConfig()

	return v, nil
}
