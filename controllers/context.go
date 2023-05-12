package controllers

import (
	"context"
	"errors"

	m "github.com/aleeXpress/cerca/models"
)

type service struct{
	name string
}
var key string = "user"
var serv service = service{"service"}
func SetContextUser(ctx context.Context, user *m.User)context.Context  {
	return context.WithValue(ctx, key, user)
}
func GetUserByContext(cxt context.Context)( *m.User,error)  {
	value := cxt.Value(key)
	u,ok := value.(*m.User)
	if !ok {
		return nil, errors.New("casting went wrong ")
	}
	return u, nil
}

func SetContextService(ctx context.Context, service *m.Services)context.Context  {
	return context.WithValue(ctx, serv, service)
}
func GetServiceByContext(cxt context.Context)( *m.Services,error)  {
	value := cxt.Value(serv)
	s,ok := value.(*m.Services)
	if !ok {
		return nil, errors.New("casting went wrong ")
	}
	return s, nil
}
