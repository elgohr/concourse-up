package iaas

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	// New versions of github.com/aws/aws-sdk-go/aws have these consts
	// but the version currently pinned by bosh-cli v2 does not

	// ErrCodeNoSuchBucket for service response error code
	// "NoSuchBucket".
	//
	// The specified bucket does not exist.
	awsErrCodeNoSuchBucket = "NoSuchBucket"

	// ErrCodeNoSuchKey for service response error code
	// "NoSuchKey".
	//
	// The specified key does not exist.
	awsErrCodeNoSuchKey = "NoSuchKey"

	// Returned when calling HEAD on non-existant bucket or object
	awsErrCodeNotFound = "NotFound"
)

func (client *AWSClient) DeleteVersionedBucket(name string) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(client.region),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
	})
	if err != nil {
		return err
	}

	s3Client := s3.New(sess, &aws.Config{Region: &client.region})

	// Delete all objects
	objects := []*s3.Object{}
	err = s3Client.ListObjectsPages(&s3.ListObjectsInput{Bucket: &name},
		func(output *s3.ListObjectsOutput, _ bool) bool {
			objects = append(objects, output.Contents...)

			return true
		})
	if err != nil {
		return err
	}

	for _, object := range objects {
		_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: &name,
			Key:    object.Key,
		})
		if err != nil {
			return nil
		}
	}

	time.Sleep(time.Second)

	_, err = s3Client.DeleteBucket(&s3.DeleteBucketInput{Bucket: &name})
	return err
}

func (client *AWSClient) EnsureBucketExists(name string) error {
	sess, err := session.NewSession(aws.NewConfig().WithCredentialsChainVerboseErrors(true))
	if err != nil {
		return err
	}

	s3Client := s3.New(sess, &aws.Config{Region: &client.region})

	_, err = s3Client.HeadBucket(&s3.HeadBucketInput{Bucket: &name})
	if err == nil {
		return nil
	}

	awsErrCode := err.(awserr.Error).Code()
	if awsErrCode != awsErrCodeNotFound && awsErrCode != awsErrCodeNoSuchBucket {
		return err
	}

	bucketInput := &s3.CreateBucketInput{
		Bucket: &name,
	}
	// NOTE the location constraint should only be set if using a bucket OTHER than us-east-1
	// http://docs.aws.amazon.com/AmazonS3/latest/API/RESTBucketPUT.html
	if client.region != "us-east-1" {
		bucketInput.CreateBucketConfiguration = &s3.CreateBucketConfiguration{
			LocationConstraint: &client.region,
		}
	}

	_, err = s3Client.CreateBucket(bucketInput)
	return err
}

func (client *AWSClient) WriteFile(bucket, path string, contents []byte) error {
	sess, err := session.NewSession(aws.NewConfig().WithCredentialsChainVerboseErrors(true))
	if err != nil {
		return err
	}
	s3Client := s3.New(sess, &aws.Config{Region: &client.region})

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &path,
		Body:   bytes.NewReader(contents),
	})
	return err
}

func (client *AWSClient) HasFile(bucket, path string) (bool, error) {
	sess, err := session.NewSession(aws.NewConfig().WithCredentialsChainVerboseErrors(true))
	if err != nil {
		return false, err
	}
	s3Client := s3.New(sess, &aws.Config{Region: &client.region})

	_, err = s3Client.HeadObject(&s3.HeadObjectInput{Bucket: &bucket, Key: &path})
	if err != nil {
		awsErrCode := err.(awserr.Error).Code()
		if awsErrCode == awsErrCodeNotFound || awsErrCode == awsErrCodeNoSuchKey {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (client *AWSClient) EnsureFileExists(bucket, path string, defaultContents []byte) ([]byte, bool, error) {
	sess, err := session.NewSession(aws.NewConfig().WithCredentialsChainVerboseErrors(true))
	if err != nil {
		return nil, false, err
	}

	s3Client := s3.New(sess, &aws.Config{Region: &client.region})

	output, err := s3Client.GetObject(&s3.GetObjectInput{Bucket: &bucket, Key: &path})
	if err == nil {
		var contents []byte
		contents, err = ioutil.ReadAll(output.Body)
		if err != nil {
			return nil, false, err
		}

		// Successfully loaded file
		return contents, true, nil
	}

	awsErrCode := err.(awserr.Error).Code()
	if awsErrCode != awsErrCodeNoSuchKey && awsErrCode != awsErrCodeNotFound {
		return nil, false, err
	}

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &path,
		Body:   bytes.NewReader(defaultContents),
	})
	if err != nil {
		return nil, false, err
	}

	// Created file from given contents
	return defaultContents, true, nil
}

func (client *AWSClient) LoadFile(bucket, path string) ([]byte, error) {
	sess, err := session.NewSession(aws.NewConfig().WithCredentialsChainVerboseErrors(true))
	if err != nil {
		return nil, err
	}

	s3Client := s3.New(sess, &aws.Config{Region: &client.region})

	output, err := s3Client.GetObject(&s3.GetObjectInput{Bucket: &bucket, Key: &path})
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(output.Body)
}

func (client *AWSClient) DeleteFile(bucket, path string) error {
	sess, err := session.NewSession(aws.NewConfig().WithCredentialsChainVerboseErrors(true))
	if err != nil {
		return err
	}

	s3Client := s3.New(sess, &aws.Config{Region: &client.region})
	_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &path,
	})

	return err
}
