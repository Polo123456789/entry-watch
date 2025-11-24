package app

type App struct {
	Config Config
	store  Store
}

func NewApp(store Store) *App {
	return &App{store: store}
}

type Store interface {
	CondominiumStore
	VisitStore
}

type Config struct{}

type Valid interface {
	Valid() error
}
