package jwt

import (
	"fmt"
	"context"
	"encoding/base64"
	"crypto/x509"
	"crypto/rsa"
    "encoding/pem"
	"github.com/golang-jwt/jwt/v4"

	"github.com/rs/zerolog/log"

	"github.com/lambda-go-authorizer-cert/pkg/observability"
	"github.com/lambda-go-authorizer-cert/internal/domain/erro"
	"github.com/lambda-go-authorizer-cert/internal/domain/model"

	repository "github.com/lambda-go-authorizer-cert/pkg/database/dynamo"
)

var childLogger = log.With().Str("useCase", "jwt").Logger()

type UseCaseJwt struct{
	repository	*repository.Repository
	JwtKey		*string
}

func NewUseCaseJwt(	repository *repository.Repository, 
					jwtKey *string) *UseCaseJwt{
	childLogger.Debug().Msg("NewUseCaseJwt")

	return &UseCaseJwt{
		repository: repository,
		JwtKey: jwtKey,
	}
}

func (u UseCaseJwt) TokenValidationRSA(ctx context.Context, bearerToken string) (bool, error){
	childLogger.Debug().Msg("TokenValidationRSA")

	span := observability.Span(ctx, "useCase.TokenValidationRSA")	
    defer span.End()

	log.Debug().Interface("bearerToken : ", bearerToken).Msg("")

	jwksDataRSABytes, err := base64.RawStdEncoding.DecodeString(bearerToken)
	if err != nil {
		childLogger.Error().Err(err).Msg("erro RawURLEncoding.DecodeString")
		return false, nil
	}

	var publicKey *rsa.PublicKey
	block, _ := pem.Decode([]byte(jwksDataRSABytes))
	if block == nil || block.Type != "PUBLIC KEY" {
		childLogger.Error().Err(erro.ErrDecodeKey).Msg("erro Decode")
		return false, nil
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		childLogger.Error().Err(err).Msg("erro ParsePKIXPublicKey")
		return false, nil
	}
	publicKey = pubInterface.(*rsa.PublicKey)

	claims := &model.JwtData{}

	tkn, err := jwt.ParseWithClaims(bearerToken, claims, func(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("error unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		fmt.Println(err)
		if err == jwt.ErrSignatureInvalid {
			return false, erro.ErrStatusUnauthorized
		}
		return false, erro.ErrTokenExpired
	}

	if !tkn.Valid {
		return false, erro.ErrStatusUnauthorized
	}

	return true ,nil
}

func (u UseCaseJwt) TokenValidation(ctx context.Context, bearerToken string, jwtKey string) (bool, error){
	childLogger.Debug().Msg("TokenValidation")

	span := observability.Span(ctx, "useCase.TokenValidation")	
    defer span.End()

	log.Debug().Interface("bearerToken : ", bearerToken).Msg("")
	log.Debug().Interface("jwtKey : ", jwtKey).Msg("")

	claims := &model.JwtData{}
	tkn, err := jwt.ParseWithClaims(bearerToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false, erro.ErrStatusUnauthorized
		}
		return false, erro.ErrTokenExpired
	}

	if !tkn.Valid {
		return false, erro.ErrStatusUnauthorized
	}

	return true, nil
}