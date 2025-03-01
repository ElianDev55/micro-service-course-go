package course

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ElianDev55/micro-service-domain-go/domain"
	"gorm.io/gorm"
)



type (

	Repository interface {
		Create(course *domain.Course) error 
		GetAll(filters Filters,  offset, limit int) ([]domain.Course,error)
		Update(id string, name *string, startDate, endDate *time.Time)error
		Get(id string) (*domain.Course,error)
		Delete(id string) error
		Count(filters Filters) (int, error)
	}

	Filters struct {
	Name string
}

	repo struct {
		log *log.Logger
		db *gorm.DB
	}

)


func NewRepo(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: log,
		db: db,
	}
}

func (repo *repo) Create(course *domain.Course) error {
	repo.log.Println("course from repo")

	if err := repo.db.Create(course).Error; err != nil {
		repo.log.Println(err)
		return err
	}

	repo.log.Println("course has been create with id ", course.ID)

	return nil
}

func (repo *repo) GetAll(filters Filters,  offset, limit int) ([]domain.Course, error) {
	repo.log.Println("Start Get all cousers with or no filters")

	var c []domain.Course

	tx := repo.db.Model(&c)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Order("created_at desc").Find(&c)

		if result.Error != nil {
		return nil, result.Error
	}

	return c, nil
	
}

func (r *repo) Update(id string, name *string, startDate, endDate *time.Time) error {

  values := make(map[string]interface{})

	if name != nil {
		values["name"] =  *name
	}

	
	if startDate != nil {
		values["start_date"] =  *startDate
	}

	
	if endDate != nil {
		values["end_date"] =  *endDate
	}
  
	result := r.db.Model(&domain.Course{}).Where("id = ?", id).Updates(values)

	if result.Error != nil {
        return result.Error
    }

  return nil

}


func (r *repo) Get(id string) (*domain.Course,error) {
	course := domain.Course{
		ID: 	id,
	}
	result := r.db.First(&course)
	if result.Error != nil {
		return nil, result.Error
	}
	return &course, nil
}


func (repo *repo) Delete(id string) error {

	repo.log.Println("Get User by id from repo")
	course := domain.Course{
		ID: 	id,
	}
	result := repo.db.Delete(&course)

		if result.Error != nil {
		return result.Error
	}

	return nil

}



func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {

	if filters.Name != "" {
		filters.Name = fmt.Sprintf("%%%s%%", strings.ToLower(filters.Name))
		tx = tx.Where("lower(name) like ?", filters.Name)
	}

	return tx
}


func (r repo) Count(filters Filters) (int, error) {
    var count int64
    tx := r.db.Model(&domain.Course{})
    tx = applyFilters(tx, filters)
    if err := tx.Count(&count).Error; err != nil {
        return 0, err
    }
    return int(count), nil
}
