package aws

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/ffajarpratama/boiler-api/config"
	custom_error "github.com/ffajarpratama/boiler-api/pkg/error"
	"github.com/ffajarpratama/boiler-api/pkg/util"
)

type IFaceAWS interface {
	UploadToS3WithSlug(ctx context.Context, files []*multipart.FileHeader, slug string) ([]*UploadResponse, error)
	// UploadToS3(file *multipart.FileHeader) (map[string]interface{}, error)
	// UploadFromBuffer(data []byte, filename string, mimetype string) (map[string]interface{}, error)
	// UploadFromFile(content io.Reader, filename string, mimetype string) (map[string]interface{}, error)
}

type AWS struct {
	cnf     config.AWSConfig
	session *session.Session
}

type UploadResponse struct {
	Path     string
	Size     int
	Mimetype string
	Name     string
	URL      string
	Slug     string
}

func NewAWSClient(cnf config.AWSConfig) (IFaceAWS, error) {
	creds := credentials.NewStaticCredentials(cnf.AccessKeyID, cnf.AccessKeySecret, "")
	_, err := creds.Get()
	if err != nil {
		log.Fatal("[aws-init-error] \n", err.Error())
		return nil, err
	}

	awsConf := aws.NewConfig().WithRegion(cnf.Region).WithCredentials(creds)
	session, err := session.NewSession(awsConf)
	if err != nil {
		log.Fatal("[aws-session-error] \n", err.Error())
		return nil, err
	}

	return &AWS{cnf: cnf, session: session}, nil
}

func (a *AWS) UploadToS3WithSlug(ctx context.Context, files []*multipart.FileHeader, slug string) ([]*UploadResponse, error) {
	var res = make([]*UploadResponse, 0)

	if len(files) == 0 {
		err := custom_error.SetCustomError(&custom_error.ErrorContext{
			HTTPCode: http.StatusBadRequest,
			Message:  "file tidak boleh kosong",
		})

		return nil, err
	}

	fileObjects := make([]s3manager.BatchUploadObject, 0, len(files))
	now := time.Now().Unix()

	for _, FileHeader := range files {
		media, err := FileHeader.Open()
		if err != nil {
			return nil, err
		}

		defer media.Close()

		ext := filepath.Ext(FileHeader.Filename)
		buff := bytes.NewBuffer(nil)

		_, err = io.Copy(buff, media)
		if err != nil {
			return nil, err
		}

		prop := util.SlugUpload(slug)
		path, err := prop.GetPath()
		if err != nil {
			return nil, err
		}

		filename := fmt.Sprintf("%d%s", now, ext)
		filepath := fmt.Sprintf("%s/%s", path, filename)
		mimetype := http.DetectContentType(buff.Bytes())

		fileObject := &s3manager.UploadInput{
			Key:          aws.String(filepath),
			Bucket:       aws.String(a.cnf.Bucket),
			Body:         bytes.NewReader(buff.Bytes()),
			ContentType:  aws.String(mimetype),
			CacheControl: aws.String(fmt.Sprintf("max-age=%d", util.IMAGE_UPLOAD_MAX_AGE)),
		}

		fileObjects = append(fileObjects, s3manager.BatchUploadObject{Object: fileObject})

		res = append(res, &UploadResponse{
			Path:     filepath,
			Size:     int(FileHeader.Size),
			Mimetype: mimetype,
			Name:     filename,
			URL:      fmt.Sprintf("%s/%s", a.cnf.URL, filepath),
			Slug:     slug,
		})
	}

	svc := s3manager.NewUploader(a.session)
	err := svc.UploadWithIterator(aws.BackgroundContext(), &s3manager.UploadObjectsIterator{Objects: fileObjects})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// // UploadToS3 implements IFaceAWS.
// func (a *AWS) UploadToS3(file *multipart.FileHeader) (map[string]interface{}, error) {
// 	f, err := file.Open()
// 	if err != nil {
// 		return nil, err
// 	}

// 	mimetype := file.Header.Get("Content-Type")
// 	extname := filepath.Ext(file.Filename)
// 	n := strings.LastIndexByte(file.Filename, '.')

// 	filename := util.RemoveSpecialCharacters(file.Filename[:n])
// 	filename = fmt.Sprintf("%s_%s%s", filename, strconv.FormatInt(util.MakeTimeStamp(), 10), extname)

// 	s3PutObject := &s3.PutObjectInput{
// 		Bucket:      aws.String(a.Conf.AWS.Bucket),
// 		Key:         aws.String(filename),
// 		Body:        f,
// 		ContentType: &mimetype,
// 		// ACL:    "public-read",
// 	}

// 	cfg, err := aws_config.LoadDefaultConfig(context.TODO(), aws_config.WithRegion(a.Conf.AWS.Region))
// 	if err != nil {
// 		log.Printf("error: %v", err)
// 		return nil, err
// 	}

// 	client := s3.NewFromConfig(cfg)
// 	uploader := manager.NewUploader(client)
// 	result, err := uploader.Upload(context.TODO(), s3PutObject)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return map[string]interface{}{
// 		"filename":  filename,
// 		"media_url": result.Location,
// 		"file_size": file.Size,
// 	}, nil
// }

// // UploadFromBuffer implements IFaceAWS.
// func (a *AWS) UploadFromBuffer(data []byte, filename string, mimetype string) (map[string]interface{}, error) {
// 	s3PutObject := &s3.PutObjectInput{
// 		Bucket: aws.String(a.Conf.AWS.Bucket),
// 		Key:    aws.String(filename),
// 		Body:   bytes.NewReader(data),
// 	}

// 	if mimetype != "" {
// 		s3PutObject.ContentType = &mimetype
// 	}

// 	cfg, err := aws_config.LoadDefaultConfig(context.TODO(), aws_config.WithRegion(a.Conf.AWS.Region))
// 	if err != nil {
// 		log.Printf("error: %v", err)
// 		return nil, err
// 	}

// 	client := s3.NewFromConfig(cfg)
// 	uploader := manager.NewUploader(client)
// 	result, err := uploader.Upload(context.TODO(), s3PutObject)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return map[string]interface{}{
// 		"filename":  filename,
// 		"media_url": result.Location,
// 	}, nil
// }

// // UploadFromFile implements IFaceAWS.
// func (a *AWS) UploadFromFile(content io.Reader, filename string, mimetype string) (map[string]interface{}, error) {
// 	s3PutObject := &s3.PutObjectInput{
// 		Bucket: aws.String(a.Conf.AWS.Bucket),
// 		Key:    aws.String(filename),
// 		Body:   content,
// 	}

// 	if mimetype != "" {
// 		s3PutObject.ContentType = &mimetype
// 	}

// 	cfg, err := aws_config.LoadDefaultConfig(context.TODO(), aws_config.WithRegion(a.Conf.AWS.Region))
// 	if err != nil {
// 		log.Printf("error: %v", err)
// 		return nil, err
// 	}

// 	client := s3.NewFromConfig(cfg)
// 	uploader := manager.NewUploader(client)
// 	result, err := uploader.Upload(context.TODO(), s3PutObject)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return map[string]interface{}{
// 		"filename":  filename,
// 		"media_url": result.Location,
// 	}, nil
// }
