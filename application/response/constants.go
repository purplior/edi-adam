package response

var (
	Status_Ok                  = 200
	Status_Created             = 201
	Status_InternalServerError = 500
	Status_BadRequest          = 400
	Status_Unauthorized        = 401
	Status_Forbidden           = 403
	Status_NotFound            = 404
	Status_NotAcceptable       = 406
	Status_Unprocessable       = 422

	Message_ErrorNormal = "internal server error"
)
