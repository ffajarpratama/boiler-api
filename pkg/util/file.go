package util

import (
	"errors"
	"net/http"
	"strings"

	custom_error "github.com/ffajarpratama/boiler-api/pkg/error"
)

//  TODO: do properly

type SlugUpload string

const (
	SlugPortalCommunityPostImage    = "image"
	SlugPortalCommunityPostVideo    = "video"
	SlugPortalCommunityPostDocument = "document"
	SlugPortalCommunityThumbnail    = "thumbnail"
	SlugPortalAvatarImage           = "avatar"
	SlugNewsAndPromotionPicture     = "picture"

	// SlugOrderPoFile      = "po"
	// SlugOrderKtpFile     = "order-ktp"
	// SlugFaceRecognition  = "face-recognition"
	// SlugOCRKtp           = "ocr-ktp"
	// SlugDealerDocument   = "dealer-document"
	// SlugVehicle          = "vehicle"
	// SlugVehicleBrandLogo = "vehicle-brand-logo"

	ExtSlugPortalCommunityPostDocument = ".pdf"
	ExtSlugPortalCommunityPostImage    = ".jpeg,.png,.jpg"
	ExtSlugPortalCommunityThumbnail    = ".jpeg,.png,.jpg"

	IMAGE_UPLOAD_MAX_AGE = 365 * 24 * 60 * 60 // 1 year
)

var (
	UploadSlug = map[string]bool{
		SlugPortalCommunityPostImage:    true,
		SlugPortalCommunityPostVideo:    true,
		SlugPortalCommunityPostDocument: true,
		SlugPortalCommunityThumbnail:    true,
		SlugPortalAvatarImage:           true,
		SlugNewsAndPromotionPicture:     true,
		// SlugOrderPoFile:                 true,
		// SlugOrderKtpFile:                true,

		// SlugDealerDocument:              true,
		// SlugFaceRecognition:             true,
		// SlugVehicle:                     true,
		// SlugVehicleBrandLogo:            true,
	}
)

func (s SlugUpload) GetPath() (string, error) {
	switch s {
	case SlugPortalCommunityPostImage:
		return "community/post/image", nil
	case SlugPortalCommunityPostVideo:
		return "community/post/video", nil
	case SlugPortalCommunityPostDocument:
		return "community/post/document", nil
	case SlugPortalCommunityThumbnail:
		return "community", nil
	case SlugPortalAvatarImage:
		return "user/avatar", nil
	case SlugNewsAndPromotionPicture:
		return "news_and_promotion/picture", nil
	// case SlugOrderPoFile:
	// 	return "order/po", nil
	// case SlugOrderKtpFile:
	// 	return "order/ktp", nil
	// case SlugDealerDocument:
	// 	return "dealer/document", nil
	// case SlugOCRKtp:
	// 	return "ocr/ktp", nil
	// case SlugFaceRecognition:
	// 	return "face_recognition", nil
	// case SlugVehicle:
	// 	return "vehicle", nil
	// case SlugVehicleBrandLogo:
	// 	return "vehicle/brand", nil

	default:
		return "", errors.New("path not found")
	}
}

func (s SlugUpload) CheckExentions(ext string) error {
	result := false
	switch s {
	case
		SlugPortalCommunityPostImage,
		SlugPortalAvatarImage:

		result = isFileExtensionAllowed(ExtSlugPortalCommunityPostImage, ext)
	case ExtSlugPortalCommunityPostDocument:
		result = isFileExtensionAllowed(ExtSlugPortalCommunityPostDocument, ext)
	case SlugPortalCommunityThumbnail:
		result = isFileExtensionAllowed(ExtSlugPortalCommunityThumbnail, ext)
	default:
		return errors.New("extension not allowed")
	}

	if !result {
		err := custom_error.SetCustomError(&custom_error.ErrorContext{
			HTTPCode: http.StatusBadRequest,
			Message:  "please check your file",
		})
		return err
	}

	return nil
}

func isFileExtensionAllowed(extAllowed, fileExtension string) bool {
	allowedFileExtensions := strings.Split(extAllowed, ",")
	if extAllowed == "" || len(allowedFileExtensions) == 0 {
		return true
	}

	for _, allowedFileExtension := range allowedFileExtensions {
		if strings.EqualFold(allowedFileExtension, fileExtension) {
			return true
		}
	}
	return false
}
