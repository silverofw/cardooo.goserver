package battleRoomMgr

import (
	"cardooo/common"
	MainEvent "cardooo/enum"
	"fmt"
	"log"
	"time"
)

const DefaultRoomID = 1000

type BattleRoomMgr struct {
	nextRoomId int
	Rooms      map[int]*BattleRoom
}

type BattleRoom struct {
	Id      int
	Tick    int
	player1 *common.Client
	player2 *common.Client
}

// 初始化 BattleRoomMgr 並返回指標
func NewBattleRoomMgr() *BattleRoomMgr {
	mgr := &BattleRoomMgr{
		nextRoomId: 1000,
		Rooms:      make(map[int]*BattleRoom),
	}
	log.Println("[enterRoom]")
	return mgr
}

func (mgr *BattleRoomMgr) EnterRoom(player *common.Client) *BattleRoom {
	log.Println("[EnterRoom]")
	var room, exit = mgr.Rooms[DefaultRoomID]
	if !exit {
		room = mgr.CreateNewRoom()
	}

	if room.player1 == nil {
		log.Printf("[EnterRoom][%v enter room %v] player 1 ", player.Id, DefaultRoomID)
		room.player1 = player
	} else if room.player2 == nil {
		log.Printf("[EnterRoom][%v enter room %v] player 2 ", player.Id, DefaultRoomID)
		room.player2 = player
		//room.player1.SendToClient(1, MainEvent.CSC_JOIN_ROOM, fmt.Sprintf("%v", player.Id))
	} else {
		// 房間已滿，可以處理相應的邏輯，例如返回錯誤或創建新房間
		log.Println("[EnterRoom] 房間已滿，無法加入新的玩家")
	}

	if room.player1 != nil && room.player2 != nil {
		log.Printf("[EnterRoom][%v] room ready ", DefaultRoomID)
		room.player1.SendToClient(1, MainEvent.CSC_BATTLE_ROOM_START_GAME, "")
		room.player2.SendToClient(1, MainEvent.CSC_BATTLE_ROOM_START_GAME, "")

		go BattleRoomTick(room)
	}
	return room
}

func (mgr *BattleRoomMgr) CreateNewRoom() *BattleRoom {
	log.Println("[CreateNewRoom]")
	var newRoom = &BattleRoom{
		Id: mgr.nextRoomId,
	}
	mgr.Rooms[mgr.nextRoomId] = newRoom
	mgr.nextRoomId++
	return newRoom
}

func (mgr *BattleRoomMgr) LeaveRoom(c *common.Client) {
	for _, room := range mgr.Rooms {
		if room.player1 != nil && room.player1.Id == c.Id {
			log.Printf("[LeaveRoom][%v]", c.Id)
			room.player1 = nil
			break
		}

		if room.player2 != nil && room.player2.Id == c.Id {
			log.Printf("[LeaveRoom][%v]", c.Id)
			room.player2 = nil
			break
		}
	}
}

func BattleRoomTick(room *BattleRoom) {
	room.Tick = 0
	log.Printf("[BattleRoomTick][%v] ROOM START TICK", room.Id)
	for {
		// 檢查玩家是否存在
		if room.player1 == nil || room.player2 == nil {
			log.Printf("[BattleRoomTick][%v] ROOM STOP TICK", room.Id)
			break
		}

		room.Tick++
		room.player1.SendToClient(1, MainEvent.CSC_BATTLE_ROOM_TICK, fmt.Sprintf("%v", room.Tick))
		room.player2.SendToClient(1, MainEvent.CSC_BATTLE_ROOM_TICK, fmt.Sprintf("%v", room.Tick))

		if room.Tick%60 == 0 {
			log.Printf("[BattleRoomTick][%v]", room.Tick)
		}
		time.Sleep(time.Second / 30)
	}
}
