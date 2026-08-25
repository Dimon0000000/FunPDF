package engine

type Engine interface {
	Translate(text, from, to string) (string, error)
	Healthy() bool
}
