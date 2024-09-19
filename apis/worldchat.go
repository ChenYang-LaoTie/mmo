package apis

import (
	"fmt"
	"mmo/core"
	"mmo/pb"

	"github.com/ganhuone/zinx/ziface"
	"github.com/ganhuone/zinx/znet"

	"google.golang.org/protobuf/proto"
)

const (
	Pid = "pid"
)

type WorldChatApi struct {
	znet.BaseRouter
}

func (w *WorldChatApi) Handle(request ziface.IRequest) {
	protoMsg := &pb.Talk{}

	err := proto.Unmarshal(request.GetData(), protoMsg)
	if err != nil {
		fmt.Println("Talk error: ", err)
		return
	}

	pid, err := request.GetConnection().GetProperty(Pid)
	if err != nil {
		fmt.Println("Talk error: ", err)
		return
	}

	player, err := core.TheWorldManager.GetPlayerByPid(pid.(int32))
	if err != nil {
		fmt.Println("Talk error: ", err)
		return
	}

	player.Talk(protoMsg.Content)
}
