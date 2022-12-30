package admin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/bars-squad/ais-user-query-service/exception"
	"github.com/bars-squad/ais-user-query-service/model"
	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
)

type CreatedAdministratorEventHandler struct {
	Logger   *logrus.Logger
	Validate *validator.Validate
	Usescase Usecase
}

func NewCreatedAdministratorEventHandler(logger *logrus.Logger, validate *validator.Validate, usecase Usecase) *CreatedAdministratorEventHandler {
	return &CreatedAdministratorEventHandler{logger, validate, usecase}
}

// Handle will process the message.
func (handler *CreatedAdministratorEventHandler) Handle(ctx context.Context, message interface{}) (err error) {
	msg, ok := message.(*sarama.ConsumerMessage)
	if !ok {
		handler.Logger.Error("Not a kafka message")
		return
	}

	var user model.AdminRegistration
	var value string

	if err = json.Unmarshal(msg.Value, &value); err != nil {
		handler.Logger.Error(err)
		return
	}

	if err = json.Unmarshal([]byte(value), &user); err != nil {
		handler.Logger.Error(err)
		return
	}

	if err = handler.validateMessage(user); err != nil {
		handler.Logger.WithField("payload", user).Error(err)
		return
	}

	err = exception.InternalError(handler.Usescase.OnCreatedAdministrator(ctx, user))
	return
}

func (handler CreatedAdministratorEventHandler) validateMessage(message interface{}) (err error) {
	err = handler.Validate.Struct(message)
	if err == nil {
		return
	}

	errorFields := err.(validator.ValidationErrors)
	errorField := errorFields[0]
	err = fmt.Errorf("invalid '%s' with value '%v'", errorField.Field(), errorField.Value())

	return
}
