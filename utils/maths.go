package utils

import (
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
	"time"
)

type uMaths struct {
}

// Maths 数据处理通用类
var Maths = uMaths{}

// SecondToHour 秒转小时
func (u *uMaths) SecondToHour(sec int) float64 {
	return u.NumberDecimal(float64(sec)/3600, 2, true)
}

// SecondToString 秒转时间
func (u *uMaths) SecondToString(sec int) string {
	arr := make([]string, 0)
	duration := time.Duration(int64(sec)) * time.Second
	hour := int(duration.Hours())
	arr = append(arr, fmt.Sprintf("%d", hour))
	minute := int(duration.Minutes()) % 60
	arr = append(arr, Strings.PadLeft(minute, "0", 2))
	second := int(duration.Seconds()) % 60
	arr = append(arr, Strings.PadLeft(second, "0", 2))
	return strings.Join(arr, ":")
}

// MoneyFen2Yuan 金额分转元
func (u *uMaths) MoneyFen2Yuan(fen interface{}) float64 {
	f := decimal.NewFromFloat(gconv.Float64(fen))
	return f.Div(decimal.NewFromFloat(100)).InexactFloat64()
}

// MoneyFen2YuanStr 金额分转元字符串
func (u *uMaths) MoneyFen2YuanStr(fen interface{}) string {
	f := decimal.NewFromFloat(gconv.Float64(fen))
	return f.Div(decimal.NewFromFloat(100)).StringFixedBank(2)
}

// MoneyYuan2Fen 金额元转分
func (u *uMaths) MoneyYuan2Fen(yuan float64) int64 {
	f := decimal.NewFromFloat(yuan)
	return f.Mul(decimal.NewFromFloat(100)).IntPart()
}

// NumberDecimal number_format()
// decimals: Sets the number of decimal points.
// decPoint: Sets the separator for the decimal point.
// thousandsSep: Sets the thousands' separator.
func (u *uMaths) NumberDecimal(number float64, decimals int32, isRound ...bool) float64 {
	a := decimal.NewFromFloat(number)
	if len(isRound) > 0 && isRound[0] {
		return a.Round(decimals).InexactFloat64()
	}
	return a.Truncate(decimals).InexactFloat64()
	// res := strconv.FormatFloat(math.Floor(number*d+rd)/d, 'f', -1, 64)
	// return gconv.Float64(res)
	// value, _ := strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(dec)+"F", number), 64)
}

// NumberFormat number_format()
// decimals: Sets the number of decimal points.
// decPoint: Sets the separator for the decimal point.
// thousandsSep: Sets the thousands' separator.
func (u *uMaths) NumberFormat(number float64, decimals uint, isRound bool, decPointAndThousandsSep ...string) string {
	neg := false
	if number < 0 {
		number = -number
		neg = true
	}
	decPoint := "."
	thousandsSep := ""
	if len(decPointAndThousandsSep) > 0 {
		decPoint = decPointAndThousandsSep[0]
	}
	if len(decPointAndThousandsSep) > 1 {
		thousandsSep = decPointAndThousandsSep[1]
	}
	dec := int(decimals)
	// Will round off
	str := fmt.Sprintf("%."+strconv.Itoa(dec)+"f", u.NumberDecimal(number, int32(dec), isRound))
	prefix, suffix := "", ""
	if dec > 0 {
		prefix = str[:len(str)-(dec+1)]
		suffix = str[len(str)-dec:]
	} else {
		prefix = str
	}
	sep := []byte(thousandsSep)
	n, l1, l2 := 0, len(prefix), len(sep)
	// thousands sep num
	c := (l1 - 1) / 3
	tmp := make([]byte, l2*c+l1)
	pos := len(tmp) - 1
	for i := l1 - 1; i >= 0; i, n, pos = i-1, n+1, pos-1 {
		if l2 > 0 && n > 0 && n%3 == 0 {
			for j := range sep {
				tmp[pos] = sep[l2-j-1]
				pos--
			}
		}
		tmp[pos] = prefix[i]
	}
	s := string(tmp)
	if dec > 0 {
		s += decPoint + suffix
	}
	if neg {
		s = "-" + s
	}

	return s
}
