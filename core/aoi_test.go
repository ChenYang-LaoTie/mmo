package core

import (
	"fmt"
	"mmo/pb"
	"testing"

	"google.golang.org/protobuf/proto"
)

var aoiManager = NewAOIManager(0, 250, 5, 0, 250, 5)

func TestNewAOIManager(t *testing.T) {

	fmt.Println(aoiManager)
}

func TestMy(t *testing.T) {
	r := make([]int, 0)

	r = append(r, []int{1, 2, 3}...)

	i := 1
	for k, v := range r {
		fmt.Println("k - ", k, " v - ", v)
		r = append(r, i)

		i++

		fmt.Println(len(r))
	}

	// aoiManager.grids[99].AddPlayer(1)
}

func TestSurroundGrids(t *testing.T) {
	for gid := range aoiManager.grids {
		surr := aoiManager.GetSurroundGridsByGId(gid)

		gIds := make([]int, 0, len(surr))

		for _, v := range surr {
			gIds = append(gIds, v.GId)
		}

		fmt.Println(gid, " - surround grid IDs are - ", gIds)
	}
}

func TestProto(t *testing.T) {
	person := &pb.Person{
		Name:  "chenyang",
		Id:    32,
		Email: "775517776@qq.com",
		Phones: []*pb.Person_PhoneNumber{
			{
				Number: "19136073671",
				Type:   pb.Person_MOBILE,
			},
			{
				Number: "18191071231",
				Type:   pb.Person_HOME,
			},
		},
	}

	data, err := proto.Marshal(person)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(data)
	fmt.Println(string(data))

	p := &pb.Person{}

	err = proto.Unmarshal(data, p)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(p)

}
