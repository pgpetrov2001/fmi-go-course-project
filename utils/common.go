package utils

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func GetIdFromRouteParam(w http.ResponseWriter, r *http.Request, idField string) (uint, error) {
	vars := mux.Vars(r)
	idVal, err := strconv.ParseUint(vars[idField], 10, 64)
	if err != nil {
		http.Error(w, "Specified value is not a valid id", http.StatusBadRequest)
		return 0, err
	}
	id := uint(idVal)
	return id, nil
}
