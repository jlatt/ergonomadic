package irc

import (
	"code.google.com/p/gcfg"
	"errors"
)

type PassConfig struct {
	Password string
}

func (conf *PassConfig) PasswordBytes() []byte {
	bytes, err := DecodePassword(conf.Password)
	if err != nil {
		Log.error.Fatalln("decode password error:", err)
	}
	return bytes
}

type Config struct {
	Server struct {
		PassConfig
		Database string
		Listen   []string
		Log      string
		MOTD     string
		Name     string
	}

	Operator map[string]*PassConfig

	Theater map[string]*PassConfig
}

func (conf *Config) Operators() map[Name][]byte {
	operators := make(map[Name][]byte)
	for name, opConf := range conf.Operator {
		operators[NewName(name)] = opConf.PasswordBytes()
	}
	return operators
}

func (conf *Config) Theaters() map[Name][]byte {
	theaters := make(map[Name][]byte)
	for s, theaterConf := range conf.Theater {
		name := NewName(s)
		if !name.IsChannel() {
			Log.error.Fatal("config uses a non-channel for a theater!")
		}
		theaters[name] = theaterConf.PasswordBytes()
	}
	return theaters
}

func LoadConfig(filename string) (config *Config, err error) {
	config = &Config{}
	err = gcfg.ReadFileInto(config, filename)
	if err != nil {
		return
	}
	if config.Server.Name == "" {
		err = errors.New("server.name missing")
		return
	}
	if config.Server.Database == "" {
		err = errors.New("server.database missing")
		return
	}
	if len(config.Server.Listen) == 0 {
		err = errors.New("server.listen missing")
		return
	}
	return
}
