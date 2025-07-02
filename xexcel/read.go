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
 @Time    : 2025/6/23 -- 17:56
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: xexcel xexcel/read.go
*/

package xexcel

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// xlsx.SetCellStr
// xlsx.SetCellValue

var (
	ExcelNullLine error = errors.New("empty col")
	NotFull       error = errors.New("missing required field")
)

type excelTool struct{}

var ExcelTool excelTool

// GetIOReaderFromUrl 通过url获取io Reader数据流
func (et excelTool) GetIOReaderFromUrl(url string) (reader io.ReadCloser, err error) {
	fun := "GetIOReaderFromUrl"
	load := func() (io.ReadCloser, error) {
		httpRes, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		// defer httpRes.Body.Close()

		if httpRes.StatusCode != 200 {
			_ = httpRes.Body.Close()
			return nil, fmt.Errorf("%s Http.Get not 200 error", fun)
		}
		return httpRes.Body, nil
	}
	for i := 1; i <= 3; i++ {
		if reader, err = load(); err == nil && reader != nil {
			break
		}
	}
	return
}

// GetExcelFileFromUrl 通过url获取并打开Excel文件
func (et excelTool) GetExcelFileFromUrl(url string) (f *excelize.File, err error) {
	fun := "GetExcelFileFromUrl"
	if url == "" {
		return nil, fmt.Errorf("%s url为空", fun)
	}
	ioReader, err := et.GetIOReaderFromUrl(url)
	if err != nil {
		return nil, err
	}
	defer ioReader.Close()
	f, err = excelize.OpenReader(ioReader)
	if err != nil {
		return nil, fmt.Errorf("%s open ExcelFile error,err:%s", fun, err)
	}
	return f, nil
}

// CreateNewExcel 创建一个新的Excel文件
func (et excelTool) CreateNewExcel() (f *excelize.File) {
	return excelize.NewFile()
}

// UploadExcelFile 上传Excel文件到CDN中,并获取URL
func (et excelTool) UploadExcelFile(ctx context.Context, f *excelize.File,
	uploader func(ctx context.Context, buf *bytes.Buffer, args ...any) (string, error),
	urlFmt func(ctx context.Context, srcUrl string) (dstUrl string),
	args ...any) (url string, err error) {

	writeBuf, err := f.WriteToBuffer()
	if err != nil {
		return "", err
	}

	url, err = uploader(ctx, writeBuf, args)
	if err != nil {
		return "", err
	}
	return urlFmt(ctx, url), nil
}

func (et excelTool) Processor(
	ctx context.Context,
	f *excelize.File,
	sheetName string,
	checkers []func([]string) (bool, error),                                                                     // 校验每行数据
	fallback func(ctx context.Context, row []string, index int, err error, errorNum *int64, f *excelize.File),   // 处理错误数据
	success func(ctx context.Context, row []string, index int, success *int64, f *excelize.File) (bool, string), // 处理正确数据
	cleaner func() error,                                                                                        // 清理现场：结果excel上传、结果上报、上传记录落库 etc.
) {
	// fun := "ExcelTool.Processor -->"
	var total, successNum, errorNum int64

	for i, row := range f.GetRows(sheetName) {
		if i < 1 { // 前1行为说明文字，不做处理
			continue
		}
		flag, reason := CheckRowData(ctx, row, checkers...)
		if !flag {
			fallback(ctx, row, i, reason, &errorNum, f)
		}
		total += 1
		success(ctx, row, i, &successNum, f)
	}

	if cleaner != nil {
		_ = cleaner()
	}
}

func CheckRowData(ctx context.Context, row []string, checkers ...func([]string) (bool, error)) (bool, error) {
	var flag bool
	var err error
	if flag, err = CheckIsNotNull(row); !flag {
		return false, err
	}
	for _, c := range checkers {
		if flag, err = c(row); !flag {
			return false, err
		}
	}
	return true, nil
}

// CheckIsNotNull 校验是否不是空行
func CheckIsNotNull(row []string) (bool, error) {
	for _, cell := range row {
		if cell != "" {
			return true, nil
		}
	}
	return false, ExcelNullLine
}
