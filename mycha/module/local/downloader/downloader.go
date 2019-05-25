package downloader

import (
	"mycha/errors"
	"mycha/helper/log"
	"mycha/module"
	"mycha/module/stub"
	"net/http"
)

var logger = log.DLogger()

func New(mid module.MID,client *http.Client,
	  scoreCalculator module.CalculateScore) (module.Downloader,error) {
	moduleBase, err := stub.NewModuleInternal(mid, scoreCalculator)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.NewCrawlerError(errors.ERROR_TYPE_PARAMETER,"空的下载客户端")
	}
	return &myDownloader{
		ModuleInternal: moduleBase,
		httpClient:     *client,
	}, nil
}


type myDownloader struct {
	stub.ModuleInternal
	httpClient http.Client
}


func (downloader *myDownloader) Download(req *module.Request) (*module.Response,error) {
	downloader.ModuleInternal.IncrHandlingNumber()
	defer downloader.ModuleInternal.DecrHandlingNumber()
	downloader.ModuleInternal.IncrCalledCount()
	if req == nil {
		return nil, errors.NewCrawlerError(errors.ERROR_TYPE_PARAMETER,"dd")
	}
	httpReq := req.HTTPReq()
	if httpReq == nil {
		return nil, errors.NewCrawlerError(errors.ERROR_TYPE_PARAMETER,"DD")
	}
	downloader.ModuleInternal.IncrAcceptedCount()
	logger.Infof("Do the request (URL: %s, depth: %d)... \n", httpReq.URL, req.Depth())
	httpResp, err := downloader.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	downloader.ModuleInternal.IncrCompletedCount()
	return module.NewResponse(httpResp, req.Depth()), nil
}













