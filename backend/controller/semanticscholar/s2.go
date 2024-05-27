package semanticscholar

import (
	"learnerai/golib/app/request"
	s2 "learnerai/logic/semanticscholar"
)

var GetPaperInfo = request.JsonController(s2.GetPaperInfo)
var GetPaperReference = request.JsonController(s2.GetPaperReference)
