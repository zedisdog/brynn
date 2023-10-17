package orm

type IModel interface {
	//TableName 返回表名
	TableName() string

	//Fill 向模型中填充数据.
	//填充时校验数据类型,如果数据校验不过返回带有详细错误信息的error.
	Fill(data map[string]any) error
}
