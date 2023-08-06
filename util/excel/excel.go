package excel

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Excel struct{}

func (e *Excel) GenUsingDefaultSheet(getDataFunc func() (data [][]any, err error)) (content []byte, err error) {
	data, err := getDataFunc()
	if err != nil {
		return
	}
	f := excelize.NewFile()
	defer f.Close()

	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	for index, item := range data {
		err = f.SetSheetRow(sheet, fmt.Sprintf("A%d", index+1), &item)
		if err != nil {
			return
		}
	}

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return
	}
	content, err = io.ReadAll(buffer)
	return
}

func (e *Excel) GetDataFromActiveSheet(content []byte, col int) (data [][]string, err error) {
	f, err := excelize.OpenReader(bytes.NewReader(content))
	if err != nil {
		return
	}
	defer f.Close()

	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	if err != nil {
		return
	}
	rows, err := f.GetRows(sheet)
	data = make([][]string, 0, len(rows))
	for index, _ := range rows {
		c := make([]string, 0, col)
		for i := 0; i < col; i++ {
			var v string
			v, err = f.GetCellValue(sheet, Cell(i+1, index+1))
			if err != nil {
				return
			}
			c = append(c, v)
		}
		data = append(data, c)
	}
	return
}

func Cell(col int, row int) string {
	return fmt.Sprintf("%s%d", ColIndexByNum(col), row)
}

func ColIndexByNum(num int) string {
	s := make([]string, 0, num/26)
	for num != 0 {
		tmp := num % 26
		num /= 26

		//此处略微关键，当为0时，其实是26，也就是Z，
		//而且当你将0调整为26后，需要从数字中去除26代表的这个数
		if tmp == 0 {
			tmp = 26
			num -= 1
		}
		s = append(s, fmt.Sprintf("%c", 'A'+tmp-1))
	}

	for i := 0; i < len(s)/2; i++ {
		temp := s[i]
		s[i] = s[len(s)-1-i]
		s[len(s)-1-i] = temp
	}

	return strings.Join(s, "")
}
