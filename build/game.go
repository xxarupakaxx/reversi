package build

import (
	"fmt"
	"github.com/xxarupakaxx/reversi/game"
	pb "github.com/xxarupakaxx/reversi/gen/pb/proto"
)

func Room(r *pb.Room) *game.Room {
	return &game.Room{
		ID:    int(r.GetId()),
		Host:  Player(r.GetHost()),
		Guest: Player(r.GetGuest()),
	}
}

func Player(p *pb.Player) *game.Player {
	return &game.Player{
		ID:    int(p.GetId()),
		Color: Color(p.GetColor()),
	}
}

func Color(p pb.Color) game.Color {
	switch p {
	case pb.Color_BLACK:
		return game.Black
	case pb.Color_WHITE:
		return game.White
	case pb.Color_WALL:
		return game.Wall
	case pb.Color_EMPTY:
		return game.Empty
	}

	panic(fmt.Sprintf("unknown color = %v", p))
}
