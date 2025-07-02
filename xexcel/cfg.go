/*
 *  ┏┓      ┏┓
 *┏━┛┻━━━━━━┛┻┓
 *┃　　　━　　  ┃
 *┃   ┳┛ ┗┳   ┃
 *┃           ┃
 *┃     ┻     ┃
 *┗━━━┓     ┏━┛
 *　　 ┃　　　┃神兽保佑
 *　　 ┃　　　┃代码无BUG！
 *　　 ┃　　　┗━━━┓
 *　　 ┃         ┣┓
 *　　 ┃         ┏┛
 *　　 ┗━┓┓┏━━┳┓┏┛
 *　　   ┃┫┫  ┃┫┫
 *      ┗┻┛　 ┗┻┛
 @Time    : 2025/6/24 -- 11:11
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: xexcel xexcel/cfg.go
*/

package xexcel

import excelizeV2 "github.com/xuri/excelize/v2"

// excel通用配置
var (
	FileTypeExcel = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	FileBussType  = "expend_doc"
)

var (
	ExcelFontStyle = &excelizeV2.Font{
		Family: "Microsoft YaHei",
		Size:   11,
		Color:  "777777",
	}

	ExcelAlignmentStyle = &excelizeV2.Alignment{
		Horizontal: "left",
		Vertical:   "bottom",
		WrapText:   true,
	}

	ExcelBorderStyle = []excelizeV2.Border{
		{
			Color: "DADEE0",
			Type:  "left",
			Style: 2,
		},
		{
			Color: "DADEE0",
			Type:  "top",
			Style: 2,
		},
		{
			Color: "DADEE0",
			Type:  "right",
			Style: 2,
		},
		{
			Color: "DADEE0",
			Type:  "bottom",
			Style: 2,
		},
	}

	ExcelFillStyle = excelizeV2.Fill{
		Type:    "pattern",
		Color:   []string{"EFEFEF"},
		Pattern: 1,
	}

	ExcelHeaderFillStyle = excelizeV2.Fill{
		Type:    "pattern",
		Color:   []string{"E6F4EA"},
		Pattern: 1,
	}

	ExcelHeaderFontStyle = &excelizeV2.Font{
		Family: "Microsoft YaHei",
		Size:   11,
		Bold:   true,
		Color:  "1f7f3b",
	}

	ExcelHeaderAlignmentStyle = &excelizeV2.Alignment{
		Horizontal:  "left",
		Indent:      1,
		ShrinkToFit: false,
		Vertical:    "center",
		WrapText:    true,
	}

	ExcelFirstColWidth = 16
	ExcelColWidth      = 25
	ExcelHeaderHeight  = 35
	ExcelHeight        = 30
)
