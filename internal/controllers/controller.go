package controllers

var (
	controllerInstance *controller
)

func GetController() (instance *controller, err error) {
	if controllerInstance == nil {
		instance, err = newController()
		if err != nil {
			return nil, err
		}
		controllerInstance = instance
	}
	return controllerInstance, nil
}

type controller struct {
	Ws WsController
}

func newController() (instance *controller, err error) {
	ws, err := NewWsController()
	if err != nil {
		return
	}

	return &controller{
		Ws: ws,
	}, nil
}
