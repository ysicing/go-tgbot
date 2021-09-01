// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/zos"
	"github.com/ergoapi/zlog"
	"github.com/spf13/cobra"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/ysicing/sb/api"
)

var (
	msgtype       string
	msgvalue      string
	msgchan       bool
	Version       string
	BuildDate     string
	GitCommitHash string
)

var rootCmd = &cobra.Command{
	Use:     "bot",
	Short:   "simple bot",
	Long:    "一个 Telegram 推送的小工具，用于调用 Bot API 发送告警等",
	Version: fmt.Sprintf("%s %s %s", Version, GitCommitHash, BuildDate),
	Run: func(cmd *cobra.Command, args []string) {
		bot, err := api.NewBot(viper.GetString("token"), viper.GetString("endpoint"))
		if err != nil {
			zlog.Error("bot init err: %v", err)
			return
		}
		var wg sync.WaitGroup
		wg.Add(1)
		if msgtype == "msg" {
			SendMsg(msgvalue, bot, &wg, msgchan)
		} else if msgtype == "file" {
			SendFile(msgvalue, bot, &wg)
		} else {
			SendImage(msgvalue, bot, &wg)
		}
		wg.Wait()
	},
}

func initConfig() {
	cfgpath := "./bot.yaml"
	if zos.IsLinux() {
		cfgpath = "/conf/bot.yaml"
	}
	if !file.CheckFileExists(cfgpath) {
		defaultcfg := `token: xxx # bot token
mode:
  debug: false
endpoint: "https://api.telegram.org/bot%s/%s"
tgchan: "@chanid" # 频道name
tguser: userid # 用户id
`
		file.Writefile(cfgpath, defaultcfg)
		zlog.Fatal("需要修改配置文件")
	}
	viper.SetConfigFile(cfgpath)
	// viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		os.Exit(-1)
	}
	// reload
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		zlog.Debug("config changed: ", in.Name)
	})
}

// SendMsg send msg
func SendMsg(msg string, botapi *api.TgApi, wg *sync.WaitGroup, chanmsg ...bool) {
	defer wg.Done()
	if len(chanmsg) > 0 && chanmsg[0] {
		// tgchan := viper.GetString("tgchan")

	} else {
		tguser := viper.GetInt64("tguser")
		err := botapi.SendMsg(msg, tguser, true)
		if err != nil {
			zlog.Error("send msg err: %v", err)
		}
	}
}

func SendFile(filepath string, botapi *api.TgApi, wg *sync.WaitGroup) {
	defer wg.Done()
	tguser := viper.GetInt64("tguser")
	if len(filepath) != 0 {
		info, err := os.Stat(filepath)
		if err != nil {
			zlog.Error("file stat err: %v", err)
			return
		}
		err = botapi.SendFile(filepath, info.Name(), "", info.Name(), tguser)
		if err != nil {
			zlog.Error("failed to send file: %v", err)
		}
	}
}

func SendImage(filepath string, botapi *api.TgApi, wg *sync.WaitGroup) {
	defer wg.Done()
	tguser := viper.GetInt64("tguser")
	if len(filepath) != 0 {
		info, err := os.Stat(filepath)
		if err != nil {
			zlog.Error("file stat err: %v", err)
			return
		}
		err = botapi.SendImage(filepath, info.Name(), tguser)
		if err != nil {
			zlog.Error("failed to send file: %v", err)
		}
	}
}

func init() {
	logcfg := zlog.Config{Simple: true, WriteLog: false, ServiceName: "tgbot"}
	zlog.InitZlog(&logcfg)
	initConfig()
	rootCmd.PersistentFlags().StringVar(&msgtype, "type", "msg", "msg或者file，默认msg")
	rootCmd.PersistentFlags().StringVar(&msgvalue, "c", "simple bot", "msg信息或者文件路径")
	rootCmd.PersistentFlags().BoolVar(&msgchan, "chan", false, "msg发送chan")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		zlog.Error("执行失败: %v", err)
	}
}
