package modules

import (
	"sync"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
)

func init() {
	bot.RegisterModule(instance)
}

var instance = &ar{}
var tem map[string]string

type ar struct {
}

func (a *ar) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "dezhiShen.reply",
		Instance: instance,
	}
}

func (a *ar) Init() {
}

func (a *ar) PostInit() {
}

func (a *ar) Serve(b *bot.Bot) {
	b.OnGroupMessage(func(c *client.QQClient, msg *message.GroupMessage) {
		requset := plugins.NewMessageRequsetFromGroupMessage(msg)
		requset.QQClient = c
		m, err := onMessage(requset)
		if err != nil {
			go c.SendGroupMessage(msg.GroupCode, message.NewSendingMessage().Append(message.NewText(err.Error())))
		}
		if m != nil && m.Elements != nil && len(m.Elements) > 0 {
			go c.SendGroupMessage(msg.GroupCode, m)
		}

	})

	b.OnPrivateMessage(func(c *client.QQClient, msg *message.PrivateMessage) {
		requset := plugins.NewMessageRequsetFromPrivateMessage(msg)
		requset.QQClient = c
		m, err := onMessage(requset)
		if err != nil {
			go c.SendGroupMessage(msg.Sender.Uin, message.NewSendingMessage().Append(message.NewText(err.Error())))
		}
		if m != nil && m.Elements != nil && len(m.Elements) > 0 {
			go c.SendPrivateMessage(msg.Sender.Uin, m)
		}
	})
}

func (a *ar) Start(bot *bot.Bot) {
}

func (a *ar) Stop(bot *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}

func onMessage(requset *plugins.MessageRequest) (*message.SendingMessage, error) {
	m := message.NewSendingMessage()
	for _, pluginID := range plugins.GlobalOnMessagePluginIDs {
		plugin := plugins.GlobalOnMessagePlugins[pluginID]
		if plugin.IsFireEvent(requset) {
			resp, err := plugin.OnMessageEvent(requset)
			if err != nil {
				return m, err
			}
			if resp != nil {
				for _, element := range resp.Elements {
					m.Append(element)
				}
			}
			if !plugin.IsFireNextEvent(requset) {
				break
			}
		}
	}
	return m, nil
}
