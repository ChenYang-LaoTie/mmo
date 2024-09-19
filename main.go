package main

import (
	"fmt"
	"mmo/apis"
	"mmo/core"

	"github.com/ganhuone/zinx/ziface"
	"github.com/ganhuone/zinx/znet"
)

func main() {
	server, err := znet.NewServer("[zinx V0.1]")

	if err != nil {
		fmt.Println(err)
		return
	}

	server.SetOnConnStart(OnConnestionAdd)
	server.SetOnConnStop(OnConnectionLost)

	server.AddRouter(2, &apis.WorldChatApi{})
	server.AddRouter(3, &apis.MoveApi{})

	server.Serve()

}

func OnConnestionAdd(conn ziface.IConnection) {

	player := core.NewPlayer(conn)
	//给客户端发消息 MsgId = 1 同步当前Player给ID客户端
	player.SyncPid()
	//给客户端发消息 MsgId = 2 同步当前Player的初始值给客户端
	player.BroadCastStartPosition()

	core.TheWorldManager.AddPlayer(player)

	conn.SetProperty(apis.Pid, player.Pid)

	player.SyncSurrounding()

	fmt.Println("--> player id = ", player.Pid, "is arrived")
}

func OnConnectionLost(conn ziface.IConnection) {
	pid, err := conn.GetProperty(apis.Pid)
	if err != nil {
		fmt.Println(err)
		return
	}

	player, err := core.TheWorldManager.GetPlayerByPid(pid.(int32))
	if err != nil {
		fmt.Println(err)
		return
	}

	player.Offline()
}
