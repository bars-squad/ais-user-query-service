package admin

import (
	"net/http"
	"strconv"

	"github.com/bars-squad/ais-user-query-service/middleware"
	"github.com/bars-squad/ais-user-query-service/responses"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	unprocessableEntityMessage = "Invalid payload format"
	internalServerErrorMessage = "Internal server error"
	badRequestMessage          = "Please check your payload"
)

var (
	httpResponse = responses.HttpResponseStatusCodesImpl{}
)

type HTTPHandler struct {
	Logger   *logrus.Logger
	Validate *validator.Validate
	Usecase  Usecase
}

func NewHTTPHandler(logger *logrus.Logger, validate *validator.Validate, router *mux.Router, basicAuth middleware.RouteMiddleware, usecase Usecase) {

	handler := &HTTPHandler{
		Logger:   logger,
		Validate: validate,
		Usecase:  usecase,
	}

	router.HandleFunc("/v1/admin/registration", basicAuth.Verify(handler.GetListAccount)).Methods(http.MethodGet)
	// router.HandleFunc("/mpv-general-registration/v1/users/registration/{nationalityId}", basicAuth.Verify(handler.GetUser)).Methods(http.MethodGet)
	// router.HandleFunc("/mpv-general-registration/v1/users/registration/{nationalityId}/subsidy-product/{subsidyProduct}", basicAuth.Verify(handler.GetUserByNationalityIDAndSubsidyProduct)).Methods(http.MethodGet)
}

func (handler *HTTPHandler) GetListAccount(w http.ResponseWriter, r *http.Request) {

	var resp responses.Responses
	ctx := r.Context()

	queryString := r.URL.Query()
	cursorQS := queryString.Get("cursor")
	actionQS := queryString.Get("action")
	sizeQS := queryString.Get("size")
	size, _ := strconv.Atoi(sizeQS)

	resp = handler.Usecase.GetListAccount(ctx, size, actionQS, cursorQS)
	responses.REST(w, resp)
}
