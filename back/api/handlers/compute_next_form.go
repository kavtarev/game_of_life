package handlers

import (
	"encoding/json"
	"game_of_life/internal"
	"math"
	"net/http"
)

type Dto struct {
	Data string `json:"data"`
}

type Response struct {
	Data string `json:"data"`
}

func HandleComputeNextForm(w http.ResponseWriter, r *http.Request) {
	var b Dto

	if r.Method != http.MethodPost {
		w.WriteHeader(404)
		w.Write([]byte("Only Post"))
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	if err := d.Decode(&b); err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Invalid Data"))
		return
	}

	root := int(math.Sqrt(float64(len(b.Data))))

	f := internal.NewField(root)
	f.Update(b.Data)
	str := f.Run()

	resp := Response{Data: str}

	m, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	w.Write(m)
}
