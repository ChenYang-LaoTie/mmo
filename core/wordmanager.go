package core

import (
	"errors"
	"sync"
)

type WorldManager struct {
	AoiManager *AOIManager

	Players map[int32]*Player

	pLock sync.RWMutex
}

var TheWorldManager *WorldManager

func init() {
	TheWorldManager = &WorldManager{
		AoiManager: NewAOIManager(0, 1000, 100, 0, 1000, 100),

		Players: make(map[int32]*Player),
	}

}

func (w *WorldManager) AddPlayer(player *Player) {
	w.pLock.Lock()
	w.Players[player.Pid] = player
	w.pLock.Unlock()

	w.AoiManager.AddToGridByPos(player.Pid, player.X, player.Z)
}

func (w *WorldManager) RemovePlayer(pId int32) {
	w.pLock.Lock()
	defer w.pLock.Unlock()
	player := w.Players[pId]
	w.AoiManager.RemoveFromGridByPos(player.Pid, player.X, player.Z)

	delete(w.Players, pId)

}

func (w *WorldManager) GetPlayerByPid(pid int32) (*Player, error) {
	w.pLock.Lock()
	defer w.pLock.Unlock()

	player, ok := w.Players[pid]
	if !ok {
		return nil, errors.New("player not found")
	}

	return player, nil
}

func (w *WorldManager) GetAllPlayers() []*Player {
	players := make([]*Player, 0)

	for _, v := range w.Players {
		players = append(players, v)
	}

	return players
}
