package client

import (
	"context"
	"github.com/MrWhok/IMK-FP-BACKEND/model"
)

type HttpBinClient interface {
	PostMethod(ctx context.Context, requestBody *model.HttpBin, response *map[string]interface{})
}
