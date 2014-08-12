package plugins

type Plugin interface {
	GetColumns() []string
	Extract() []float64
}
