package dynamo

import(
	"context"
	"github.com/rs/zerolog/log"

	"github.com/lambda-go-authorizer-cert/pkg/observability"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var childLogger = log.With().Str("pkg", "database").Logger()

type Repository struct {
	client 		*dynamodb.Client
	tableName   *string
}

func NewRepository(ctx context.Context, configAWS *aws.Config, tableName string) (*Repository, error){
	childLogger.Debug().Msg("NewRepository")

	span := observability.Span(ctx, "repository.NewRepository")	
    defer span.End()

	client := dynamodb.NewFromConfig(*configAWS)

	return &Repository {
		client: client,
		tableName: aws.String(tableName),
	}, nil
}