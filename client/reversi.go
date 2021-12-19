package client

import (
	"bufio"
	"context"
	"fmt"
	"github.com/xxarupakaxx/reversi/build"
	"github.com/xxarupakaxx/reversi/game"
	pb "github.com/xxarupakaxx/reversi/gen/pb/proto"
	"google.golang.org/grpc"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Reversi struct {
	sync.RWMutex
	started  bool
	finished bool
	me       *game.Player
	room     *game.Room
	game     *game.Game
}

func NewReversi() *Reversi {
	return &Reversi{}
}

func (r *Reversi) run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("faield to connect to grpc server:%w", err)
	}
	defer conn.Close()

	err = r.matching(ctx, pb.NewMatchingServiceClient(conn))
	if err != nil {
		return err
	}
	r.game = game.NewGame(r.me.Color)

	return r.play(ctx, pb.NewGameServiceClient(conn))

}

func (r *Reversi) matching(ctx context.Context, client pb.MatchingServiceClient) error {
	stream, err := client.JoinRoom(ctx, &pb.JoinRoomRequest{})
	if err != nil {
		return err
	}

	defer stream.CloseSend()

	fmt.Println("マッチング相手を探しております...")

	for true {
		resp, err := stream.Recv()
		if err != nil {
			return err
		}

		if resp.GetStatus() == pb.JoinRoomResponse_MATCHED {
			r.room = build.Room(resp.GetRoom())
			r.me = build.Player(resp.GetMe())
			fmt.Printf("Matched roomID=%d\n", resp.GetRoom().GetId())
			return nil
		} else if resp.GetStatus() == pb.JoinRoomResponse_WAITING {
			fmt.Println("waiting matching..")
		}
	}
	return nil
}

func (r *Reversi) play(ctx context.Context, client pb.GameServiceClient) error {
	c, cancel := context.WithCancel(ctx)
	defer cancel()
	stream, err := client.Play(c)
	if err != nil {
		return err
	}
	defer stream.CloseSend()

	go func() {
		err = r.send(c, stream)
		if err != nil {
			cancel()
		}
	}()

	err = r.receive(c, stream)
	if err != nil {
		cancel()
		return err
	}

	return nil
}

func (r *Reversi) send(ctx context.Context, stream pb.GameService_PlayClient) error {
	for true {
		r.RLock()

		if r.finished {
			r.RUnlock()
			return nil
		} else if !r.started {
			err := stream.Send(&pb.PlayerRequest{
				RoomId: int32(r.room.ID),
				Player: build.PBPlayer(r.me),
				Action: &pb.PlayerRequest_Start{Start: &pb.PlayerRequest_StartAction{}},
			})
			r.RUnlock()
			if err != nil {
				return err
			}

			for true {
				r.RLock()
				if r.started {
					r.RUnlock()
					fmt.Printf("対戦見つかったね")
					break
				}
				r.RUnlock()
				fmt.Println("対戦相手が見つかるまで待とうね")
				time.Sleep(time.Second * 1)

			}
		} else {
			r.RUnlock()
			fmt.Println("どの石を動かす？")
			stdin := bufio.NewScanner(os.Stdin)
			stdin.Scan()

			text := stdin.Text()
			x, y, err := parseInput(text)
			if err != nil {
				fmt.Println(err)
				continue
			}

			r.Lock()
			_, err = r.game.Move(x, y, r.me.Color)
			r.Unlock()
			if err != nil {
				fmt.Println(err)
				continue
			}
			go func() {
				err = stream.Send(&pb.PlayerRequest{
					RoomId: int32(r.room.ID),
					Player: build.PBPlayer(r.me),
					Action: &pb.PlayerRequest_Move{
						Move: &pb.PlayerRequest_MoveAction{
							Move: &pb.Move{
								X: int32(x),
								Y: int32(y),
							},
						},
					},
				})
				if err != nil {
					fmt.Println(err)
				}
			}()

			ch := make(chan int)
			go func(ch chan int) {
				fmt.Println("")
				for i := 0; i < 5; i++ {
					fmt.Printf("%d秒間止まります \n", 5-i)
					time.Sleep(1 * time.Second)
				}
				fmt.Println("")
				ch <- 0
			}(ch)
			<-ch

		}

		select {
		case <-ctx.Done():
			return nil
		default:

		}

	}
	return nil
}

func parseInput(text string) (int, int, error) {
	ss := strings.Split(text, "-")
	if len(ss) != 2 {
		return 0, 0, fmt.Errorf("入力が不正です 例:A-1")
	}
	xs := ss[0]
	xrs := []rune(strings.ToUpper(xs))
	x := int(xrs[0]-rune('A')) + 1
	if x < 1 || 8 < x {
		return 0, 0, fmt.Errorf("入力が不正です。例：A-1")
	}

	ys := ss[1]
	y, err := strconv.ParseInt(ys, 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("入力が不正です。例：A-1")
	}
	if y < 1 || 8 < y {
		return 0, 0, fmt.Errorf("入力が不正です。例：A-1")
	}

	return x, int(y), nil
}

func (r *Reversi) receive(ctx context.Context, stream pb.GameService_PlayClient) error {
	for true {
		res, err := stream.Recv()
		if err != nil {
			return err
		}

		r.Lock()
		switch res.GetEvent().(type) {
		case *pb.PlayerResponse_Waiting:

		case *pb.PlayerResponse_Ready:
			r.started = true
			r.game.Display(r.me.Color)
		case *pb.PlayerResponse_Move:
			color := build.Color(res.GetMove().GetPlayer().GetColor())
			if color != r.me.Color {
				move := res.GetMove().GetMove()
				_, err = r.game.Move(int(move.GetX()), int(move.GetY()), color)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Print("石をどこに動かしますか 例(A-1)")
			}
		case *pb.PlayerResponse_Finished:
			r.finished = true

			winner := build.Color(res.GetFinished().Winner)
			fmt.Println("")
			if winner == game.None {
				fmt.Println("Draw!")
			} else if winner == r.me.Color {
				fmt.Println("you win!")
			} else {
				fmt.Println("You Lose!")
			}

			r.Unlock()
			return nil

		}
		r.Unlock()

		select {
		case <-ctx.Done():
			return nil
		default:

		}
	}
	return nil
}
