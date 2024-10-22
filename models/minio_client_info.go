package models

type MinioClientInfo struct {
	AccessKey string
	SecretKey string
	Host      string
	Port      string
	Secure    bool
	Bucket    string
}
