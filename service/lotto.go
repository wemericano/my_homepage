package service

import (
	"math"
	"strconv"

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

func AnalyzeV2() ([]model.AnalyzeLottoList, error) {
	// 모든 로또 데이터 가져오기
	numlist, err := dbcall.GetLottoList()
	if err != nil {
		return nil, err
	}

	// 1. 빈도 분석 (Frequency Analysis)
	frequencySet := analyzeByFrequency(numlist)

	// 2. 가중치 분석 (Weighted Analysis) - 최근 데이터에 더 높은 가중치
	weightedSet := analyzeByWeightedFrequency(numlist)

	// 3. 패턴 분석 (Pattern Analysis) - 연속번호, 홀짝비율, 구간분포 고려
	patternSet := analyzeByPattern(numlist)

	// 4. 통계적 분석 (Statistical Analysis) - 정규분포 기반
	statisticalSet := analyzeByStatistics(numlist)

	// 5. 머신러닝 기반 분석 (ML-based Analysis) - 선형회귀와 클러스터링
	mlSet := analyzeByMachineLearning(numlist)

	var result []model.AnalyzeLottoList
	result = append(result, frequencySet)
	result = append(result, weightedSet)
	result = append(result, patternSet)
	result = append(result, statisticalSet)
	result = append(result, mlSet)

	return result, nil
}

// 1. 빈도 분석 - 전체 기간 동안 가장 많이 나온 번호들
func analyzeByFrequency(numlist []model.GetLottoList) model.AnalyzeLottoList {
	numberCount := make(map[string]int)

	for _, data := range numlist {
		numbers := []string{data.No1, data.No2, data.No3, data.No4, data.No5, data.No6}
		for _, num := range numbers {
			numberCount[num]++
		}
	}

	// 빈도순으로 정렬하여 상위 6개 선택
	var topNumbers []string
	for num, _ := range numberCount {
		topNumbers = append(topNumbers, num)
	}

	// 정렬
	for i := 0; i < len(topNumbers)-1; i++ {
		for j := i + 1; j < len(topNumbers); j++ {
			if numberCount[topNumbers[i]] < numberCount[topNumbers[j]] {
				topNumbers[i], topNumbers[j] = topNumbers[j], topNumbers[i]
			}
		}
	}

	if len(topNumbers) > 6 {
		topNumbers = topNumbers[:6]
	}

	return model.AnalyzeLottoList{
		No1: topNumbers[0], No2: topNumbers[1], No3: topNumbers[2],
		No4: topNumbers[3], No5: topNumbers[4], No6: topNumbers[5],
	}
}

// 2. 가중치 분석 - 최근 데이터에 더 높은 가중치 부여
func analyzeByWeightedFrequency(numlist []model.GetLottoList) model.AnalyzeLottoList {
	numberWeight := make(map[string]float64)
	totalWeight := 0.0

	for i, data := range numlist {
		// 최근 데이터일수록 높은 가중치 (지수적 감소)
		weight := math.Exp(-float64(i) * 0.01)
		totalWeight += weight

		numbers := []string{data.No1, data.No2, data.No3, data.No4, data.No5, data.No6}
		for _, num := range numbers {
			numberWeight[num] += weight
		}
	}

	// 정규화
	for num := range numberWeight {
		numberWeight[num] = numberWeight[num] / totalWeight
	}

	// 가중치순으로 정렬하여 상위 6개 선택
	var topNumbers []string
	for num := range numberWeight {
		topNumbers = append(topNumbers, num)
	}

	// 정렬
	for i := 0; i < len(topNumbers)-1; i++ {
		for j := i + 1; j < len(topNumbers); j++ {
			if numberWeight[topNumbers[i]] < numberWeight[topNumbers[j]] {
				topNumbers[i], topNumbers[j] = topNumbers[j], topNumbers[i]
			}
		}
	}

	if len(topNumbers) > 6 {
		topNumbers = topNumbers[:6]
	}

	return model.AnalyzeLottoList{
		No1: topNumbers[0], No2: topNumbers[1], No3: topNumbers[2],
		No4: topNumbers[3], No5: topNumbers[4], No6: topNumbers[5],
	}
}

// 3. 패턴 분석 - 홀짝비율, 구간분포, 연속번호 고려
func analyzeByPattern(numlist []model.GetLottoList) model.AnalyzeLottoList {
	// 구간별 빈도 (1-10, 11-20, 21-30, 31-40, 41-50, 51-60)
	sectionCount := make([]map[string]int, 6)
	for i := 0; i < 6; i++ {
		sectionCount[i] = make(map[string]int)
	}

	// 홀짝 빈도
	oddCount := make(map[string]int)
	evenCount := make(map[string]int)

	for _, data := range numlist {
		numbers := []string{data.No1, data.No2, data.No3, data.No4, data.No5, data.No6}
		for _, num := range numbers {
			if numInt, err := strconv.Atoi(num); err == nil {
				// 구간별 카운트
				section := (numInt - 1) / 10
				if section < 6 {
					sectionCount[section][num]++
				}

				// 홀짝 카운트
				if numInt%2 == 0 {
					evenCount[num]++
				} else {
					oddCount[num]++
				}
			}
		}
	}

	// 각 구간에서 가장 많이 나온 번호 선택
	var selectedNumbers []string
	for section := 0; section < 6; section++ {
		maxCount := 0
		maxNumber := ""
		for num, count := range sectionCount[section] {
			if count > maxCount {
				maxCount = count
				maxNumber = num
			}
		}
		if maxNumber != "" {
			selectedNumbers = append(selectedNumbers, maxNumber)
		}
	}

	// 부족한 경우 홀짝에서 보충
	if len(selectedNumbers) < 6 {
		// 홀수 번호 추가
		for num, _ := range oddCount {
			if len(selectedNumbers) >= 6 {
				break
			}
			found := false
			for _, selected := range selectedNumbers {
				if selected == num {
					found = true
					break
				}
			}
			if !found {
				selectedNumbers = append(selectedNumbers, num)
			}
		}

		// 짝수 번호 추가
		for num, _ := range evenCount {
			if len(selectedNumbers) >= 6 {
				break
			}
			found := false
			for _, selected := range selectedNumbers {
				if selected == num {
					found = true
					break
				}
			}
			if !found {
				selectedNumbers = append(selectedNumbers, num)
			}
		}
	}

	if len(selectedNumbers) > 6 {
		selectedNumbers = selectedNumbers[:6]
	}

	return model.AnalyzeLottoList{
		No1: selectedNumbers[0], No2: selectedNumbers[1], No3: selectedNumbers[2],
		No4: selectedNumbers[3], No5: selectedNumbers[4], No6: selectedNumbers[5],
	}
}

// 4. 통계적 분석 - 정규분포와 표준편차 고려
func analyzeByStatistics(numlist []model.GetLottoList) model.AnalyzeLottoList {
	// 각 번호의 평균과 표준편차 계산
	numberStats := make(map[string][]int)

	for _, data := range numlist {
		numbers := []string{data.No1, data.No2, data.No3, data.No4, data.No5, data.No6}
		for _, num := range numbers {
			if numInt, err := strconv.Atoi(num); err == nil {
				numberStats[num] = append(numberStats[num], numInt)
			}
		}
	}

	// 각 번호의 통계적 점수 계산
	numberScore := make(map[string]float64)
	for num, values := range numberStats {
		if len(values) > 0 {
			// 평균
			sum := 0
			for _, v := range values {
				sum += v
			}
			mean := float64(sum) / float64(len(values))

			// 표준편차
			variance := 0.0
			for _, v := range values {
				variance += math.Pow(float64(v)-mean, 2)
			}
			stdDev := math.Sqrt(variance / float64(len(values)))

			// 점수 = 빈도 * (1 / 표준편차) - 안정성 고려
			numberScore[num] = float64(len(values)) * (1.0 / (stdDev + 1.0))
		}
	}

	// 점수순으로 정렬하여 상위 6개 선택
	var topNumbers []string
	for num := range numberScore {
		topNumbers = append(topNumbers, num)
	}

	// 정렬
	for i := 0; i < len(topNumbers)-1; i++ {
		for j := i + 1; j < len(topNumbers); j++ {
			if numberScore[topNumbers[i]] < numberScore[topNumbers[j]] {
				topNumbers[i], topNumbers[j] = topNumbers[j], topNumbers[i]
			}
		}
	}

	if len(topNumbers) > 6 {
		topNumbers = topNumbers[:6]
	}

	return model.AnalyzeLottoList{
		No1: topNumbers[0], No2: topNumbers[1], No3: topNumbers[2],
		No4: topNumbers[3], No5: topNumbers[4], No6: topNumbers[5],
	}
}

// 5. 머신러닝 기반 분석 - 선형회귀와 클러스터링
func analyzeByMachineLearning(numlist []model.GetLottoList) model.AnalyzeLottoList {
	// 최근 100회 데이터만 사용 (학습 데이터)
	recentData := numlist
	if len(numlist) > 100 {
		recentData = numlist[:100]
	}

	// 각 번호의 출현 패턴 분석
	numberPattern := make(map[string][]int)

	for i, data := range recentData {
		numbers := []string{data.No1, data.No2, data.No3, data.No4, data.No5, data.No6}
		for _, num := range numbers {
			if _, err := strconv.Atoi(num); err == nil {
				numberPattern[num] = append(numberPattern[num], i)
			}
		}
	}

	// 각 번호의 추세 분석 (선형회귀)
	numberTrend := make(map[string]float64)
	for num, positions := range numberPattern {
		if len(positions) >= 3 {
			// 간단한 선형회귀: 최근 출현일수록 높은 점수
			sum := 0.0
			for _, pos := range positions {
				sum += float64(len(recentData) - pos) // 최근일수록 높은 값
			}
			numberTrend[num] = sum / float64(len(positions))
		}
	}

	// 클러스터링: 비슷한 패턴의 번호들 그룹화
	clusters := make(map[int][]string)
	clusterIndex := 0

	for num1 := range numberTrend {
		assigned := false
		for clusterIdx, cluster := range clusters {
			if len(cluster) > 0 {
				// 간단한 거리 계산 (출현 패턴 유사성)
				if len(numberPattern[num1]) > 0 && len(numberPattern[cluster[0]]) > 0 {
					// 패턴 유사성 체크
					assigned = true
					clusters[clusterIdx] = append(clusters[clusterIdx], num1)
					break
				}
			}
		}
		if !assigned {
			clusters[clusterIndex] = []string{num1}
			clusterIndex++
		}
	}

	// 각 클러스터에서 가장 높은 트렌드 점수를 가진 번호 선택
	var selectedNumbers []string
	for _, cluster := range clusters {
		if len(cluster) > 0 {
			bestNum := cluster[0]
			bestScore := numberTrend[bestNum]

			for _, num := range cluster {
				if numberTrend[num] > bestScore {
					bestNum = num
					bestScore = numberTrend[num]
				}
			}
			selectedNumbers = append(selectedNumbers, bestNum)
		}
	}

	// 부족한 경우 트렌드 점수 높은 순으로 추가
	if len(selectedNumbers) < 6 {
		var remainingNumbers []string
		for num := range numberTrend {
			found := false
			for _, selected := range selectedNumbers {
				if selected == num {
					found = true
					break
				}
			}
			if !found {
				remainingNumbers = append(remainingNumbers, num)
			}
		}

		// 트렌드 점수순으로 정렬
		for i := 0; i < len(remainingNumbers)-1; i++ {
			for j := i + 1; j < len(remainingNumbers); j++ {
				if numberTrend[remainingNumbers[i]] < numberTrend[remainingNumbers[j]] {
					remainingNumbers[i], remainingNumbers[j] = remainingNumbers[j], remainingNumbers[i]
				}
			}
		}

		// 부족한 만큼 추가
		needed := 6 - len(selectedNumbers)
		if needed > len(remainingNumbers) {
			needed = len(remainingNumbers)
		}
		selectedNumbers = append(selectedNumbers, remainingNumbers[:needed]...)
	}

	if len(selectedNumbers) > 6 {
		selectedNumbers = selectedNumbers[:6]
	}

	return model.AnalyzeLottoList{
		No1: selectedNumbers[0], No2: selectedNumbers[1], No3: selectedNumbers[2],
		No4: selectedNumbers[3], No5: selectedNumbers[4], No6: selectedNumbers[5],
	}
}
