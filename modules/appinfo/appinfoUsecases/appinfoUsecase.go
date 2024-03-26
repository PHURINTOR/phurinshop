package appinfoUsecases

import (
	"github.com/PHURINTOR/phurinshop/modules/appinfo"
	appinfoRepositories "github.com/PHURINTOR/phurinshop/modules/appinfo/appinfoRepositories"
)

// ======================================= Interface =========================================
type IAppinfoUsecase interface {
	FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error)
	InsertCategory(req []*appinfo.Category) error
	DeleteCategory(categoryId int) error // Delete Category
}

// ======================================= Struct ============================================
type appinfoUsecase struct {
	appinfoRepository appinfoRepositories.IAppinfoRepository
}

// ======================================= Constructor =======================================
func AppinfoUsecase(appinfoRepository appinfoRepositories.IAppinfoRepository) IAppinfoUsecase {
	return &appinfoUsecase{
		appinfoRepository: appinfoRepository,
	}
}

// =======================================Missing Function =======================================

// -------------------------------------- Find Category -------------------------------------
func (u *appinfoUsecase) FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error) {
	category, err := u.appinfoRepository.FindCategory(req)
	if err != nil {
		return nil, err
	}
	return category, nil
}

// -------------------------------------- Insert Category -------------------------------------
func (u *appinfoUsecase) InsertCategory(req []*appinfo.Category) error { // Category Array เป็น pointer อยู่แล้วไม่จำเป็นต้อง return ออกไป
	return u.appinfoRepository.InsertCategory(req)

}

// -------------------------------------- Delete Category -------------------------------------
func (u *appinfoUsecase) DeleteCategory(categoryId int) error {
	if err := u.appinfoRepository.DeleteCategory(categoryId); err != nil {
		return err
	}
	return nil
}
