package filesHandlers

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"

	"github.com/PHURINTOR/phurinshop/config"
	"github.com/PHURINTOR/phurinshop/modules/entities"
	"github.com/PHURINTOR/phurinshop/modules/files"
	fileUsecases "github.com/PHURINTOR/phurinshop/modules/files/filesUsecases"
	"github.com/PHURINTOR/phurinshop/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// ======================================= Enum  ============================================
type fileHandlerErrCode string

const (
	uploadFileErr fileHandlerErrCode = "files-001"
	deleteFileErr fileHandlerErrCode = "files-002"
)

// ======================================= Interface =========================================
type IFilesHandler interface {
	UploadFiles(c *fiber.Ctx) error
	DeleteFile(c *fiber.Ctx) error
}

// ======================================= Struct ============================================
type filesHandler struct {
	cfg          config.IConfig
	fileUsecases fileUsecases.IFilesUsecase
}

// ======================================= Constructor =======================================
func FilesHandler(cfg config.IConfig, fileUsecases fileUsecases.IFilesUsecase) IFilesHandler {
	return &filesHandler{
		cfg:          cfg,
		fileUsecases: fileUsecases,
	}
}

// ======================================= Missing Func =======================================
func (h *filesHandler) UploadFiles(c *fiber.Ctx) error {

	//input =  multipath form
	//Step1 : Req multipath Form
	req := make([]*files.FileReq, 0)
	form, err := c.MultipartForm()
	if err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(uploadFileErr),
			err.Error(),
		).Res()
	}

	filesReq := form.File["files"] //อัพหลายๆ รูป  [files] = ชื่อฟิวตอน upload
	destination := c.FormValue("destination")

	// Files Extention Validation
	extMap := map[string]string{
		"png":  "png",
		"jpg":  "jpg",
		"jpeg": "jpeg",
	}
	// validate extention
	for _, file := range filesReq {
		ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		if extMap[ext] != ext || extMap[ext] == "" {
			return entities.NewErrorResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(uploadFileErr),
				"extension is not acceptable",
			).Res()
		}
		// validate File Size
		if file.Size > int64(h.cfg.App().FileLimit()) {
			return entities.NewErrorResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(uploadFileErr),
				fmt.Sprintf("file size must less than %d MiB", int(math.Ceil(float64(h.cfg.App().FileLimit())/math.Pow(1024, 2)))),
			).Res()
		}

		// Gen name file before add to array
		filesname := utils.RanFileName(ext)
		req = append(req, &files.FileReq{
			File:        file,
			Destination: destination + "/" + filesname,
			FileName:    filesname,
			Extension:   ext,
		})
	}
	// ---------------------- Upload after validation

	res, err := h.fileUsecases.UploadToGCP(req)
	if err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(uploadFileErr),
			err.Error(),
		).Res()

	}

	return entities.NewErrorResponse(c).Success(fiber.StatusCreated, res).Res()
}

// ------------------------------- Delete File
func (h *filesHandler) DeleteFile(c *fiber.Ctx) error {
	req := make([]*files.DeleteFileReq, 0)

	//validate
	if err := c.BodyParser(&req); err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(deleteFileErr),
			err.Error(),
		).Res()
	}
	// ---------------------- deletefile after validation
	if err := h.fileUsecases.DeleteFileGCP(req); err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteFileErr),
			err.Error(),
		).Res()
	}

	return entities.NewErrorResponse(c).Success(fiber.StatusCreated, nil).Res()
}
