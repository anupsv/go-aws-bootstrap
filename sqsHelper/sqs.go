package sqsHelper

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsSqs "github.com/aws/aws-sdk-go/service/sqs"
	"go-aws-bootstrap/config"
)

type Init struct {
	SqsInfo               config.SqsInfo
	localFileDownloadPath string
}

type SqsHelper struct {
	SqsInfo    config.SqsInfo
	awsSession *session.Session
	sqsSession *awsSqs.SQS
}

func (init Init) init(awsConfig *aws.Config) (*SqsHelper, error) {
	sqs := SqsHelper{}
	sqs.SqsInfo = init.SqsInfo
	awsSession, err := session.NewSession(awsConfig)
	sqs.awsSession = awsSession
	sqs.sqsSession = awsSqs.New(awsSession)
	return &sqs, err
}

func (sqs SqsHelper) sendMessage(awsSendMessageInput *awsSqs.SendMessageInput) (*awsSqs.SendMessageOutput, error) {
	awsSendMessageInput.SetQueueUrl(sqs.SqsInfo.SqsUrl)
	return sqs.sqsSession.SendMessage(awsSendMessageInput)
}

func (sqs SqsHelper) longPollingReceiveMessage(receiveTimeOut int64, awsReceiveMessageInput *awsSqs.ReceiveMessageInput) (*awsSqs.ReceiveMessageOutput, error) {
	awsReceiveMessageInput.SetWaitTimeSeconds(receiveTimeOut)
	return sqs.sqsSession.ReceiveMessage(awsReceiveMessageInput)
}
