package core

import (
	"errors"
	"fmt"
	"math/rand"
	"mmo/pb"
	"sync"
	"zinx/ziface"

	"google.golang.org/protobuf/proto"
)

const (
	SyncPid          = 1   //同步玩家本次的登陆的ID
	Talk             = 2   //世界聊天
	Position         = 3   //移动
	BroadCast        = 200 //广播消息(Tp 1 世界聊天 Tp 2 坐标 Tp 3 动作)
	SyncPidDisappear = 201 //广播消息 掉线/aoi消失视野
	SyncPlayers      = 202 //同步周围的人位置信息
)

type Player struct {
	Pid  int32
	Conn ziface.IConnection
	X    float32
	Y    float32
	Z    float32
	V    float32
}

var PidGen int32 = 1
var IdLock sync.Mutex

func NewPlayer(conn ziface.IConnection) *Player {
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	return &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)), //(横)随机在160坐标点，基于X轴若干偏移
		Y:    0,                            //(垂直)
		Z:    float32(140 + rand.Intn(20)), //(纵)随机在140坐标点，基于Y轴若干偏移
		V:    0,                            //面向的角度
	}

}

func (p *Player) SendMsg(msgId uint32, data proto.Message) error {
	msg, err := proto.Marshal(data)
	if err != nil {
		return err
	}

	if p.Conn == nil {
		return errors.New("connection in player is nil")
	}

	err = p.Conn.SendMsg(msgId, msg)
	if err != nil {
		return err
	}

	return nil
}

func (p *Player) SyncPid() {
	data := &pb.SyncPid{
		Pid: p.Pid,
	}

	p.SendMsg(SyncPid, data)
}

func (p *Player) BroadCastStartPosition() {
	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	p.SendMsg(BroadCast, protoMsg)
}

func (p *Player) Talk(content string) {
	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1,
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	for _, v := range TheWorldManager.GetAllPlayers() {
		v.SendMsg(BroadCast, protoMsg)
	}

}

func (p *Player) SyncSurrounding() {
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	pids := TheWorldManager.AoiManager.GetPidsByPos(p.X, p.Z)

	player_proto_msg := make([]*pb.Player, 0, len(pids))

	for _, v := range pids {
		player, err := TheWorldManager.GetPlayerByPid(v)
		if err != nil {
			fmt.Println(err)
			continue
		}

		player.SendMsg(BroadCast, proto_msg)

		player_proto_msg = append(player_proto_msg, &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		})
	}

	p.SendMsg(SyncPlayers, &pb.SyncPlayers{
		Ps: player_proto_msg,
	})

}

func (p *Player) UpdatePos(x, y, z, v float32) {
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  4,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	pids := TheWorldManager.AoiManager.GetPidsByPos(p.X, p.Z)
	for _, v := range pids {
		player, err := TheWorldManager.GetPlayerByPid(v)
		if err != nil {
			fmt.Println(err)
			continue
		}

		player.SendMsg(BroadCast, proto_msg)

	}

}

func (p *Player) Offline() {
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}

	pids := TheWorldManager.AoiManager.GetPidsByPos(p.X, p.Z)

	for _, v := range pids {
		player, err := TheWorldManager.GetPlayerByPid(v)
		if err != nil {
			fmt.Println(err)
			continue
		}

		player.SendMsg(SyncPidDisappear, proto_msg)
	}

	TheWorldManager.AoiManager.RemoveFromGridByPos(p.Pid, p.X, p.Z)
	TheWorldManager.RemovePlayer(p.Pid)
	fmt.Println("Player pid = ", p.Pid, " offline...")
}
