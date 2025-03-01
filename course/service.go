package course

import (
	"log"
	"time"

	"github.com/ElianDev55/micro-service-domain-go/domain"
)


type (

	Service interface {
		Create(name, startDate, endDate string) (*domain.Course, error)
		GetAll(filters Filters,  offset, limit int) ([]domain.Course,error)
		Get(id string) (*domain.Course,error)
		Update(id string, name *string, startDate *string, endDate *string) error
		Delete(id string) error
		Count(filters Filters) (int, error)
	}

	service struct {
		log *log.Logger
		repo Repository
	}

)

func NewService (l *log.Logger, repo Repository) Service {
	return &service{
		log:  l,
		repo: repo,
	}
}


func (s service) Count(filters Filters) (int, error) {
	return s.repo.Count(filters)
}


func (s service) Create(name, startDate, endDate string) (*domain.Course, error) {
	
	startDateParsed, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}

	endDateParsed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}

	course := &domain.Course{
		Name: name,
		StartDate: startDateParsed ,
		EndDate:  endDateParsed,
	}

	errDb := s.repo.Create(course)

	if errDb != nil {
		s.log.Println(errDb)
		return nil,errDb
	}

	return course, nil


}


func (s service) GetAll(filters Filters,  offset, limit int) ([]domain.Course,error) {
	s.log.Println("Start GetAll on service")

	courses, error := s.repo.GetAll(filters, offset,limit)


	if error != nil {
		return nil,error
	}

	return courses,nil
}

func (s service) Get(id string) (*domain.Course, error) {

	course, err := s.repo.Get(id)

	if err != nil {
		return nil,err
	}

	return course,nil

}


func (s service) Update(id string, name *string, startDate *string , endDate *string) error {

	var startDateParsed, endDateParsed * time.Time

	if startDate != nil {
		data , err := time.Parse("2006-01-2", *startDate)
		if err != nil {
			s.log.Println(err)
			return err
		}
		startDateParsed = &data
	}

	
	if endDate != nil {
		data , err := time.Parse("2006-01-2", *endDate)
		if err != nil {
			s.log.Println(err)
			return err
		}
		endDateParsed = &data
	}

	err := s.repo.Update(id,name, startDateParsed, endDateParsed)

	if err != nil {
	s.log.Println("ERORORRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRR")
		return err
	}
	return nil
}

func (s service) Delete(id string) error {

	err := s.repo.Delete(id)

	
	if err != nil {
		return err
	}

	return nil

} 
