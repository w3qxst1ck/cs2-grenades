package main

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/w3qxst1ck/cs2-grenades/internal/data"
)

func (app *application) createClient() (*s3.Client, error) {
	// Создаем кастомный обработчик эндпоинтов, который для сервиса S3 и региона ru-central1 выдаст корректный URL
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID && region == "ru-1" {
			return aws.Endpoint{
				// PartitionID:   "yc",
				URL:           app.config.storageS3.URL,
				SigningRegion: app.config.storageS3.Region,
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})

	// Подгружаем конфигрурацию из ~/.aws/*
	cfg, err := s3config.LoadDefaultConfig(context.TODO(), s3config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		return nil, err
	}

	// Создаем клиента для доступа к хранилищу S3
	client := s3.NewFromConfig(cfg)

	return client, nil
}

func (app *application) uploadImage(image multipart.File, filename string) error {
	client, err := app.createClient()
	if err != nil {
		return err
	}

	uploader := manager.NewUploader(client)

	input := &s3.PutObjectInput{
		Bucket: aws.String("test-bucket-2"),
		Key:    aws.String(filename),
		Body:   image,
	}

	res, err := uploader.Upload(context.TODO(), input)

	if err != nil {
		return err
	}
	_ = res

	return nil
}

func (app *application) deleteImagesFromStorage(images []*data.Image) error {
	client, err := app.createClient()
	if err != nil {
		return err
	}

	for _, image := range images {
		input := &s3.DeleteObjectInput{
			Bucket: aws.String(app.config.storageS3.Bucket), // s3 bucket name
			Key:    aws.String(image.Name),                    // file name
		}
	
		_, err = client.DeleteObject(context.TODO(), input)
		if err != nil {
			return err
		}
	}

	return nil
}
