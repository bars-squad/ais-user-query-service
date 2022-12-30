package admin

import (
	"context"
	"encoding/json"

	"github.com/bars-squad/ais-user-query-service/entity"
	"github.com/bars-squad/ais-user-query-service/jwt"
	"github.com/bars-squad/ais-user-query-service/model"
	"github.com/bars-squad/ais-user-query-service/responses"
	"github.com/sirupsen/logrus"
)

var (
	httpResponses = responses.HttpResponseStatusCodesImpl{}
)

type Usecase interface {
	GetListAccount(ctx context.Context, size int, actionQS string, cursorQS string) responses.Responses
	OnCreatedAdministrator(ctx context.Context, user model.AdminRegistration) responses.Responses
}

type UsecaseImpl struct {
	ServiceName string
	Logger      *logrus.Logger
	Repository  Repository
	// SubsidyProductRepository subsidyproduct.Repository
	JSONWebToken *jwt.JSONWebToken
	// Session                  session.Session
	// Publisher                pubsub.Publisher
}

func NewUsecase(property *Property) Usecase {
	return &UsecaseImpl{
		ServiceName:  property.ServiceName,
		Logger:       property.Logger,
		Repository:   property.Repository,
		JSONWebToken: property.JSONWebToken,
		// sess:                     property.Session,
		// publisher:                property.Publisher,
	}
}

func (u UsecaseImpl) GetListAccount(ctx context.Context, size int, actionQS string, cursorQS string) responses.Responses {
	return httpResponses.Ok("").NewResponses(nil, "")
}

func (u UsecaseImpl) OnCreatedAdministrator(ctx context.Context, payload model.AdminRegistration) responses.Responses {
	var user entity.Admin

	createdTpUserBuff, _ := json.Marshal(payload)
	if err := json.Unmarshal(createdTpUserBuff, &user); err != nil {
		u.Logger.Error(err)
		return httpResponses.InternalServerError("").NewResponses(nil, err.Error())
	}

	user.ID = payload.ID

	if err := u.Repository.Save(ctx, user); err != nil {
		u.Logger.WithField("payload", payload).Error(err)
		return httpResponses.InternalServerError("").NewResponses(nil, err.Error())
	}

	return httpResponses.Ok("").NewResponses(nil, "")

}
