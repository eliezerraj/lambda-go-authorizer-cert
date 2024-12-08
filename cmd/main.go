package main

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/lambda-go-authorizer-cert/internal/domain/usecase/certs"
	"github.com/lambda-go-authorizer-cert/internal/domain/usecase/policy"
	"github.com/lambda-go-authorizer-cert/internal/domain/usecase/jwt"
	"github.com/lambda-go-authorizer-cert/internal/domain/model"
	"github.com/lambda-go-authorizer-cert/internal/config"

	"github.com/lambda-go-authorizer-cert/pkg/util"
	"github.com/lambda-go-authorizer-cert/pkg/aws_bucket_s3"
	"github.com/lambda-go-authorizer-cert/pkg/aws_parameter"

	"github.com/lambda-go-authorizer-cert/pkg/handler/apigw"
	"github.com/lambda-go-authorizer-cert/pkg/observability"

	repository "github.com/lambda-go-authorizer-cert/pkg/database/dynamo"

	"github.com/aws/aws-lambda-go/lambda"

	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
 	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/otel/trace"
)

var (
	logLevel = zerolog.DebugLevel // InfoLevel DebugLevel
	appServer	model.AppServer
	tracer 		trace.Tracer
)

func init(){
	log.Info().Msg("init")

	zerolog.SetGlobalLevel(logLevel)

	infoApp := util.GetAppInfo()
	configOTEL := util.GetOtelEnv()

	appServer.InfoApp = &infoApp
	appServer.ConfigOTEL = &configOTEL

	log.Info().Interface("appServer : ", appServer).Msg(".")
}

func main(){
	log.Info().Msg("main")
	
	ctx := context.Background()
	configAWS, err := config.GetAWSConfig(ctx, appServer.InfoApp.AWSRegion)
	if err != nil {
		panic("configuration error create new aws session " + err.Error())
	}

	//Load CRL
	clientS3 := aws_bucket_s3.NewClientS3Bucket(*configAWS)
	crl_pem, err := clientS3.GetObject(	ctx, 
										appServer.InfoApp.CrlBucketNameKey,
										appServer.InfoApp.CrlFilePath,
										appServer.InfoApp.CrlFileKey)
	if err != nil {
		log.Error().Err(err).Msg("Erro NewClientS3Bucket")
	}
	
	//Load symetric key
	clientSsm := aws_parameter.NewClientParameterStore(*configAWS)
	jwtKey, err := clientSsm.GetParameter(ctx, appServer.InfoApp.SSMJwtKey)
	if err != nil {
		panic("Error GetParameter " + err.Error())
	}

	// Create a repository
	repository, err := repository.NewRepository(ctx, configAWS, appServer.InfoApp.TableName)
	if err != nil {
		panic("configuration NewAuthRepository, " + err.Error())
	}

	tp := observability.NewTracerProvider(ctx, appServer.ConfigOTEL, appServer.InfoApp)
	defer func(ctx context.Context) {
			err := tp.Shutdown(ctx)
			if err != nil {
				log.Error().Err(err).Msg("Error shutting down tracer provider")
			}
	}(ctx)
	
	otel.SetTextMapPropagator(xray.Propagator{})
	otel.SetTracerProvider(tp)
	tracer = tp.Tracer("lambda-go-authorizer-cert")

	usecaseCerts := certs.NewUseCaseCerts(crl_pem)
	useCaseJwt := jwt.NewUseCaseJwt(repository, jwtKey)
	usecasePolicy := policy.NewUseCaseCPolicy()

	handler := apigw.InitializeLambdaHandler(*usecaseCerts, *usecasePolicy, *useCaseJwt)
	lambda.Start(otellambda.InstrumentHandler(handler.LambdaHandlerRequest, xrayconfig.WithRecommendedOptions(tp)... ))
	//lambda.Start(handler.LambdaHandlerRequest)
}