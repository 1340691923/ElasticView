package my_error

const (
	IndexNameNullError = 200001
	AliasNameNullError = 200002
)

var ParmasNullError = map[int]string{
	IndexNameNullError: "索引名不能为空",
	AliasNameNullError: "别名不能为空",
}
