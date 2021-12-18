package game


type Room struct {
	ID    int
	Host  *Player
	Guest *Player
}