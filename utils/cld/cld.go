package cld

import (
	"BE-hi-SPEC/config"
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/labstack/gommon/log"
)

type params struct {
	cloudname string
	key       string
	secret    string
}

func InitCloudnr(c config.AppConfig) (*cloudinary.Cloudinary, context.Context, string) {
	var config params
	config.cloudname = c.CLOUDINARY_CLD
	config.key = c.CLOUDINARY_KEY
	config.secret = c.CLOUDINARY_SECRET
	param := c.CLOUDINARY_FOLDER
	cld, err := cloudinary.NewFromParams(config.cloudname, config.key, config.secret)
	if err != nil {
		log.Error("terjadi kesalahan pada koneksi cloudinary, error:", err.Error())
		return nil, nil, ""
	}

	ctx := context.Background()
	return cld, ctx, param
}

func UploadImage(cld *cloudinary.Cloudinary, ctx context.Context, image multipart.File, filePath string) (string, error) {

	resp, err := cld.Upload.Upload(ctx, image, uploader.UploadParams{
		Folder: filePath,
	})
	if err != nil {
		log.Error("terjadi kesalahan pada upload, error:", err.Error())
		return "", err
	}
	return resp.SecureURL, nil
}
