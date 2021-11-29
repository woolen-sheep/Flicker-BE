package i18n

type StatusDict map[int]string

var Status = map[string]StatusDict{
	"zh-cn": HTTPStatusZHCN,
}
