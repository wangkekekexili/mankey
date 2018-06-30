package object

type HashKey string

type HashKeyer interface {
	HashKey() HashKey
}
