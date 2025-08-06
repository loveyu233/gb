package examples

import (
	"fmt"
	"github.com/loveyu233/gb"
	"github.com/xuri/excelize/v2"
	"testing"
	"time"
)

func TestExcelMath(t *testing.T) {
	for i := 0; i < 50; i++ {
		for j := 0; j < 50; j++ {
			position := gb.ExcelGetPosition(i, j)
			row, col, _ := gb.ExcelParsePosition(position)
			fmt.Printf("i:%d j:%d position:%s row:%d col:%d\n", i, j, position, row, col)
		}
	}
}

func TestExcelMapper(t *testing.T) {
	/*
		type User struct {
		    ID    int64  `excel:"ID"`        // 精确匹配列名
		    Name  string `excel:"A"`         // Excel列名匹配
		    Email string `excel:"2"`         // 列索引匹配
		    Age   int    `excel:"年龄"`       // 中文列名匹配
		}
	*/
	// 定义目标结构体
	//type User struct {
	//	ID       int64     `excel:"ID"`
	//	Name     string    `excel:"姓名"`
	//	Email    string    `excel:"邮箱"`
	//	Age      int       `excel:"年龄"`
	//	Salary   float64   `excel:"工资"`
	//	IsActive bool      `excel:"状态"`
	//	Birthday time.Time `excel:"生日"`
	//	Phone    *string   `excel:"电话"` // 指针类型，支持空值
	//}
	// 定义目标结构体
	type User struct {
		Name  string `excel:"姓名"`
		Phone string `excel:"电话"`
	}

	// 创建映射器
	mapper := gb.InitExcelMapper()

	// 执行映射
	var users []User
	err := mapper.MapToStructs("test.xlsx", &users)
	if err != nil {
		fmt.Printf("映射失败: %v\n", err)
		return
	}

	// 输出结果
	fmt.Printf("成功映射 %d 条记录\n", len(users))
	for i, user := range users {
		fmt.Printf("用户 %d: %+v\n", i+1, user)
	}

	// 检查错误
	if errors := mapper.GetErrors(); len(errors) > 0 {
		fmt.Printf("\n映射过程中发生 %d 个错误:\n", len(errors))
		for _, err := range errors {
			fmt.Printf("- %s\n", err.Error())
		}
	}
}

func TestExcelExport(t *testing.T) {
	// 定义数据结构
	type User struct {
		ID       int64     `excel:"ID,title:用户ID"`
		Name     string    `excel:"姓名"`
		Email    string    `excel:"邮箱"`
		Age      int       `excel:"年龄"`
		Salary   float64   `excel:"工资"`
		IsActive bool      `excel:"状态"`
		Birthday time.Time `excel:"生日"`
		Phone    *string   `excel:"电话"`
	}

	// 准备测试数据
	phone1 := "13800138001"
	users := []User{
		{
			ID:       1,
			Name:     "张三",
			Email:    "zhangsan@example.com",
			Age:      25,
			Salary:   8500.50,
			IsActive: true,
			Birthday: time.Date(1998, 5, 15, 0, 0, 0, 0, gb.ShangHaiTimeLocation),
			Phone:    &phone1,
		},
		{
			ID:       2,
			Name:     "李四",
			Email:    "lisi@example.com",
			Age:      30,
			Salary:   12000.00,
			IsActive: false,
			Birthday: time.Date(1993, 8, 20, 0, 0, 0, 0, gb.ShangHaiTimeLocation),
			Phone:    nil, // 空值测试
		},
		{
			ID:       3,
			Name:     "王五",
			Email:    "wangwu@example.com",
			Age:      28,
			Salary:   9800.75,
			IsActive: true,
			Birthday: time.Date(1995, 12, 3, 0, 0, 0, 0, gb.ShangHaiTimeLocation),
			Phone:    nil,
		},
	}

	// 创建导出器,添加自定义数据
	//exporter := gb.InitExcelExporter(
	//	gb.WithExcelExporterSheetName("用户数据"),
	//	gb.WithExcelExporterIncludeHeader(gb.PtrType(true)),
	//	gb.WithExcelExporterColumnWidths(map[string]float64{
	//		"姓名": 15.0,
	//		"邮箱": 25.0,
	//		"工资": 12.0,
	//	}),
	//	gb.WithExcelExporterHeaderStyle(&gb.HeaderStyle{
	//		Bold:            true,
	//		BackgroundColor: "4F81BD", // 蓝色
	//		FontColor:       "FFFFFF", // 白色
	//		FontSize:        13,
	//		Alignment:       "center",
	//	}),
	//	gb.WithExcelExporterDataStyle(&gb.DataStyle{
	//		FontSize:     11,
	//		Alignment:    "left",
	//		NumberFormat: "General",
	//		DateFormat:   "yyyy-mm-dd",
	//	}),
	//)

	// 使用默认数据
	exporter := gb.InitExcelExporter()
	// 导出到文件
	err := exporter.ExportToFile(users, "users_export.xlsx")
	if err != nil {
		fmt.Printf("导出失败: %v\n", err)
		return
	}

	// 获取统计信息
	rows, cols := exporter.GetStats()
	fmt.Printf("导出成功! 共导出 %d 行 %d 列数据到文件: users_export.xlsx\n", rows, cols)

	// 测试导出到现有文件的新工作表
	file := excelize.NewFile()
	err = exporter.ExportToSheet(users, file, "Sheet1")
	if err != nil {
		fmt.Printf("导出到工作表失败: %v\n", err)
		return
	}

	err = file.SaveAs("users_export_custom.xlsx")
	if err != nil {
		fmt.Printf("保存文件失败: %v\n", err)
		return
	}

	fmt.Println("自定义导出也成功完成!")
}
