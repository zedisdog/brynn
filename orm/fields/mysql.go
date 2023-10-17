package fields

type IField interface {
	//Type 返回字段的类型名称
	Type() string

	//Length 返回char或者varchar的长度
	Length() int

	Primary() string

	Name() string
}
