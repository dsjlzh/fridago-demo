package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/dsjlzh/fridago"
)

var chQuit chan bool

func on_message(message string, data []byte) {
	fmt.Println(message)
}

func main() {
	chQuit = make(chan bool)
	dev, err := fridago.GetUsbDevice()
	if err != nil {
		fmt.Println(err)
		return
	}
	pid, err := dev.Spawn("com.android.chrome")
	if err != nil {
		fmt.Println(err)
		return
	}
	sess, err := dev.Attach(pid)
	if err != nil {
		fmt.Println(err)
		return
	}
	agent, err := ioutil.ReadFile("test.js")
	if err != nil {
		fmt.Println(err)
		return
	}
	script, err := sess.CreateScriptSync("example", string(agent))
	if err != nil {
		fmt.Println(err)
		return
	}
	script.On("message", on_message)
	script.Load()

	go func() {
		for {
			_, err := sess.Dev.FindProcessByPidSync(sess.Pid)
			if err != nil {
				fmt.Println("process has quitted")
				chQuit <- true
			}
			time.Sleep(time.Second * 1)
		}
	}()

	tick := time.Tick(5 * time.Second)
	dev.Resume(pid)
	select {
	case <-tick:
		fmt.Println("timeout")
		script.UnLoad()
		sess.Detach()
		dev.Kill(pid)
	case <-chQuit:
		fmt.Println("quitted")
	}

	fmt.Println("done")
}
