package core

import "fmt"

type AOIManager struct {
	MinX int

	MaxX int

	CntsX int

	MinY int

	MaxY int

	CntsY int

	grids map[int]*Grid
}

func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	a := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,

		grids: make(map[int]*Grid),
	}

	gridWidth := a.gridWidth()
	gridLength := a.gridLength()

	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			gid := a.GetGIdFromXY(x, y)

			gminX := a.MinX + x*gridWidth
			gMaxX := gminX + gridWidth

			gMinY := a.MinY + y*gridLength
			gMaxY := gMinY + gridLength

			a.grids[gid] = NewGrid(x, y, gid, gminX, gMaxX, gMinY, gMaxY)
		}
	}

	return a
}

func (a *AOIManager) GetGIdFromXY(x, y int) int {
	return x + a.CntsX*y
}

func (a *AOIManager) gridWidth() int {
	return (a.MaxX - a.MinX) / a.CntsX
}

func (a *AOIManager) gridLength() int {
	return (a.MaxY - a.MinY) / a.CntsY
}

func (a *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager:\n MinX: %d, MaxX: %d, cntsX: %d, MinY: %d, MaxY: %d, cntsY: %d, Grids in AOIManager:\n", a.MinX, a.MaxX, a.CntsX, a.MinY, a.MaxY, a.CntsY)

	for _, grid := range a.grids {
		s += fmt.Sprintln(grid)
	}

	return s
}

func (a *AOIManager) GetSurroundGridsByGId(gId int) []*Grid {

	grid, ok := a.grids[gId]
	if !ok {
		return nil
	}

	gridSurround := make([]*Grid, 0)
	gridSurround = append(gridSurround, grid)

	xy := make([][2]int, 0)

	xy = append(xy, [2]int{grid.IdX - 1, grid.IdY + 1})
	xy = append(xy, [2]int{grid.IdX - 1, grid.IdY})
	xy = append(xy, [2]int{grid.IdX - 1, grid.IdY - 1})

	xy = append(xy, [2]int{grid.IdX, grid.IdY + 1})
	xy = append(xy, [2]int{grid.IdX, grid.IdY - 1})

	xy = append(xy, [2]int{grid.IdX + 1, grid.IdY + 1})
	xy = append(xy, [2]int{grid.IdX + 1, grid.IdY})
	xy = append(xy, [2]int{grid.IdX + 1, grid.IdY - 1})

	for _, v := range xy {
		if a.XYCompliance(v[0], v[1]) {
			gId := a.GetGIdFromXY(v[0], v[1])

			gridSurround = append(gridSurround, a.grids[gId])
		}
	}

	return gridSurround
}

func (a *AOIManager) XYCompliance(x, y int) bool {
	if x < 0 || x > a.CntsX-1 {
		return false
	}

	if y < 0 || y > a.CntsY-1 {
		return false
	}

	return true
}

func (a *AOIManager) GetGidByPos(xPos, yPos float32) int {
	x := (int(xPos) - a.MinX) / a.gridWidth()
	y := (int(yPos) - a.MinY) / a.gridLength()

	return a.GetGIdFromXY(x, y)
}

func (a *AOIManager) GetPidsByPos(xPos, yPos float32) []int32 {
	gId := a.GetGidByPos(xPos, yPos)

	grids := a.GetSurroundGridsByGId(gId)

	playerIDs := make([]int32, 0)

	for _, v := range grids {
		playerIDs = append(playerIDs, v.GetPlayerIds()...)

	}

	fmt.Println("players Id : --> ", playerIDs)
	return playerIDs
}

func (a *AOIManager) AddPidToGrid(pId int32, gId int) {
	grid, ok := a.grids[gId]
	if !ok {
		return
	}

	grid.AddPlayer(pId)
}

func (a *AOIManager) RemovePidFromGrid(pId int32, gId int) {
	grid, ok := a.grids[gId]
	if !ok {
		return
	}

	grid.RemovePlayer(pId)
}

func (a *AOIManager) GetPidsByGid(gId int) []int32 {
	grid, ok := a.grids[gId]
	if !ok {
		return nil
	}

	return grid.GetPlayerIds()
}

func (a *AOIManager) AddToGridByPos(pId int32, xPos, yPos float32) {
	gid := a.GetGidByPos(xPos, yPos)

	a.AddPidToGrid(pId, gid)
}

func (a *AOIManager) RemoveFromGridByPos(pId int32, xPos, yPos float32) {
	gid := a.GetGidByPos(xPos, yPos)

	a.RemovePidFromGrid(pId, gid)
}
