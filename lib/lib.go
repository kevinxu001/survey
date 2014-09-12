package lib

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

//create md5 string
func StrToMD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	rs := hex.EncodeToString(h.Sum(nil))
	return rs
}

//根据传入的time.Time，返回对应的时间描述
// “后天”
// “明天”
// “今天”
// “昨天”
// “前天”
// "8月17日"
func TimeToDesc(t time.Time) string {
	today := time.Now()
	year, month, day := today.Date()
	today = time.Date(year, month, day, 0, 0, 0, 0, time.Local)

	ty, tm, td := t.Date()
	t = time.Date(ty, tm, td, 0, 0, 0, 0, time.Local)

	switch {
	case time.Date(year, month, day+2, 0, 0, 0, 0, time.Local).Equal(t):
		return "后天"
	case time.Date(year, month, day+1, 0, 0, 0, 0, time.Local).Equal(t):
		return "明天"
	case time.Date(year, month, day, 0, 0, 0, 0, time.Local).Equal(t):
		return "今天"
	case time.Date(year, month, day-1, 0, 0, 0, 0, time.Local).Equal(t):
		return "昨天"
	case time.Date(year, month, day-2, 0, 0, 0, 0, time.Local).Equal(t):
		return "前天"
	default:
		return t.Format("2006年1月2日")
	}
}

func StringsToJson(str string) string {
	rs := []rune(str)
	jsons := ""
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			jsons += string(r)
		} else {
			jsons += "\\u" + strconv.FormatInt(int64(rint), 16) // json
		}
	}

	return jsons
}

func CheckErr(err error) {
	fmt.Println(err)
}
