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
 @Time    : 2025/6/23 -- 18:41
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: xexcel xexcel/write.go
*/

package xexcel

import (
	"bytes"
	"context"

	excelizeV2 "github.com/xuri/excelize/v2"
)

func GenerateExcel() {

}

func GenerateXlsx(ctx context.Context, content map[string][][]interface{}, validation map[string][]*excelizeV2.DataValidation) (*bytes.Buffer, error) {

	var (
		file         = excelizeV2.NewFile()
		fileContent  *bytes.Buffer
		streamWriter *excelizeV2.StreamWriter
		err          error
	)

	_ = file.DeleteSheet("Sheet1")

	for sheet, sheetContent := range content {

		var endCol int

		_, _ = file.NewSheet(sheet)

		if streamWriter, err = file.NewStreamWriter(sheet); err != nil {
			return nil, err
		}
		for index, rowValue := range sheetContent {
			cell, _ := excelizeV2.CoordinatesToCellName(1, index+1)
			if err = streamWriter.SetRow(cell, rowValue); err != nil {
				return nil, err
			}
			if index == 0 {
				endCol = len(rowValue)
			}
		}
		if err = streamWriter.Flush(); err != nil {
			return nil, err
		}

		styleId, err := file.NewStyle(&excelizeV2.Style{
			Font:      ExcelFontStyle,
			Alignment: ExcelAlignmentStyle,
			Border:    ExcelBorderStyle,
		})
		if err != nil {
			return nil, err
		}
		fillStyleId, err := file.NewStyle(&excelizeV2.Style{
			Font:      ExcelFontStyle,
			Alignment: ExcelAlignmentStyle,
			Border:    ExcelBorderStyle,
			Fill:      ExcelFillStyle,
		})
		if err != nil {
			return nil, err
		}
		headerStyleId, err := file.NewStyle(&excelizeV2.Style{
			Font:      ExcelHeaderFontStyle,
			Alignment: ExcelHeaderAlignmentStyle,
			Border:    ExcelBorderStyle,
			Fill:      ExcelHeaderFillStyle,
		})
		if err != nil {
			return nil, err
		}
		if endCol > 0 {
			endColString, err := excelizeV2.ColumnNumberToName(endCol)
			endCell, _ := excelizeV2.CoordinatesToCellName(endCol, len(sheetContent))
			endHeaderCell, _ := excelizeV2.CoordinatesToCellName(endCol, 1)
			if err != nil {
				return nil, err
			}

			_ = file.SetColWidth(sheet, "A", endColString, float64(ExcelColWidth))
			if err = file.SetCellStyle(sheet, "A1", endCell, styleId); err != nil {
				return nil, err
			}
			for rowId := 1; rowId <= len(sheetContent); rowId += 2 {

				startCell, _ := excelizeV2.CoordinatesToCellName(1, rowId)
				endCell, _ := excelizeV2.CoordinatesToCellName(endCol, rowId)
				if err = file.SetCellStyle(sheet, startCell, endCell, fillStyleId); err != nil {
					return nil, err
				}

			}
			if err = file.SetCellStyle(sheet, "A1", endHeaderCell, headerStyleId); err != nil {
				return nil, err
			}

			_ = file.SetRowHeight(sheet, 1, float64(ExcelHeaderHeight))
		}

		// 添加校验
		if validation != nil {
			if sheetValidations, ok := validation[sheet]; ok {
				for _, sheetValidation := range sheetValidations {
					if sheetValidation != nil {
						_ = file.AddDataValidation(sheet, sheetValidation)
					}
				}
			}
		}
	}
	if fileContent, err = file.WriteToBuffer(); err != nil {
		return nil, err
	}
	return fileContent, nil
}
