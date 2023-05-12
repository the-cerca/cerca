package controllers

import "net/http"

func InternalServerError(w http.ResponseWriter)  {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}