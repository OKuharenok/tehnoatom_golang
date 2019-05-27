package main

import (
	"sort"
	"strconv"
)

func IntSliceToString(sl []int) string {
	res := ""
	for _, val := range sl {
		res += strconv.Itoa(val)
	}
	return res
}

func MergeSlices(sl1 []float32, sl2 []int32) []int {
	var newSlice []int
	for _, val := range sl1 {
		newSlice = append(newSlice, int(val))
	}
	for _, val := range sl2 {
		newSlice = append(newSlice, int(val))
	}
	return newSlice
}

func GetMapValuesSortedByKey (input map[int]string) []string {
	var res []string
	var k  []int
	for key := range input {
		k = append(k, key)
	}
	sort.Ints(k)
	for _, key := range k {
		res = append(res, input[key])
	}
	return res
}
