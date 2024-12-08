package policy

import(
	"context"

	"github.com/rs/zerolog/log"
	"github.com/lambda-go-authorizer-cert/pkg/observability"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lambda-go-authorizer-cert/internal/domain/model"
)

var childLogger = log.With().Str("useCase", "policy").Logger()

type UseCasePolicy struct{
}

func NewUseCaseCPolicy() *UseCasePolicy{
	childLogger.Debug().Msg("NewUseCaseCerts")
	return &UseCasePolicy{
	}
}

func(u *UseCasePolicy) GeneratePolicyFromClaims(ctx context.Context, policyData model.PolicyData) (events.APIGatewayCustomAuthorizerResponse){
	childLogger.Debug().Msg("GeneratePolicyFromClaims")
	
	span := observability.Span(ctx, "useCase.GeneratePolicyFromClaims")	
    defer span.End()

	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: policyData.PrincipalID}
	authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
		Version: "2012-10-17",
		Statement: []events.IAMPolicyStatement{
			{
				Action:   []string{"execute-api:Invoke"},
				Effect:   policyData.Effect,
				Resource: []string{policyData.MethodArn},
			},
		},
	}

	log.Debug().Interface("authResponse : ", authResponse).Msg("")

	return authResponse
}