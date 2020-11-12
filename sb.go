// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/ysicing/ext/logger"
	"github.com/ysicing/ext/utils/exfile"
	"os"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"github.com/ysicing/ext/utils/exos"
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
	Use:     "sb",
	Short:   "simple bot",
	Long:    "一个 Telegram 推送的小工具，用于调用 Bot API 发送告警等",
	Version: fmt.Sprintf("%s %s %s", Version, GitCommitHash, BuildDate),
	Run: func(cmd *cobra.Command, args []string) {
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

// NewBot new bot
func NewBot() *tgbotapi.BotAPI {
	endpoint := tgbotapi.APIEndpoint
	getep := viper.GetString("botapi")
	if len(getep) != 0 && strings.HasPrefix(getep, "http") {
		endpoint = getep
	}

	bot, err := tgbotapi.NewBotAPIWithAPIEndpoint(viper.GetString("tgbot"), endpoint)
	if err != nil {
		return nil
	}
	bot.Debug = viper.GetBool("mode.debug")
	return bot
}

func initConfig() {
	cfgpath := "./bot.yaml"
	if exos.IsLinux() {
		cfgpath = "/conf/bot.yaml"
	}
	if !exfile.CheckFileExistsv2(cfgpath) {
		defaultcfg := `tgbot: xxx # bot token
mode:
  debug: false
botapi: "https://api.telegram.org/bot%s/%s"
tgchan: "@chanid" # 频道name
tguser: userid # 用户id
`
		exfile.WriteFile(cfgpath, defaultcfg)
		logger.Slog.Fatal("需要修改配置文件")
	}
	viper.SetConfigFile(cfgpath)
	// viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		os.Exit(-1)
	}
	// reload
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		logger.Slog.Debug("config changed: ", in.Name)
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
		logger.Slog.Fatalf("send msg err: %v, msg: %v", err, msg)
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
		logger.Slog.Fatalf("send msg err: %v, msg: %v", err, filepath)
	}
}

func init() {
	logcfg := logger.Config{Simple: true, ConsoleOnly: true}
	logger.InitLogger(&logcfg)
	initConfig()
	rootCmd.PersistentFlags().StringVar(&msgtype, "type", "msg", "msg或者file，默认msg")
	rootCmd.PersistentFlags().StringVar(&msgvalue, "c", "simple bot", "msg信息或者文件路径")
	rootCmd.PersistentFlags().BoolVar(&msgchan, "chan", false, "msg发送chan")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		logger.Slog.Fatalf("执行失败: %v", err)
	}
}
