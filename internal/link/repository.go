package link

import (
	"golang/advanced/pkg/db"
	"gorm.io/gorm/clause"
)

type LinkRepository struct {
	DataBase *db.Db
}

func NewLinkRepository(dataBase *db.Db) *LinkRepository {
	return &LinkRepository{
		DataBase: dataBase,
	}
}

func (repo *LinkRepository) Create(link *Link) (*Link, error) {
	result := repo.DataBase.DB.Create(link)
	if result.Error != nil {
		return nil, result.Error
	}
	return link, nil
}

func (repo *LinkRepository) GetByHash(hash string) (*Link, error) {
	var link Link
	result := repo.DataBase.DB.First(&link, "hash = ?", hash)
	if result.Error != nil {
		return nil, result.Error
	}
	return &link, nil
}

func (repo *LinkRepository) Update(link *Link) (*Link, error) {
	result := repo.DataBase.DB.Clauses(clause.Returning{}).Updates(link)
	if result.Error != nil {
		return nil, result.Error
	}
	return link, nil
}

func (repo *LinkRepository) Delete(id uint) error {
	result := repo.DataBase.DB.Delete(&Link{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *LinkRepository) GetById(id uint) (*Link, error) {
	var link Link
	result := repo.DataBase.DB.First(&link, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &link, nil
}

func (repo *LinkRepository) Count() int64 {
	var count int64
	repo.DataBase.Table("links").
		Where("deleted_at is NULL").
		Count(&count)
	return count
}

func (repo *LinkRepository) GetAll(limit, offset int) []Link {
	var links []Link
	repo.DataBase.
		Table("links").
		Where("deleted_at is NULL").
		Order("id asc").
		Limit(limit).
		Offset(offset).
		Scan(&links)

	return links
}
