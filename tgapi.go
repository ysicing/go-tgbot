package main

import (
	"fmt"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

type TgApi struct {
	bot *tb.Bot
}

// NewBot new bot
func NewBot(token, endpoint string) (*TgApi, error) {
	tbcfg := tb.Settings{
		Poller: &tb.LongPoller{Timeout: time.Second * 15},
	}
	if len(endpoint) > 0 && strings.HasPrefix(endpoint, "http") {
		tbcfg.URL = endpoint
	}
	tbcfg.Token = token
	bot, err := tb.NewBot(tbcfg)
	if err != nil {
		return nil, err
	}
	return &TgApi{bot: bot}, nil
}

func (tg *TgApi) SendMsg(msg string, to int64, makedown bool) error {
	opt := &tb.SendOptions{}
	if makedown {
		opt.ParseMode = tb.ModeMarkdown
		msg = fmt.Sprintf("```\n%s\n```", msg)
	}
	_, err := tg.bot.Send(tb.ChatID(to), msg, opt)
	return err
}

func (tg *TgApi) SendFile(filePath, fileName, mime, caption string, to int64) error {
	_, err := tg.bot.Send(tb.ChatID(to), &tb.Document{
		File:     tb.FromDisk(filePath),
		Caption:  caption,
		MIME:     mime,
		FileName: fileName,
	})
	return err
}

func (tg *TgApi) SendImage(imagePath, caption string, to int64) error {
	_, err := tg.bot.Send(tb.ChatID(to), &tb.Photo{
		File:    tb.FromDisk(imagePath),
		Caption: caption,
	})
	return err
}
