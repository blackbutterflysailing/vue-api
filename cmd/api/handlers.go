package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"vue-api/internal/data"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Numa    int64  `json:"numa"`
}

func (app *application) SaveCalc(w http.ResponseWriter, r *http.Request) {
	type nums struct {
		Numa int64 `json:"numa"`
	}

	var numas nums
	var payload jsonResponse

	fmt.Println("SaveCalc...")
	app.infoLog.Println(r.Body)
	err := json.NewDecoder(r.Body).Decode(&numas)
	if err != nil {
		// send back an error message
		app.errorLog.Println("invalid json")
		payload.Error = true

		payload.Message = "invalid json"

		out, err := json.MarshalIndent(payload, "", "\t")
		if err != nil {
			app.errorLog.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(out)
		return
	}

	// save to the database
	app.infoLog.Printf("save to database %d", numas.Numa)

	//var calculated data.Calculated
	calculated := data.Calculated{
		ID: 1,
		// ID:     app.models.Calculated.ID,
		Result: int(numas.Numa),
	}
	fmt.Printf("Handlers:SaveCalc: calculated %v", calculated)
	err = calculated.Update()
	if err != nil {
		return
	}

	payload.Error = false
	payload.Message = "saved"

	// send back a response
	app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) ReadCalc(w http.ResponseWriter, r *http.Request) {

	// read from the database
	var calculated data.Calculated
	calc, err := calculated.GetCalculated()
	if err != nil {
		return
	}

	// send back a response
	app.writeJSON(w, http.StatusOK, calc)
}
