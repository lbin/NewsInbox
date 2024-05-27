package paper

import (
	"fmt"
	"learnerai/golib/app/logger"
	"learnerai/models/mgodb/paper"
	"learnerai/req"

	"github.com/gin-gonic/gin"
)

var fieldCanReset = map[string]bool{
	"abstract_ch": true,
	"keywords":    true,
	"keywords_ch": true,
	"title_ch":    true,
}

func UnsetDailyInfoField(ctx *gin.Context, request req.UnsetDailyInfoField) error {
	if len(request.Fields) == 0 {
		logger.Debugf(ctx, "no field specified in request: %+v", request)
		return nil
	}

	for _, f := range request.Fields {
		if !fieldCanReset[f] {
			return fmt.Errorf("not allowed to unset field: %v", f)
		}
	}

	return paper.UnsetDailyInfoField(ctx, request.Id, request.Fields)
}
