package message

import (
	"Go-Bat/abstraction"
	"Go-Bat/api"
	"Go-Bat/config"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Message interface {
	// 发送消息
	Send(MessageString Data)
	// 接收消息
	receive(w http.ResponseWriter, r *http.Request)
	//	开始监听消息
	Start()
	//	已读消息
	read()

	//
	Serve()
}

type GoBat struct {
	name    string
	version float64
	time    string
}

type Data struct {
	User_id     int64  `json:"user_id"`
	Message     string `json:"message"`
	Auto_escape bool   `json:"auto_escape"`
}

var once sync.Once
var bat *GoBat
var Mess config.Messages
var MessageChan = make(chan []byte, 100)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

func NewGoBat() *GoBat {
	once.Do(func() {
		bat = &GoBat{name: "Go-Bat", version: 0.1, time: time.Now().Format("2006-01-02 15:04:05")}
	})
	return bat
}

func (b *GoBat) Send(d Data) {
	fmt.Println("d", d)
	marshal, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	resp, err := http.Post("http://127.0.0.1:5000/send_private_msg", "application/json", bytes.NewBuffer(marshal))
	if err != nil {
		panic(err)
	}
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println("res", string(all))
	defer resp.Body.Close()
}

// websocket
func (b *GoBat) receive(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	go func() {
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				continue
			}
			if messageType == 1 {
				err := json.Unmarshal(message, &Mess)
				fmt.Println(Mess)
				if string(message) != "" {
					b.Send(Data{User_id: Mess.User_id, Message: Mess.Message, Auto_escape: false})
				}
				if err != nil {
					log.Panicln(err)
				}
				Gobat := abstraction.GoBat{}
				Gobat.SetStrategy(new(api.PrivateText))
				Gobat.Deal(Mess)
				b.read()
				Mess = config.Messages{}
			}
		}
	}()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			continue
		}
		if messageType == 1 {
			MessageChan <- message
			fmt.Println("送到通道")
			if err != nil {
				log.Panicln(err)
			}

		}
	}

}

// 启动监听
func (b *GoBat) Start() {
	log.Printf("[INFO]: %v  %v  %v", b.name, b.version, "机器人启动")
	http.HandleFunc("/", b.receive)
	err := http.ListenAndServe(":5700", nil)
	if err != nil {
		log.Panicln(err)
	}
	b.read()
}

func (b *GoBat) Serve() {
	Gobat := new(abstraction.GoBat)
	for {
		select {
		case c := <-MessageChan:
			// 如果chan1成功读到数据，则进行该case处理语句
			err := json.Unmarshal(c, &Mess)
			if err != nil {
				log.Panicln(err)
			}
			fmt.Println("Mess", Mess)
			if Mess.Message != "" {
				b.read()
				Gobat.SetStrategy(new(api.PrivateText))
				fmt.Println(Gobat.Deal(Mess).(string))
				b.Send(Data{User_id: Mess.User_id, Message: Gobat.Deal(Mess).(string), Auto_escape: false})
				Mess = config.Messages{}
				//b.Send(Data{User_id: Mess.User_id, Message: Mess.Message, Auto_escape: false})
			}
		default:
			// 如果上面都没有成功，则进入default处理流程
		}

	}
}

func (b *GoBat) read() {
	_, err := http.Get("http://127.0.0.1:5000/get_forward_msg?message_id=" + strconv.FormatInt(Mess.Message_id, 10))
	if err != nil {
		log.Panicln(err)
		return
	}
}

func (b *GoBat) GetName() string {
	return b.name
}

func (b *GoBat) GetVersion() float64 {
	return b.version
}