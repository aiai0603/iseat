package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-yaml/yaml"
)

const (
	DeviceManagerBaseKey = "DeviceManagerBaseKey"
	etc_dir              = "/etc/"
	data_dir             = "/data/"
	ConfName             = "conf.yaml"
	CaName               = "ca.crt"
	LicName              = "emqx.lic"
	LoadFileType         = "relative"
)

var (
	config *DeviceManagerConf
	once   sync.Once
	crypto cipher.AEAD
)

type TlsConf struct {
	Certfile string `yaml:"certfile"`
	Keyfile  string `yaml:"keyfile"`
}

type AccountLockConf struct {
	MaximumFailAttempts int `yaml:"maximumFailAttempts"`
	DurationMinutes     int `yaml:"durationMinutes"`
	CountRangeMinutes   int `yaml:"countRangeMinutes"`
}

type BasicAuthConf struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
}

type DB struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	Port     int    `yaml:"port"`
}

type DeviceManagerConf struct {
	Basic struct {
		Port    int      `yaml:"port"`
		Tlsconf *TlsConf `yaml:"tls"`
		Aeskey  string   `yaml:"aeskey"`
	}
	Logger    *LoggerConfig  `yaml:"logger"`
	Db        *DB            `yaml:"db"`
	BasicAuth *BasicAuthConf `yaml:"basicAuth"`
	Jwt       struct {
		Secret            string `yaml:"secret"`
		Expiration        int    `yaml:"expiration"`
		RefreshSecret     string `yaml:"refreshSecret"`
		RefreshExpiration int    `yaml:"refreshExpiration"`
	}
	AccountLock *AccountLockConf `yaml:"accountLock"`
}

func LoadConf(confName string) ([]byte, error) {
	confDir, err := GetConfLoc()
	if err != nil {
		return nil, err
	}

	file := path.Join(confDir, confName)
	//	file := confDir + confName
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GetConf() *DeviceManagerConf {
	once.Do(initConf)
	return config
}

func initConf() {
	b, err := LoadConf(ConfName)
	if err != nil {
		panic(err)
	}
	cfg := DeviceManagerConf{}
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		panic(err)
	}

	if printable, err := json.Marshal(cfg); err != nil {
		Log.Infof("Init with configuration %+v", cfg)
	} else {
		Log.Infof("Init with configuration %s", printable)
	}
	config = &cfg

	if config.Basic.Aeskey != "" {
		if len(config.Basic.Aeskey) != 32 {
			panic(fmt.Errorf("invalid config `aeskey`, must be a 32 bytes string"))
		} else {
			c, err := aes.NewCipher([]byte(config.Basic.Aeskey))
			if err != nil {
				panic(fmt.Errorf("invalid config `aeskey`: %s", err))
			}
			crypto, err = cipher.NewGCM(c)
			if err != nil {
				panic(fmt.Errorf("invalid config `aeskey`: %s", err))
			}
		}
	}

	initLogConf(config.Logger)
	if e := readMsgDir(); nil != e {
		panic(fmt.Errorf("load msgDir fail:%v", e))
	}
}

func GetConfLoc() (string, error) {
	return GetLoc(etc_dir)
}

func GetDataLoc() (string, error) {
	return GetLoc(data_dir)
}

func absolutePath(subdir string) (dir string, err error) {
	subdir = strings.TrimLeft(subdir, `/`)
	subdir = strings.TrimRight(subdir, `/`)
	switch subdir {
	case "etc":
		dir = "/etc/device_manager/"
		break
	case "log":
		dir = "/var/log/device_manager/"
		break
	}
	if 0 == len(dir) {
		return "", fmt.Errorf("no find such file : %s", subdir)
	}
	return dir, nil
}

func GetLoc(subdir string) (string, error) {
	if "relative" == LoadFileType {
		return relativePath(subdir)
	}

	if "absolute" == LoadFileType {
		return absolutePath(subdir)
	}
	return "", fmt.Errorf("Unrecognized loading method.")
}

func relativePath(subdir string) (dir string, err error) {
	dir, err = os.Getwd()
	if err != nil {
		return "", err
	}

	if base := os.Getenv(DeviceManagerBaseKey); base != "" {
		Log.Infof("Specified device manager base folder at location %s.\n", base)
		dir = base
	}
	confDir := dir + subdir
	if _, err := os.Stat(confDir); os.IsNotExist(err) {
		lastdir := dir
		for len(dir) > 0 {
			dir = filepath.Dir(dir)
			if lastdir == dir {
				break
			}
			confDir = dir + subdir
			if _, err := os.Stat(confDir); os.IsNotExist(err) {
				lastdir = dir
				continue
			} else {
				//Log.Printf("Trying to load file from %s", confDir)
				return confDir, nil
			}
		}
	} else {
		//Log.Printf("Trying to load file from %s", confDir)
		return confDir, nil
	}

	return "", fmt.Errorf("dir %s not found, please set DeviceManagerBaseKey program environment variable correctly.", dir)
}
