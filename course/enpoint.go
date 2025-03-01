package course

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ElianDev55/micro-service-meta-go/metaservice"
	"github.com/gorilla/mux"
)



type (
	Controller func (w http.ResponseWriter, r *http.Request)

	EndPoints struct {
		Create 		Controller
		GetAll 		Controller
		Get 			Controller
		Update		Controller
		Delete 		Controller
	}

	CreateReq struct {
		Name 				string 			`json:"name"`
		StartDate 	string 			`json:"start_date"`
		EndDate 		string 			`json:"end_date"`
	}

	UpdateReq struct {
		Name 				*string 			`json:"name"`
		StartDate 	*string 			`json:"start_date"`
		EndDate 		*string 			`json:"end_date"`
	}

	Response struct {
		Status	 int   					`json:"status"` 
		Data 		interface{} 		`json:"data,omitempty"`
		Err 		string					`json:"error,omitempty"`
		Meta 		*metaservice.Meta			`json:"meta,omitempty"`
	}

)

func MakeEndPoints (s Service)  EndPoints {
	return  EndPoints{
		Create: 		makeCreateEndpoint(s),
		GetAll: 		makeGetAllEndpoint(s),
		Get: 				makeGetEnpoint(s),
		Update: 		makeUpdateEndpoint(s),
		Delete: 		makeDeleteEndPoint(s),
	}
}


func makeCreateEndpoint(s Service) Controller{
	return func(w http.ResponseWriter, r *http.Request) {

		var rq CreateReq

		errJson :=  json.NewDecoder(r.Body).Decode(&rq)

		if errJson != nil {
				w.WriteHeader(400)
		}

		
		if rq.Name == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err: "Name is required",
			})
			return
		}

		if rq.StartDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err: "StartDate is required",
			})
			return
		}

		if rq.EndDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err: "EndDate is required",
			})
			return
		}

		course, err := s.Create(rq.Name,rq.StartDate,rq.EndDate)

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err: err.Error(),
			})
			return
		}
		
		json.NewEncoder(w).Encode(&Response{
			Status: 200,
			Data: course,
			})

	}
}


func makeGetAllEndpoint (s Service) Controller {
	return func (w http.ResponseWriter, r *http.Request) {

		q := r.URL.Query()

		filters := Filters{
			Name: q.Get("name"),
		}

		limit, _ := strconv.Atoi(q.Get("limit"))
		page, _ := strconv.Atoi(q.Get("page"))
		
		count, errCount := s.Count(filters)
		if errCount != nil {
				w.WriteHeader(400)
				json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err: errCount.Error(),
			})
			return
		}
		
		metaservice, errMeta := metaservice.New(page, limit,count)
		if errMeta != nil {
				w.WriteHeader(400)
				json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err: errMeta.Error(),
			})
			return
		}

		courses, err := s.GetAll(filters, metaservice.Offset(), metaservice.Limit())
		if err != nil {
				w.WriteHeader(400)
				json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err: err.Error(),
			})
			return
		}
		
		json.NewEncoder(w).Encode(&Response{
			Status: 200,
				Data: courses,
				Meta: metaservice,
			})
	}

}



	func makeGetEnpoint(s Service) Controller {
		return func (w http.ResponseWriter, r *http.Request) {
			
			path := mux.Vars(r)
			id := path["id"]
			
			course, err := s.Get(id)

			
			if err != nil {
				w.WriteHeader(400)
				json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err: err.Error(),
			})
			return
			}
		
			json.NewEncoder(w).Encode(&Response{
			Status: 200,
				Data: course,
			})
		}
	}


	func makeUpdateEndpoint(s Service) Controller {
		return func (w http.ResponseWriter, r *http.Request) {

			var rq UpdateReq

			err := json.NewDecoder(r.Body).Decode(&rq)
			
			if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err: "Invalided request",
			})
			return
		}


		if rq.Name == nil && rq.StartDate == nil && rq.EndDate == nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err: "You have to send info to update",
			})
			return
		}

		path := mux.Vars(r)
		id := path["id"]

		erro := s.Update(id, rq.Name, rq.StartDate, rq.EndDate)

			if erro != nil {
				w.WriteHeader(400)
				json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err: erro.Error(),
			})
			return
		}

			json.NewEncoder(w).Encode(map[string]bool{"Update": true})

	}
}

	func makeDeleteEndPoint(s Service) Controller {
		return func (w http.ResponseWriter, r *http.Request) {

			
		path := mux.Vars(r)
		id := path["id"]

		err := s.Delete(id)
		
		if err != nil {
				w.WriteHeader(400)
				json.NewEncoder(w).Encode(&Response{
				Status: 400,
				Err: err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(map[string]bool{"Delete": true })


	}
}
