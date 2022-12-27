package tinyurl

type App struct {
	store   UrlStore
	version string
}

type TinyRequestPayload struct {
	Url string
}

type TinyRequestResponse struct {
	Code string
}

func NewApp(version string, store UrlStore) App {
	return App{version: version, store: store}
}

type UrlStore interface {
	GetById(int) (string, error)
	Store(string) (int, error)
}
