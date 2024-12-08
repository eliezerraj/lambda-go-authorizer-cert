package model

import(
	"crypto/rsa"
	
	"github.com/golang-jwt/jwt/v4"
)

type AppServer struct {
	InfoApp 		*InfoApp 		`json:"info_app"`
	ConfigOTEL		*ConfigOTEL		`json:"otel_config"`
}

type InfoApp struct {
	AppName				string `json:"app_name,omitempty"`
	AWSRegion			string `json:"aws_region,omitempty"`
	ApiVersion			string `json:"version,omitempty"`
	AvailabilityZone 	string `json:"availabilityZone,omitempty"`
	TableName			string `json:"table_name,omitempty"`
	JwtKey				string `json:"jwt_key,omitempty"`
	SSMJwtKey			string `json:"ssm_jwt_key,omitempty"`
	ScopeValidation		bool `json:"scope_validation"`
	CrlValidation		bool `json:"crl_validation"`
	CrlBucketNameKey	string `json:"crl_bucket_name_key"`
	CrlFilePath			string `json:"crl_file_path"`
	CrlFileKey			string `json:"crl_file_key"`
	Env					string `json:"env,omitempty"`
	AccountID			string `json:"account,omitempty"`
}

type ConfigOTEL struct {
	OtelExportEndpoint		string
	TimeInterval            int64    `mapstructure:"TimeInterval"`
	TimeAliveIncrementer    int64    `mapstructure:"RandomTimeAliveIncrementer"`
	TotalHeapSizeUpperBound int64    `mapstructure:"RandomTotalHeapSizeUpperBound"`
	ThreadsActiveUpperBound int64    `mapstructure:"RandomThreadsActiveUpperBound"`
	CpuUsageUpperBound      int64    `mapstructure:"RandomCpuUsageUpperBound"`
	SampleAppPorts          []string `mapstructure:"SampleAppPorts"`
}

type PolicyData struct {
	PrincipalID		string
	Effect			string
	MethodArn		string
	UsageIdentifierKey	*string		
	Message			string		
}

type JwtData struct {
	TokenUse	string 	`json:"token_use"`
	ISS			string 	`json:"iss"`
	Version		string 	`json:"version"`
	JwtId		string 	`json:"jwt_id"`
	Username	string 	`json:"username"`
	Scope	  	[]string `json:"scope"`
	jwt.RegisteredClaims
}


type RSA_Key struct{
	SecretNameH256		string 	`json:"secret_name_h256"`
	RSAPublicKey		string 	`json:"rsa_public_key"`
	RSAPublicKeyByte 	[]byte 	`json:"rsa_public_key_byte"`
	RSAPrivateKey		string 	`json:"rsa_private_key"`
	RSAPrivateKeyByte 	[]byte 	`json:"rsa_private_key_byte"`
	PrivateKeyPem		*rsa.PrivateKey
	HS256				[]byte 	`json:"h256_key_byte"`		
}