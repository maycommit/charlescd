package v1

import (
	"charlescd/internal/manager/project"
	"encoding/json"
	"k8s.io/client-go/dynamic"
	"net/http"
)

func ProjectCreate(client dynamic.Interface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := project.CreateProject(client, r.Body)
		if err != nil {
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
			return
		}
	}
}