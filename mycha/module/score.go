package module

type CalculateScore func (counts Counts) uint32


func CalculateScoreSimple(counts Counts) uint32 {
	return counts.CallNum+counts.AcceptedNum+counts.CompletedNum+counts.handlingNum
}

//先计算一遍 当新分数是新的话 就设置组件的分数
func SetScore(module Module) bool {
	calculator := module.ScoreCalculator()
	if calculator == nil {
		calculator = CalculateScoreSimple
	}
	newScore := calculator(module.Counts())
	if newScore == module.Score() {
		return false
	}
	module.SetScore(newScore)
	return true
}