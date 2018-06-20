package object

type ObjectType string

type Object interface {
	Type() ObjectType
	String() string
}
