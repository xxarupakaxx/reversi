package handler

import (
	"fmt"
	"github.com/xxarupakaxx/reversi/build"
	"github.com/xxarupakaxx/reversi/game"
	pb "github.com/xxarupakaxx/reversi/gen/pb/proto"
	"sync"
)

type GameHandler struct {
	sync.RWMutex
	games  map[int]*game.Game
	client map[int][]pb.GameService_PlayServer
}

func (g *GameHandler) Play(stream pb.GameService_PlayServer) error {
	for true {
		req, err := stream.Recv()
		if err != nil {
			return err
		}

		roomID := req.GetRoomId()
		player := build.Player(req.GetPlayer())

		switch req.GetAction().(type) {
		case *pb.PlayerRequest_Start:
			err = g.start(stream, roomID, player)
			if err != nil {
				return err
			}
		case *pb.PlayerRequest_Move:
			action := req.GetMove()
			x := action.GetMove().GetX()
			y := action.GetMove().GetY()
			err = g.move(roomID, x, y, player)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *GameHandler) start(stream pb.GameService_PlayServer, id int32, player *game.Player) error {
	g.Lock()
	defer g.Unlock()

	ga := g.games[int(id)]
	if ga == nil {
		ga = game.NewGame(game.None)
		g.games[int(id)] = ga
		g.client[int(id)] = make([]pb.GameService_PlayServer, 0, 2)
	}

	g.client[int(id)] = append(g.client[int(id)], stream)

	if len(g.client[int(id)]) == 2 {
		for _, server := range g.client[int(id)] {
			err := server.Send(&pb.PlayerResponse{Event: &pb.PlayerResponse_Ready{
				Ready: &pb.PlayerResponse_ReadyEvent{},
			}})
			if err != nil {
				return err
			}

		}
		fmt.Printf("ゲームが始まりました roomID = %v\n", id)
	} else {
		err := stream.Send(&pb.PlayerResponse{
			Event: &pb.PlayerResponse_Waiting{
				Waiting: &pb.PlayerResponse_WaitingEvent{},
			}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *GameHandler) move(id int32, x int32, y int32, player *game.Player) error {
	g.Lock()
	defer g.Unlock()

	ga := g.games[int(id)]

	finished, err := ga.Move(int(x), int(y), player.Color)
	if err != nil {
		return err
	}

	for _, server := range g.client[int(id)] {
		err = server.Send(&pb.PlayerResponse{
			Event: &pb.PlayerResponse_Move{Move: &pb.PlayerResponse_MoveEvent{
				Player: build.PBPlayer(player),
				Move: &pb.Move{
					X: x,
					Y: y,
				},
				Board: build.PBBoard(ga.Board),
			}},
		})
		if err != nil {
			return err
		}

		if finished {
			err = server.Send(&pb.PlayerResponse{
				Event: &pb.PlayerResponse_Finished{
					Finished: &pb.PlayerResponse_FinishedEvent{
						Winner: build.PBColor(ga.Winner()),
						Board:  build.PBBoard(ga.Board),
					},
				},
			})
			delete(g.client,int(id))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func NewGameHandler() *GameHandler {
	return &GameHandler{
		games:  make(map[int]*game.Game),
		client: make(map[int][]pb.GameService_PlayServer),
	}
}
