package s3helper

import (
	"../config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

type S3helper struct {
	S3Info                config.S3Info
}

type S3Helpers struct {
	S3Info                config.S3Info
	awsSession            *session.Session
}

func (s3helper S3helper) init(awsConfig *aws.Config) (*S3Helpers, error) {

	s3Helpers := S3Helpers{}
	s3Helpers.S3Info = s3helper.S3Info
	awsSession, err := session.NewSession(awsConfig)
	s3Helpers.awsSession = awsSession

	return &s3Helpers, err
}

func (s3helpers S3Helpers) downloadFromS3(localFileDownloadPath string, awsConfig *aws.Config) (int64, error) {

	file, err := os.Create(localFileDownloadPath)
	if err != nil {
		return -1, err
	}

	defer file.Close()

	sess, _ := session.NewSession(awsConfig)

	downloader := s3manager.NewDownloader(sess)

	numBytes, downloadErr := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(s3helpers.S3Info.Bucket),
			Key:    aws.String(s3helpers.S3Info.Prefix),
		})

	if err := file.Close(); err != nil {
		return -1, err
	}

	if downloadErr != nil {
		return -1, downloadErr
	}

	return numBytes, nil
}

func (s3helpers S3Helpers) uploadFileToS3(absoluteLocalFilePath string, s3FullKey string) error {

	file, err := os.Open(absoluteLocalFilePath)
	if err != nil {
		return err
	}

	reader, writer := io.Pipe()

	go func() <-chan error {

		goRoutineOut := make(chan error)
		_, err := io.Copy(writer, file)

		if err != nil {
			goRoutineOut <- err
			return goRoutineOut
		}

		err = file.Close()

		if err != nil {
			goRoutineOut <- err
			return goRoutineOut
		}

		err = writer.Close()

		if err != nil {
			goRoutineOut <- err
			return goRoutineOut
		}

		goRoutineOut <- nil
		return goRoutineOut
	}()

	uploader := s3manager.NewUploader(s3helpers.awsSession)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   reader,
		Bucket: aws.String(s3helpers.S3Info.Bucket),
		Key:    aws.String(s3FullKey),
	})
	if err != nil {
		return err
	}

	log.WithField("result-location", result.Location).Info("Successfully uploaded to.")

	return nil
}

func (s3helpers S3Helpers) deleteFromS3(s3FullKey string) (*s3.DeleteObjectOutput, error) {

	svc := s3.New(s3helpers.awsSession)

	deletedObject, err := svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(s3helpers.S3Info.Bucket), Key: aws.String(s3FullKey)})

	if err != nil {
		return nil, err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(s3helpers.S3Info.Bucket),
		Key:    aws.String(s3FullKey),
	})

	if err != nil {
		return nil, err
	}

	log.WithField("S3FullKey", s3FullKey).Info("S3 FullKey was deleted.")

	return deletedObject, nil
}
