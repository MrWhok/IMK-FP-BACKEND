package impl

import (
	"context"
	"github.com/MrWhok/IMK-FP-BACKEND/client"
	"github.com/MrWhok/IMK-FP-BACKEND/common"
	"github.com/MrWhok/IMK-FP-BACKEND/model"
	"github.com/MrWhok/IMK-FP-BACKEND/service"
)

func NewHttpBinServiceImpl(httpBinClient *client.HttpBinClient) service.HttpBinService {
	return &httpBinServiceImpl{HttpBinClient: *httpBinClient}
}

type httpBinServiceImpl struct {
	client.HttpBinClient
}

func (h *httpBinServiceImpl) PostMethod(ctx context.Context) {
	httpBin := model.HttpBin{
		Name: "rizki",
	}
	var response map[string]interface{}
	h.HttpBinClient.PostMethod(ctx, &httpBin, &response)
	common.NewLogger().Info("log response service ", response)
}
