package pubsub

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-logs/pkg/logs"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	clientID = "driverconsole"
)

type pubsubActor struct {
	ctx           actor.Context
	client        mqtt.Client
	subscriptions map[string]*subscribeMSG
	mux           sync.Mutex
}

var instance *pubsubActor
var once sync.Once

// getInstance create pubsub Gateway
func getInstance(ctx *actor.RootContext) *pubsubActor {

	once.Do(func() {
		instance = &pubsubActor{}
		// instance.mux = sync.Mutex{}
		instance.subscriptions = make(map[string]*subscribeMSG)
		if ctx == nil {
			ctx = actor.NewActorSystem().Root
		}
		props := actor.PropsFromFunc(instance.Receive)
		pid, err := ctx.SpawnNamed(props, "pubsub-actor")
		if err != nil {
			logs.LogError.Panic(err)
		}
		ctx.RequestFuture(pid, &ping{}, 10*time.Second).Wait()
	})
	return instance
}

// Init init pubsub instance
func Init(ctx *actor.RootContext) error {
	// defer time.Sleep(3 * time.Second)
	if getInstance(ctx) == nil {
		return fmt.Errorf("error instance")
	}
	return nil
}

type ping struct{}
type pong struct{}

type publishMSG struct {
	topic string
	msg   []byte
}

type subscribeMSG struct {
	pid   *actor.PID
	parse func([]byte) interface{}
}

// Publish function to publish messages in pubsub gateway
func Publish(topic string, msg []byte) {
	getInstance(nil).ctx.Send(instance.ctx.Self(), &publishMSG{topic: topic, msg: msg})
}

// Subscribe subscribe to topics
func Subscribe(topic string, pid *actor.PID, parse func([]byte) interface{}) error {
	instance := getInstance(nil)
	subs := &subscribeMSG{pid: pid, parse: parse}
	instance.mux.Lock()
	instance.subscriptions[topic] = subs
	instance.mux.Unlock()
	if !instance.client.IsConnected() {
		return fmt.Errorf("pubsub is not connected")
	}
	logs.LogBuild.Printf("subscription in topic -> %q -> %#v", topic, subs)
	instance.subscribe(topic, subs)
	return nil
}

func (ps *pubsubActor) subscribe(topic string, subs *subscribeMSG) error {
	handler := func(client mqtt.Client, m mqtt.Message) {
		fmt.Printf("local topic -> %q, local payload - > %s\n", m.Topic(), m.Payload())
		m.Ack()
		msg := subs.parse(m.Payload())
		// logs.LogBuild.Printf("parse payload-> %s", msg)
		if msg != nil {
			ps.ctx.Send(subs.pid, msg)
		}
	}
	if tk := instance.client.Subscribe(topic, 1, handler); !tk.WaitTimeout(3 * time.Second) {
		if err := tk.Error(); err != nil {
			return err
		}
	}
	return nil
}

// Receive function
func (ps *pubsubActor) Receive(ctx actor.Context) {
	logs.LogBuild.Printf("Message arrived in pubsubActor: %s, %T, %s",
		ctx.Message(), ctx.Message(), ctx.Sender())
	ps.ctx = ctx
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		logs.LogInfo.Printf("Starting, actor, pid: %v\n", ctx.Self())
		ps.client = client()
		if err := connect(ps.client); err != nil {
			time.Sleep(3 * time.Second)
			logs.LogError.Panic(err)
		}
		for k, v := range ps.subscriptions {
			ps.subscribe(k, v)
		}
	case *ping:
		if ctx.Sender() != nil {
			ctx.Respond(&pong{})
		}
	case *publishMSG:
		// fmt.Printf("publish msg: %s", msg.msg)
		logs.LogBuild.Printf("publish msg: %s", msg.msg)
		tk := ps.client.Publish(msg.topic, 1, false, msg.msg)
		if !tk.WaitTimeout(3 * time.Second) {
			if tk.Error() != nil {
				logs.LogError.Printf("end error: %s, with messages -> %v", tk.Error(), msg)
			} else {
				logs.LogError.Printf("timeout error with message -> %v", msg)
			}
		}
	case *actor.Stopping:
		ps.client.Disconnect(600)
		logs.LogError.Println("Stopping, actor is about to shut down")
	case *actor.Stopped:
		logs.LogError.Println("Stopped, actor and its children are stopped")
	case *actor.Restarting:
		logs.LogError.Println("Restarting, actor is about to restart")
	case error:
		fmt.Printf("error pubsub actor: %s\n", msg)
	}
}

func client() mqtt.Client {
	opt := mqtt.NewClientOptions().AddBroker("tcp://127.0.0.1:1883")
	opt.SetAutoReconnect(true)
	buff := make([]byte, 4)
	rand.Read(buff)
	opt.SetClientID(fmt.Sprintf("%s-%X", clientID, buff))
	opt.SetKeepAlive(30 * time.Second)
	opt.SetConnectRetryInterval(10 * time.Second)
	client := mqtt.NewClient(opt)
	return client
}

func connect(c mqtt.Client) error {
	tk := c.Connect()
	if !tk.WaitTimeout(10 * time.Second) {
		return fmt.Errorf("connect wait error")
	}
	if err := tk.Error(); err != nil {
		return err
	}
	return nil
}
