package v1

import (
	"charlescd/internal/manager/circle"
	"encoding/json"
	"net/http"

	"k8s.io/client-go/dynamic"
)

func CircleCreate(client dynamic.Interface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := circle.CreateCircle(client, r.Body)
		if err != nil {
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
			return
		}
	}
}

func CircleFindAll(client dynamic.Interface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		circles, err := circle.ListCircles(client)
		if err != nil {
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
			return
		}

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(circles)
	}
}
