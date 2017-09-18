package util

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// NewSession ...
func NewSession() *session.Session {
	region := "eu-west-1"
	awsConfig := &aws.Config{
		Region: &region,
	}
	sess, _ := session.NewSession(awsConfig)
	return sess
}
