package core

import (
	"fmt"
	"sync"
)

type Grid struct {
	IdX int

	IdY int
	GId int

	MinX int

	MaxX int

	MinY int

	MaxY int

	playerIds map[int32]bool

	pIdLock sync.RWMutex
}

func NewGrid(idX, idY, gId, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		IdX:       idX,
		IdY:       idY,
		GId:       gId,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIds: make(map[int32]bool),
	}
}

func (g *Grid) AddPlayer(playerId int32) {
	g.pIdLock.Lock()
	defer g.pIdLock.Unlock()

	g.playerIds[playerId] = true
}

func (g *Grid) RemovePlayer(playerId int32) {
	g.pIdLock.Lock()
	defer g.pIdLock.Unlock()

	delete(g.playerIds, playerId)
}

func (g *Grid) GetPlayerIds() []int32 {
	playerIds := make([]int32, 0)

	for k, v := range g.playerIds {
		if v {
			playerIds = append(playerIds, k)
		}
	}

	return playerIds
}

func (g *Grid) String() string {
	return fmt.Sprintf("IdX: %d, IdY: %d, Grid id: %d, minX: %d, maxX: %d, minY: %d, maxX: %d", g.IdX, g.IdY, g.GId, g.MinX, g.MaxX, g.MinY, g.MaxY)
}
