package appinfoHandlers

import (
	"github.com/PHURINTOR/phurinshop/config"
	"github.com/PHURINTOR/phurinshop/modules/appinfo/appinfoUsecases"
)

// -------------------------------------- Interface ---------------------------------
type IAppinfoHandler interface {
}

// -------------------------------------- Struct     ---------------------------------
type appinfoHandler struct {
	cfg            config.IConfig
	appinfoUsecase appinfoUsecases.IAppinfoUsecase
}

// -------------------------------------- Constructor--------------------------------
func AppinfoHandler(cfg config.IConfig, appinfoUsecase appinfoUsecases.IAppinfoUsecase) IAppinfoHandler {
	return &appinfoHandler{
		cfg:            cfg,
		appinfoUsecase: appinfoUsecase,
	}
}

//===========================================================================================
