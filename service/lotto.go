package service

import (
	dbcall "my-homepage/dbcall"
	model "my-homepage/struct"
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
