package services

type WsService interface { }

type wsService struct { }

func NewWsService() (WsService, error) {
	return &wsService{}, nil
}
