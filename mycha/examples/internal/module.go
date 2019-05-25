package internal

import (
	"mycha/module"
	"mycha/module/local/analyzer"
	"mycha/module/local/downloader"
	"mycha/module/local/pipline"
)

var snGen = module.NewSNGenertor(1,0) //生产一个序列号

func GetDownloaders(number uint32) ([]module.Downloader,error) {
	downloaders := []module.Downloader{}
	if number == 0 {
		return downloaders,nil
	}
	for i:= uint32(0);i<number; i++ {
		mid,err := module.GenMID(module.TYPE_DOWNLOADER,snGen.Get(),nil)
		if err != nil {
			return downloaders,err
		}
		d,err := downloader.New(mid,genHTTPClient(),module.CalculateScoreSimple)
		if err != nil {
			return downloaders,err
		}
		downloaders = append(downloaders,d)
	}
	return downloaders,nil
}

// GetAnalyzers 用于获取分析器列表。
func GetAnalyzers(number uint8) ([]module.Analyzer, error) {
	analyzers := []module.Analyzer{}
	if number == 0 {
		return analyzers, nil
	}
	for i := uint8(0); i < number; i++ {
		mid, err := module.GenMID(
			module.TYPE_ANALYZER, snGen.Get(), nil)
		if err != nil {
			return analyzers, err
		}
		a, err := analyzer.New(mid, genResponseParsers(), module.CalculateScoreSimple)
		if err != nil {
			return analyzers, err
		}
		analyzers = append(analyzers, a)
	}
	return analyzers, nil
}

// GetPipelines 用于获取条目处理管道列表。
func GetPipelines(number uint8, dirPath string) ([]module.Pipeline, error) {
	pipelines := []module.Pipeline{}
	if number == 0 {
		return pipelines, nil
	}
	for i := uint8(0); i < number; i++ {
		mid, err := module.GenMID(
			module.TYPE_PIPELINE, snGen.Get(), nil)
		if err != nil {
			return pipelines, err
		}
		a, err := pipeline.New(mid, genItemProcessors(dirPath), module.CalculateScoreSimple)
		if err != nil {
			return pipelines, err
		}
		a.SetFailFast(true)
		pipelines = append(pipelines, a)
	}
	return pipelines, nil
}