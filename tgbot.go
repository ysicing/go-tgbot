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
		bot, err := NewBot(viper.GetString("token"), viper.GetString("endpoint"))
		if err != nil {
			zlog.Error("bot init err: %v", err)
			return
		}
		var wg sync.WaitGroup
		wg.Add(1)
		if msgtype == "msg" {
			SendMsg(msgvalue, &wg, msgchan)
		} else {
			SendFile(msgvalue, &wg)
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
		defaultcfg := `tgbot: xxx # bot token
mode:
  debug: false
botapi: "https://api.telegram.org/bot%s/%s"
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
func SendMsg(msg string, wg *sync.WaitGroup, chanmsg ...bool) {
	defer wg.Done()
	botapi := NewBot()
	var err error
	if len(chanmsg) > 0 && chanmsg[0] {
		tgchan := viper.GetString("tgchan")
		tgmsg := tgbotapi.NewMessageToChannel(tgchan, msg)
		_, err = botapi.Send(tgmsg)
	} else {
		tguser := viper.GetInt64("tguser")
		tgmsg := tgbotapi.NewMessage(tguser, msg)
		_, err = botapi.Send(tgmsg)
	}
	if err != nil {
		zlog.Error("send msg err: %v, msg: %v", err, msg)
	}
}

func SendFile(filepath string, wg *sync.WaitGroup) {
	defer wg.Done()
	botapi := NewBot()
	var err error
	tguser := viper.GetInt64("tguser")
	tgfile := tgbotapi.NewDocumentUpload(tguser, filepath)
	_, err = botapi.Send(tgfile)
	if err != nil {
		zlog.Error("send msg err: %v, msg: %v", err, filepath)
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
