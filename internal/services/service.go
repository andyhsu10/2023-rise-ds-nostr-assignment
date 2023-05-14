package services

var (
	serviceInstance *service
)

func GetService() (instance *service, err error) {
	if serviceInstance == nil {
		instance, err = newService()
		if err != nil {
			return nil, err
		}
		serviceInstance = instance
	}
	return serviceInstance, nil
}

type service struct {
	Ws WsService
}

func newService() (instance *service, err error) {
	ws, err := NewWsService()
	if err != nil {
		return
	}

	return &service{
		Ws: ws,
	}, nil
}
