package appinfoUsecases

import (
	appinfoRepositories "github.com/PHURINTOR/phurinshop/modules/appinfo/appinfoRepositories"
)

// -------------------------------------- Interface ---------------------------------
type IAppinfoUsecase interface {
}

// -------------------------------------- Struct     ---------------------------------
type appinfoUsecase struct {
	appinfoRepository appinfoRepositories.IAppinfoRepository
}

// -------------------------------------- Constructor--------------------------------
func AppinfoUsecase(appinfoRepository appinfoRepositories.IAppinfoRepository) IAppinfoUsecase {
	return &appinfoUsecase{
		appinfoRepository: appinfoRepository,
	}
}

//===========================================================================================
