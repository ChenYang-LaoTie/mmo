package apis

import (
	"fmt"
	"mmo/core"
	"mmo/pb"
	"zinx/ziface"
	"zinx/znet"

	"google.golang.org/protobuf/proto"
)

type MoveApi struct {
	znet.BaseRouter
}

func (m *MoveApi) Handle(request ziface.IRequest) {
	proto_msg := &pb.Position{}

	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("Move: Postion Unmarshal error: ", err)
		return
	}

	pid, err := request.GetConnection().GetProperty(Pid)
	if err != nil {
		fmt.Println("GetProperty pid error: ", err)
		return
	}

	// fmt.Printf("Player pid = %d, move(%f, %f, %f, %f)\n", pid, proto_msg.GetX(), proto_msg.GetY(), proto_msg.GetZ(), proto_msg.GetV())

	player, err := core.TheWorldManager.GetPlayerByPid(pid.(int32))
	if err != nil {
		fmt.Println(err)
		return
	}

	player.UpdatePos(proto_msg.GetX(), proto_msg.GetY(), proto_msg.GetZ(), proto_msg.GetV())
}
