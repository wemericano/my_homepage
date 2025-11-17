package service

import (
	"math/rand/v2"
	dbcall "my-homepage/dbcall"
	model "my-homepage/struct"
	"sort"
	"strconv"
)

func GetLottoList() ([]model.GetLottoList, error) {
	res, err := dbcall.GetLottoList()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func AnalyzeV1() ([]model.AnalyzeLottoList, error) {
	// 모든 로또 데이터 가져오기
	numlist, err := dbcall.GetLottoList()
	if err != nil {
		return nil, err
	}

	// 각 자리수별 번호 카운트
	positionCounts := make([]map[string]int, 6)
	for i := 0; i < 6; i++ {
		positionCounts[i] = make(map[string]int)
	}

	// 데이터 분석
	for _, data := range numlist {
		numbers := []string{data.No1, data.No2, data.No3, data.No4, data.No5, data.No6}
		for pos, num := range numbers {
			if pos < 6 {
				positionCounts[pos][num]++
			}
		}
	}

	// 각 자리수별로 가장 많이 나온 번호 1개씩 찾아서 하나의 로또 번호로 만들기
	var mostFrequentNumbers []string

	// 1번 자리부터 6번 자리까지 각각 가장 많이 나온 번호 찾기
	for pos := 0; pos < 6; pos++ {
		var maxCount = 0
		var maxNumber = ""

		// 해당 자리수에서 가장 많이 나온 번호 찾기
		for num, count := range positionCounts[pos] {
			if count > maxCount {
				maxCount = count
				maxNumber = num
			}
		}

		// 가장 많이 나온 번호 추가
		mostFrequentNumbers = append(mostFrequentNumbers, maxNumber)
	}

	// 하나의 로또 번호로 만들기
	analyzeItem := model.AnalyzeLottoList{
		No1: mostFrequentNumbers[0],
		No2: mostFrequentNumbers[1],
		No3: mostFrequentNumbers[2],
		No4: mostFrequentNumbers[3],
		No5: mostFrequentNumbers[4],
		No6: mostFrequentNumbers[5],
	}

	var result []model.AnalyzeLottoList
	result = append(result, analyzeItem)

	return result, nil
}

func AnalyzeV2() ([]model.AnalyzeLottoList, error) {
	numlist, err := dbcall.GetLottoList()
	if err != nil {
		return nil, err
	}

	if len(numlist) > 10 {
		numlist = numlist[len(numlist)-10:]
	}

	type Pattern struct {
		Consecutive int
		OddCount    int
		HighCount   int
	}

	var patterns []Pattern

	for _, data := range numlist {
		nums := []int{
			toInt(data.No1), toInt(data.No2), toInt(data.No3),
			toInt(data.No4), toInt(data.No5), toInt(data.No6),
		}
		sort.Ints(nums)

		consecutive := 0
		for i := 0; i < len(nums)-1; i++ {
			if nums[i]+1 == nums[i+1] {
				consecutive++
			}
		}

		odd := 0
		high := 0
		for _, n := range nums {
			if n%2 != 0 {
				odd++
			}
			if n >= 21 {
				high++
			}
		}

		patterns = append(patterns, Pattern{
			Consecutive: consecutive,
			OddCount:    odd,
			HighCount:   high,
		})
	}

	var totalConsec, totalOdd, totalHigh int
	for _, p := range patterns {
		totalConsec += p.Consecutive
		totalOdd += p.OddCount
		totalHigh += p.HighCount
	}
	avgConsec := totalConsec / len(patterns)
	avgOdd := totalOdd / len(patterns)
	avgHigh := totalHigh / len(patterns)

	var results []model.AnalyzeLottoList

	for len(results) < 5 {
		candidate := rand.Perm(45)[:6]
		for i := range candidate {
			candidate[i] += 1
		}
		sort.Ints(candidate)

		consec := 0
		odd := 0
		high := 0
		for i := 0; i < len(candidate)-1; i++ {
			if candidate[i]+1 == candidate[i+1] {
				consec++
			}
		}
		for _, n := range candidate {
			if n%2 != 0 {
				odd++
			}
			if n >= 21 {
				high++
			}
		}

		if abs(consec-avgConsec) <= 1 &&
			abs(odd-avgOdd) <= 1 &&
			abs(high-avgHigh) <= 1 {

			result := model.AnalyzeLottoList{
				No1: toStr(candidate[0]),
				No2: toStr(candidate[1]),
				No3: toStr(candidate[2]),
				No4: toStr(candidate[3]),
				No5: toStr(candidate[4]),
				No6: toStr(candidate[5]),
			}
			results = append(results, result)
		}
	}

	return results, nil
}

func toInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func toStr(n int) string {
	return strconv.Itoa(n)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
