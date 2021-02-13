package oauth

import "github.com/spf13/viper"

type clientOptions struct {
	ClientID     string
	ClientSecret string
}

type oauth2Config struct {
	Google clientOptions
}

// Config : configurations for oauth2 authentication system
var Config oauth2Config
 
// ReadConfig reads oauth2 configurations from a yaml file
func ReadConfig(baseDir string) error {
	v := viper.New()
	v.SetConfigName("oauth2")
	v.SetConfigType("yml")
	v.AddConfigPath(baseDir)
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		return err
	}
	return v.Unmarshal(&Config)
}