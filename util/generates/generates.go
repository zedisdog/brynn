package generates

type IntBase interface {
	~int | ~int8 | ~int16 | ~int64 | ~int32
}

type UintBase interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type AllIntBase interface {
	~int | ~int8 | ~int16 | ~int64 | ~int32 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type FloatBase interface {
	~float32 | ~float64
}
