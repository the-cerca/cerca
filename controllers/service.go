package controllers

import (
	"io"
	"net/http"
	"os"
	"strconv"

	m "github.com/aleeXpress/cerca/models"
	"github.com/aleeXpress/cerca/utils"
)

type ServiveC struct {
	Sm      *m.ServiceManager
	Im      *m.ImageManager
	Cs      *m.CategoryManager
	Session *m.SessionManager
}

func (sc *ServiveC) CreateService(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserByContext(r.Context())
	if err != nil {
		InternalServerError(w)
		return 
	}
	// Parse and create file 
	r.ParseMultipartForm(10 << 20)
	file,FileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "can't parse file", http.StatusBadRequest)
	}
	defer file.Close()
	path := "./images/images_service/" + FileHeader.Filename
	f, err := os.Create(path)
	if err != nil {
		InternalServerError(w)
	}
	defer f.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		InternalServerError(w)
	}
	if _, err := f.Write(bytes); err != nil {
		InternalServerError(w)
	}
	img, err := sc.Im.CreateToDatabase(FileHeader.Filename, path, "")
	if err != nil {
		    panic(err)
	}
	// Parse information service
	name := r.FormValue("name")
	description := r.FormValue("description")
	category := r.FormValue("category")
	startingAt := r.FormValue("startingAt")
	c, err := sc.Cs.SearchCategory(category)
	if err != nil {
		http.Error(w, "search category went wrong", http.StatusBadRequest)
	}
	float, err := strconv.ParseFloat(startingAt, 64)
	if err != nil {
		InternalServerError(w)
	}
	if !user.IsFreelance{
		http.Error(w,"not allowed", http.StatusBadRequest)
		return
	}
	service,err := sc.Sm.CreateService(user.ID, name, description, float, img.ID,c.ID,)
	if err != nil {InternalServerError(w)}
	utils.Encode(w,service)
}


func (sc *ServiveC)UpdateService(w http.ResponseWriter, r *http.Request)  {
	r.ParseMultipartForm(10 << 20)
	file,FileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "can't parse file", http.StatusBadRequest)
	}
	defer file.Close()
	path := "./images/images_service/" + FileHeader.Filename
	f, err := os.Create(path)
	if err != nil {
		InternalServerError(w)
	}
	defer f.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		InternalServerError(w)
	}
	if _, err := f.Write(bytes); err != nil {
		InternalServerError(w)
	}
	img, err := sc.Im.CreateToDatabase(FileHeader.Filename, path, "")
	if err != nil {
		    panic(err)
	}

	id := r.FormValue("id")
	name := r.FormValue("name")
	description := r.FormValue("description")
	category := r.FormValue("category")
	startingAt := r.FormValue("startingAt")
	float, err := strconv.ParseFloat(startingAt, 64)
	if err != nil {
		InternalServerError(w)
	}
	c, err := sc.Cs.SearchCategory(category)
	if err != nil {
		http.Error(w, err.Error(),http.StatusBadRequest)
	}
	s, err := sc.Sm.UpdateService(id, name, description,float,c.ID, img.ID )
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	
	utils.Encode(w, s)
}


func (sc *ServiveC)DeleteService(w http.ResponseWriter, r *http.Request)  {}

func (sc *ServiveC)GetServices(w http.ResponseWriter, r *http.Request)  {}

func (sc *ServiveC)GetAllServices(w http.ResponseWriter, r *http.Request)  {}

func (sc *ServiveC)GetAllCategories(w http.ResponseWriter, r *http.Request) {

  categories, err := sc.Cs.GetAllCategories()
  if err!= nil {
    InternalServerError(w)
    return 
  }
  utils.Encode(w, categories)
}