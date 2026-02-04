package entry

import "log/slog"

type App struct {
	Config Config
	store  Store
	logger *slog.Logger
}

func NewApp(logger *slog.Logger, store Store) *App {
	return &App{
		store:  store,
		logger: logger,
		Config: Config{},
	}
}

type Store interface {
	CondominiumStore
	VisitStore
}

type Config struct{}

type Valid interface {
	Valid() error
}
