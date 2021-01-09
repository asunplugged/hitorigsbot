package twitchirc

import (
	"crypto/tls"
	"fmt"
	"time"

	"git.trj.tw/golang/mtfosbot/model"
	"github.com/go-irc/irc"

	"git.trj.tw/golang/mtfosbot/module/config"
)

var client *irc.Client
var queue *QueueList
var channels []string
var queueRunning bool

func init() {
	queueRunning = false
}

// InitIRC -
func InitIRC() {
	conf := config.GetConf()
	tlsConf := &tls.Config{}
	conn, err := tls.Dial("tcp", conf.Twitch.ChatHost, tlsConf)
	// conn, err := net.Dial("tcp", conf.Twitch.ChatHost)
	if err != nil {
		fmt.Println("create irc connect fail ", err)
		time.Sleep(time.Second * 3)
		go InitIRC()
		return
	}
	defer conn.Close()

	channels = make([]string, 0)
	queue = NewQueue()
	runQueue()
	ReJoin()

	config := irc.ClientConfig{
		Nick:    conf.Twitch.BotUser,
		Pass:    conf.Twitch.BotOauth,
		Handler: irc.HandlerFunc(ircHandle),
	}

	client = irc.NewClient(conn, config)

	err = client.Run()
	if err != nil {
		fmt.Println("twitch chat connect fail")
		// reconnect after 3sec
		time.Sleep(time.Second * 3)
		client = nil
		channels = channels[:0]
		queue.Clear()
		go InitIRC()
		return
	}
}

// SendMessage -
func SendMessage(ch, msg string) {
	if len(ch) == 0 {
		return
	}

	if indexOf(channels, ch) == -1 {
		return
	}

	m := &MsgObj{
		Command: "PRIVMSG",
		Params: []string{
			fmt.Sprintf("#%s", ch),
			fmt.Sprintf("%s", msg),
		},
	}
	queue.Add(m)
}

// ReJoin -
func ReJoin() {
	ch, err := model.GetAllTwitchChannel()
	if err != nil {
		return
	}
	LeaveAllChannel()
	for _, v := range ch {
		if v.Join {
			JoinChannel(v.Name)
		}
	}
}

// JoinChannel -
func JoinChannel(ch string) {
	if len(ch) == 0 {
		return
	}

	m := &MsgObj{
		Command: "JOIN",
		Params: []string{
			fmt.Sprintf("#%s", ch),
		},
	}
	queue.Add(m)
}

// LeaveChannel -
func LeaveChannel(ch string) {
	if len(ch) == 0 {
		return
	}

	if indexOf(channels, ch) == -1 {
		return
	}

	m := &MsgObj{
		Command: "PART",
		Params: []string{
			fmt.Sprintf("#%s", ch),
		},
	}
	queue.Add(m)
}

// LeaveAllChannel -
func LeaveAllChannel() {
	if len(channels) == 0 {
		return
	}
	for _, v := range channels {
		m := &MsgObj{
			Command: "PART",
			Params: []string{
				fmt.Sprintf("#%s", v),
			},
		}
		queue.Add(m)
	}
}

func runQueue() {
	if queueRunning == true {
		return
	}
	queueRunning = true
	go func() {
		cnt := 0
		for {
			if !queue.IsEmpty() && client != nil {
				m := queue.Get()
				msg := &irc.Message{}
				msg.Command = m.Command
				msg.Params = m.Params

				if m.Command == "JOIN" {
					if indexOf(channels, m.Params[0][1:]) != -1 {
						continue
					}
					channels = append(channels, m.Params[0][1:])
				} else if m.Command == "PART" {
					if indexOf(channels, m.Params[0][1:]) == -1 {
						continue
					}
					idx := indexOf(channels, m.Params[0][1:])
					channels = append(channels[:idx], channels[idx+1:]...)
				}
				fmt.Println("< ", msg.String())
				client.WriteMessage(msg)
			}
			cnt++
			if cnt > 3600 {
				// call rejoin
				ReJoin()
				cnt = 0
			}
			time.Sleep(time.Second * 1)
		}
	}()
}

func ircHandle(c *irc.Client, m *irc.Message) {
	fmt.Println("> ", m.String())
	if m.Command == "PING" {
		tmp := &irc.Message{
			Command: "PONG",
			Params: []string{
				m.Params[0],
			},
		}
		fmt.Println("< ", tmp.String())
		client.WriteMessage(tmp)
	}
}

func indexOf(c []string, data string) int {
	if len(c) == 0 {
		return -1
	}
	for k, v := range c {
		if v == data {
			return k
		}
	}
	return -1
}
