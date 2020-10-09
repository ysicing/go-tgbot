module github.com/ysicing/sb

go 1.15

require (
	github.com/fsnotify/fsnotify v1.4.9
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.0.0-rc1
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.7.1
	github.com/ysicing/ext v0.0.0-20201006083949-adf2bbb1a9b4
	golang.org/x/sys v0.0.0-20201007082116-8445cc04cbdf // indirect
	golang.org/x/text v0.3.3 // indirect
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	gopkg.in/fsnotify.v1 v1.4.7
)

replace github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.0.0-rc1 => github.com/kunnos/telegram-bot-api/v5 v5.0.0-rc1.0.20201009142551-2b6c852d6586
