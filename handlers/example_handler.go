package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gohttp/auth"
	"gohttp/database"
	"gohttp/snailmail"
	"gohttp/views"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	identity := sessionStore.GetIdentityFromCtx(r)

	viewData := make(map[string]interface{})
	viewData["Title"] = "Example Page"

	model := views.NewViewModel(identity, viewData)
	views.RenderWebpage(w, "example", model)
}

func exampleCounterHandler(w http.ResponseWriter, r *http.Request) {
	viewData := make(map[string]interface{})

	formVal := r.FormValue("number")

	num, _ := strconv.Atoi(formVal)

	viewData["Number"] = num + 1

	model := views.NewViewModel(nil, viewData)
	views.RenderWebpage(w, "example_counter", model)
}

func examplePassgenHandler(w http.ResponseWriter, r *http.Request) {
	identity := sessionStore.GetIdentityFromCtx(r)

	if identity.Claims["CanGeneratePasswords"] != "true" {
		return
	}

	val1 := r.FormValue("val1")

	var password string

	if val1 == "" {
		password = "empty"
	} else {
		password, _ = auth.HashPassword(val1)
	}

	viewData := make(map[string]interface{})
	viewData["Password"] = password

	model := views.NewViewModel(identity, viewData)

	views.RenderWebpage(w, "example_passgen", model)
}

func exampleDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	identity := sessionStore.GetIdentityFromCtx(r)

	viewData := make(map[string]interface{})
	viewData["Title"] = "Database Example"

	model := views.NewViewModel(identity, viewData)

	views.RenderWebpage(w, "example_database", model)
}

func exampleAdduserHandler(w http.ResponseWriter, r *http.Request) {
	viewData := make(map[string]interface{})
	var toomany bool

	users := database.FetchUsers()
	if len(users) >= 5 {
		viewData["TooMany"] = true
		toomany = true
	}

	if r.Method == http.MethodPost && !toomany {
		email := r.FormValue("email")

		_, err := mail.ParseAddress(email)

		if err != nil {
			viewData["Error"] = err
		} else {
			database.AddUser(email)
			viewData["SuccessMsg"] = "Successfully added user " + email + ". ✓"
		}
	}

	users = database.FetchUsers()
	if len(users) >= 5 {
		viewData["TooMany"] = true
		toomany = true
	}

	viewData["Users"] = users

	model := views.NewViewModel(nil, viewData)

	views.RenderWebpage(w, "example_adduser", model)
}

func exampleDeleteallHandler(w http.ResponseWriter, r *http.Request) {
	viewData := make(map[string]interface{})

	database.DeleteAllUsers()
	viewData["SuccessMsg"] = "Successfully deleted all users. ✓"

	model := views.NewViewModel(nil, viewData)

	views.RenderWebpage(w, "example_adduser", model)
}

func exampleUploadHandler(w http.ResponseWriter, r *http.Request) {
	viewData := make(map[string]interface{})
	viewData["Title"] = "File Upload"

	if r.Method == http.MethodPost {
		r.ParseMultipartForm(10 << 20)

		file, fileHeader, err := r.FormFile("file")

		if err != nil {
			log.Println("Error Retrieving the File")
			log.Println(err)
			return
		}

		defer file.Close()

		err = os.MkdirAll("./uploads", os.ModePerm)
		if err != nil {
			log.Println(err)
			return
		}

		dst, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer dst.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println(err)
		}

		dst.Write(fileBytes)

		viewData["UploadSuccessMsg"] = "Successfully uploaded file."
	}

	model := views.NewViewModel(nil, viewData)
	views.RenderWebpage(w, "example_upload", model)
}

func exampleMailHandler(w http.ResponseWriter, r *http.Request) {
	viewData := make(map[string]interface{})
	viewData["Title"] = "Email client"

	if r.Method == http.MethodPost {
		recipients := r.FormValue("to")
		subject := r.FormValue("subject")

		var b bytes.Buffer

		mailViewData := make(map[string]interface{})
		mailViewData["Paragraph"] = r.FormValue("body")
		mailModel := views.NewViewModel(nil, mailViewData)

		views.RenderBytes(&b, "email", mailModel)

		message := snailmail.Email{
			Recipients: []string{recipients},
			Subject:    subject,
			Body:       &b,
		}

		err := snailmail.SendMail(message, snailmail.TYPE_HTML)

		if err != nil {
			viewData["Error"] = "Error sending mail"
		} else {
			viewData["Success"] = "Mail sent successfully"
		}

		model := views.NewViewModel(nil, viewData)
		views.RenderWebpage(w, "example_mail", model)
		return
	}

	model := views.NewViewModel(nil, viewData)
	views.RenderWebpage(w, "example_mail", model)
}

// Struct for deserializing json from 3rd party API
// http://api.open-notify.org/astros.json
type astroModel struct {
	Message string `json:"message"`
	People  []struct {
		Name  string `json:"name"`
		Craft string `json:"craft"`
	} `json:"people"`
	Number int `json:"number"`
}

// This endpoint fetches a JSON response from a third party API,
// serializes it to an 'astro' struct, turns the struct back
// into json, then responds with the result
func exampleFetchHandler(w http.ResponseWriter, r *http.Request) {
	client := http.Client{
		Timeout: time.Second * 5,
	}

	req, reqErr := http.NewRequest(http.MethodGet, "http://api.open-notify.org/astros.json", nil)

	if reqErr != nil {
		log.Println(reqErr)
	}

	req.Header.Set("User-Agent", "Example-Api")

	res, resErr := client.Do(req)

	if resErr != nil {
		log.Println(resErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)

	if readErr != nil {
		log.Println(readErr)
	}

	jsonOutput := astroModel{}

	jsonErr := json.Unmarshal(body, &jsonOutput)

	if jsonErr != nil {
		log.Println(jsonErr)
	}

	viewData := make(map[string]interface{})
	viewData["Title"] = "3rd Party Api Fetch"
	viewData["Data"] = jsonOutput

	model := views.NewViewModel(nil, viewData)
	views.RenderWebpage(w, "api_fetch", model)
}
