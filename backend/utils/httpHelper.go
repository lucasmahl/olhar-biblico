package utils

import (
	"encoding/json"
	"net"
	"net/http"
	"olhar-biblico/models/search_filter"
	"strconv"
)

func HttpResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if status >= 400 {
		data = struct {
			Error interface{} `json:"error"`
		}{
			Error: data,
		}
	}

	json.NewEncoder(w).Encode(&data)
}

func ReturnFilter(request *http.Request) search_filter.SearchFilter {
	searchFilter := search_filter.SearchFilter{}

	query := request.URL.Query()

	var limit uint64
	var skip uint64
	sort := query.Get("sort")

	limit, _ = strconv.ParseUint(query.Get("limit"), 10, 64)

	const maxLimit uint64 = 500
	if limit <= 0 {
		limit = 15
	} else if limit >= maxLimit { // Máximo de objetos.
		limit = maxLimit
	}

	skip, _ = strconv.ParseUint(query.Get("skip"), 10, 64)
	if skip <= 0 {
		skip = 0
	}

	order, _ := strconv.Atoi(query.Get("order"))
	if order >= 1 {
		order = 1
	} else {
		order = -1
	}

	deleted := query.Get("deleted")

	searchFilter = search_filter.SearchFilter{
		Limit:   limit,
		Skip:    skip,
		Sort:    sort,
		Order:   order,
		Deleted: deleted,
	}

	return searchFilter
}

func GetIPAddress(r *http.Request) string {
	ipAddress, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ipAddress
}
