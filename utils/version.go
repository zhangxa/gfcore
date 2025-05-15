package utils

import (
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	"strconv"
	"strings"
)

type uVersion struct {
	version0 int // 等于
	version1 int // 大于
	version2 int // 小于
}

// Version 版本工具类
var Version = &uVersion{
	version0: 0,
	version1: 1,
	version2: -1,
}

// verFormat 去除额外空格、转换
func (u *uVersion) verFormat(v1str, v2str string) (v1, v2 string) {
	v1 = strings.ReplaceAll(strings.TrimSpace(strings.ToUpper(v1str)), "。", ".")
	v1 = strings.TrimPrefix(v1, "V")
	v1 = strings.TrimPrefix(v1, ".")
	v1 = strings.TrimSuffix(v1, ".")
	v2 = strings.ReplaceAll(strings.TrimSpace(strings.ToUpper(v2str)), "。", ".")
	v2 = strings.TrimPrefix(v2, "V")
	v2 = strings.TrimPrefix(v2, ".")
	v2 = strings.TrimSuffix(v2, ".")
	return
}
func (u *uVersion) compareSlice(v1slice, v2slice []string) int {
	for index := range v1slice {
		idxV1, _ := strconv.Atoi(v1slice[index])
		idxV2, _ := strconv.Atoi(v2slice[index])
		if idxV1 > idxV2 {
			return u.version1
		}
		if idxV1 < idxV2 {
			return u.version2
		}
		if len(v1slice)-1 == index {
			return u.version0
		}
	}
	return u.version0
}

func (u *uVersion) compareSlice1(v1slice, v2slice []string, flag int) int {
	for index := range v1slice {
		idxV1, _ := strconv.Atoi(v1slice[index])
		idxV2, _ := strconv.Atoi(v2slice[index])
		//按照正常逻辑v1slice 长度小
		if idxV1 > idxV2 {
			if flag == 2 {
				return u.version2
			}
			return u.version1

		}
		if idxV1 < idxV2 {
			if flag == 2 {
				return u.version1
			}
			return u.version2
		}
		if len(v1slice)-1 == index {
			if flag == 2 {
				return u.version1
			} else if flag == 1 {
				return u.version2
			}
		}
	}
	return u.version0
}

// CompareVersion 对比版本号
func (u *uVersion) CompareVersion(v1, v2 string) (res int) {
	s1, s2 := u.verFormat(v1, v2)
	v1Times := 0
	v2Times := 0
	v1Arr := strings.Split(s1, "+")
	if len(v1Arr) > 1 {
		v1Times = gconv.Int(v1Arr[1])
	}
	v2Arr := strings.Split(s2, "+")
	if len(v2Arr) > 1 {
		v2Times = gconv.Int(v2Arr[1])
	}
	v1slice := strings.Split(v1Arr[0], ".")
	v2slice := strings.Split(v2Arr[0], ".")
	//长度不相等直接退出
	if len(v1slice) != len(v2slice) {
		if len(v1slice) > len(v2slice) {
			res = u.compareSlice1(v2slice, v1slice, 2)
		} else {
			res = u.compareSlice1(v1slice, v2slice, 1)
		}
	} else {
		res = u.compareSlice(v1slice, v2slice)
	}
	if res == 0 {
		fmt.Sprintln(v1Times, v2Times)
		// 等于情况下，比对 次数 值
		if v1Times > v2Times {
			return u.version1
		} else if v1Times < v2Times {
			return u.version2
		}
	}
	return res

}
