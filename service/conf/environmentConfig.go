package conf

import (
	"fmt"
	"github.com/stevenroose/gonfig"
	"github.com/v2rayA/v2rayA/pkg/util/log"
	log2 "log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Params struct {
	Address                 string   `id:"address" short:"a" default:"0.0.0.0:2017" desc:"Listening address"`
	Config                  string   `id:"config" short:"c" desc:"v2rayA configuration directory"`
	V2rayBin                string   `id:"v2ray-bin" desc:"Executable v2ray binary path. Auto-detect if put it empty."`
	V2rayConfigDirectory    string   `id:"v2ray-confdir" desc:"Additional v2ray config directory, files in it will be combined with config generated by v2rayA"`
	WebDir                  string   `id:"webdir" desc:"v2rayA web files directory. use embedded files if not specify."`
	VlessGrpcInboundCertKey []string `id:"vless-grpc-inbound-cert-key" desc:"Specify the certification path instead of automatically generating a self-signed certificate. Example: /etc/v2raya/grpc_certificate.crt,/etc/v2raya/grpc_private.key"`
	ForceIPV6On             bool     `id:"force-ipv6-on" desc:"Force to turn ipv6 support on"`
	PassCheckRoot           bool     `desc:"Skip privilege checking. Use it only when you cannot start v2raya but confirm you have root privilege"`
	ResetPassword           bool     `id:"reset-password"`
	LogLevel                string   `id:"log-level" default:"info" desc:"Optional values: trace, debug, info, warn or error"`
	LogFile                 string   `id:"log-file" desc:"The path of log file"`
	LogMaxDays              int64    `id:"log-max-days" default:"3" desc:"Maximum number of days to keep log files"`
	LogDisableColor         bool     `id:"log-disable-color"`
	Lite                    bool     `id:"lite" desc:"Lite mode for non-root and non-linux users"`
	ShowVersion             bool     `id:"version"`
}

var params Params

var dontLoadConfig bool

func initFunc() {
	defer SetServiceControlMode()
	if dontLoadConfig {
		return
	}
	err := gonfig.Load(&params, gonfig.Conf{
		FileDisable:       true,
		FlagIgnoreUnknown: false,
		EnvPrefix:         "V2RAYA_",
	})
	if err != nil {
		if err.Error() != "unexpected word while parsing flags: '-test.v'" {
			log2.Fatal(err)
		}
	}
	if params.ShowVersion {
		fmt.Println(Version)
		os.Exit(0)
	}
	if params.Lite {
		params.PassCheckRoot = true
	}
	if params.Config == "" {
		if params.Lite {
			params.Config = "$HOME/.config/v2raya"
		} else {
			params.Config = "/etc/v2raya"
		}
	}
	// replace all dots of the filename with underlines
	params.Config = filepath.Join(
		filepath.Dir(params.Config),
		strings.ReplaceAll(filepath.Base(params.Config), ".", "_"),
	)
	if strings.Contains(params.Config, "$HOME") {
		if h, err := os.UserHomeDir(); err == nil {
			params.Config = strings.ReplaceAll(params.Config, "$HOME", h)
		}
	}
	logWay := "console"
	if params.LogFile != "" {
		logWay = "file"
	}
	log.InitLog(logWay, params.LogFile, params.LogLevel, params.LogMaxDays, params.LogDisableColor)
}

var once sync.Once

func GetEnvironmentConfig() *Params {
	once.Do(initFunc)
	return &params
}

func SetConfig(config Params) {
	params = config
}

func DontLoadConfig() {
	dontLoadConfig = true
}
