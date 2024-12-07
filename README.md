# lambda-go-authorizer-cert

lambda-go-authorizer-cert

## Lambda Env Variables

      APP_NAME: lambda-go-authorizer-cert
      VERSION:  1.0
      CRL_BUCKET_NAME_KEY:  eliezerraj-908671954593-mtls-truststore
      CRL_FILE_KEY: crl_ca.pem
      OTEL_EXPORTER_OTLP_ENDPOINT:  localhost:4317
      CRL_FILE_PATH:    /
      CRL_VALIDATION:   false
      TABLE_NAME:   user_login_2
      
      JWT_KEY	            my_secret_key
      SCOPE_VALIDATION	    true
      SSM_JWT_KEY	         key-secret

## Test Locally

1 Download

    mkdir -p .aws-lambda-rie && curl -Lo .aws-lambda-rie/aws-lambda-rie https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie && chmod +x .aws-lambda-rie/aws-lambda-rie

2 Run

    /local-test$ ./start.sh

3 Invoke

    curl -X POST http://localhost:9000/2015-03-31/functions/function/invocations -d '{"headers":{"authorization":"teste"},"methodArn":"arn:aws:execute-api:us-east-2:908671954593:k0ng1bdik7/qa/GET/account/info"}'

## Compile lambda

   Manually compile the function

      New Version
      GOARCH=amd64 GOOS=linux go build -o ../build/bootstrap main.go
      zip -jrm ../build/main.zip ../build/bootstrap

        aws lambda update-function-code \
        --region us-east-2 \
        --function-name lambda-go-authorizer-cert \
        --zip-file fileb:///mnt/c/Eliezer/workspace/github.com/lambda-go-authorizer-cert/build/main.zip \
        --publish