package utils

import (
	"encoding/json"
	"sort"

	"github.com/agnivade/levenshtein"
)

// 用于判断两个种子间读集或写集之间是否一致
// 一致：true 不一致:false
func StringArraysEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	sort.Strings(a)
	sort.Strings(b)

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// 主要作用：string或byte类型转为[]byte时不添加额外字符
func MarshalInterfaceToBytes(data interface{}) ([]byte, error) {
	switch v := data.(type) {
	case []byte:
		// 如果是 []byte，直接返回
		return v, nil
	case string:
		// 如果是 string，直接转换为 []byte
		return []byte(v), nil

	default:
		// 使用 json.Marshal 进行序列化
		marshaledData, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		// 返回序列化结果（包含 JSON 符号）
		return marshaledData, nil
	}
}

// 获取两个数中的较大值
func Max[T int | float64](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// 计算两个字符串之间的相似度
func CalculateSimilarity(str1, str2 string) float64 {
	distance := levenshtein.ComputeDistance(str1, str2)
	maxLen := Max(len(str1), len(str2))
	if maxLen == 0 {
		return 1.0 // 如果两个字符串都是空的，认为它们是完全相似的
	}
	return 1 - float64(distance)/float64(maxLen)
}
