package authenticator

import (
	"context"
	"errors"
	"github.com/coreos/go-oidc/v3/oidc"

	"github.com/src/go-auth-bff/config"
	"golang.org/x/oauth2"
)

// Authenticator is used to authenticate our users
type Authenticator struct {
	*oidc.Provider //
	oauth2.Config
}

// New instantiates the *Authenticator
func New()(*Authenticator,error){
provider,err :=oidc.NewProvider(
	context.Background(),
	"https://"+config.EnvVariables.Auth0Domain+"/",
)
if err != nil{
	return nil,err
}
conf :=oauth2.Config{
	ClientID: config.EnvVariables.Auth0ClientID,
	ClientSecret: config.EnvVariables.Auth0ClientSecret,
	RedirectURL: config.EnvVariables.Auth0CallbackURL,
	Endpoint: provider.Endpoint(),
	Scopes: []string{oidc.ScopeOpenID,"profile","email"},
}
return &Authenticator{
	Provider: provider,
	Config: conf,
},nil

}

// verifyId token verifies that an *oath.Token is valid *oidc.IDToken
func (a *Authenticator)VerifyIDToken(ctx context.Context,token *oauth2.Token)(*oidc.IDToken,error){
	rawIDToken,ok :=token.Extra("id_token").(string)

	if !ok {
		return nil,errors.New("no id_token field in ouath2 token")
	}
	oidcConfig :=&oidc.Config{
		ClientID: a.ClientID,
	}
	return a.Verifier(oidcConfig).Verify(ctx,rawIDToken)
}