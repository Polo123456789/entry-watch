package app

type App struct {
	Config Config
}

func NewApp() *App {
	return &App{}
}

type Config struct{}

type Valid interface {
	Valid() error
}
