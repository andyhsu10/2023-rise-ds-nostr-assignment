package middlewares

var (
	middlewareInstance *middleware
)

func GetMiddleware() (*middleware, error) {
	if middlewareInstance == nil {
		m, err := newMiddleware()
		if err != nil {
			return nil, err
		}
		middlewareInstance = m
	}
	return middlewareInstance, nil
}

type middleware struct {
	Cors CorsMiddleware
}

func newMiddleware() (*middleware, error) {
	cors, err := NewCorsMiddleware()
	if err != nil {
		return nil, err
	}

	return &middleware{
		Cors: cors,
	}, nil
}
