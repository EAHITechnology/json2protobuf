package proto

type Kind string

func NewKind(kindName string) Kind {
	return Kind(kindName)
}

func (k Kind) String() string {
	return string(k)
}

const (
	Int64Kind  Kind = "int64"
	Int32Kind  Kind = "int32"
	FloatKind  Kind = "float"
	StringKind Kind = "string"
	BoolKind   Kind = "bool"
	ArrayKind  Kind = "array"
	MapKind    Kind = "map"
)

type Typ struct {
	Name        Kind `json:"name"`
	ElementType Kind `json:"elementType"`
	KeyType     Kind `json:"keyType"`
	ValueType   Kind `json:"valueType"`
}

type Field struct {
	Name string `json:"name"`
	Type Typ    `json:"type"`
}

type JsonSchema struct {
	MsgName string  `json:"msg_name"`
	Fields  []Field `json:"fields"`
}
