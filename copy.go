package gb

import "github.com/jinzhu/copier"

/*
结构体标签:
	copier:"-"	在复制过程中明确忽略该字段。
	copier:"must"	强制复制该字段；如果未复制该字段，复印机将会崩溃或返回错误。
	copier:"nopanic"	复印机将返回错误而不是恐慌。
	copier:"override"	即使设置了，也会强制复制字段IgnoreEmpty。用于用空值覆盖现有值。
	FieldName	当结构之间的字段名称不匹配时，指定用于复制的自定义字段名称。
用法:
	from结构体的方法如果和to中的结构体字段名称一样会被调用并赋值到to的对应字段
	to结构体的方法如何和form中的结构体字段名称一样会被调用
*/

// Copy 把from的值赋值到to中
func Copy(from, to any) error {
	return copier.Copy(to, from)
}

// DeepCopy 嵌套型结构复制
func DeepCopy(from, to any) error {
	return copier.CopyWithOption(to, from, copier.Option{DeepCopy: true, IgnoreEmpty: true})
}
