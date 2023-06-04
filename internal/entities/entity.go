package entities

type Interface interface {
	GetID() int64
}

type Filter []func(Interface) bool
