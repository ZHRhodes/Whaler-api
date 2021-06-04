package models

var (
	Consumer ModelChangeConsumer
)

type ModelChangeConsumer interface {
	ModelChanged(id string)
}
