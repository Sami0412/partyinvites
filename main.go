package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Rsvp struct {
	Name, Email, Phone string
	WillAttend         bool
}

var responses = make([]*Rsvp, 0, 10)
var templates = make(map[string]*template.Template, 3)

func loadTemplates() {
	//Declare an array of the template file names
	templateNames := [5]string{"welcome", "form", "thanks", "sorry", "list"}
	//Loop through the array of names
	for index, name := range templateNames {
		//Parse each template file with the layout to create the 5 webpages
		t, err := template.ParseFiles("layout.html", name+".html") 
		if err == nil {
			//Save the parsed template file (type *Template.template) as the value against the current template name key (type string e.g. "welcome") in the templates array
			templates[name] = t
			fmt.Println("Loaded template", index, name)
		} else {
			//else throw an error
			panic(err)
		}
	}
}

func welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	//Read the Template.template value of the "welcome" key and invoke the in-built Execute() method
	templates["welcome"].Execute(writer, nil)
}

func listHandler(writer http.ResponseWriter, request *http.Request) {
	templates["list"].Execute(writer, responses)
}

//formData struct has a pointer to the existing Rsvp struct, so can be used to define the fields inside the Rsvp struct - can create an instance of formData using an existing Rsvp value
type formData struct {
	*Rsvp
	Errors []string
}

func formHandler(writer http.ResponseWriter, request *http.Request) {
	//There is no data to use when responding to GET requests, but template still expects formData data structure - use default values for fields
	if request.Method == http.MethodGet {
		templates["form"].Execute(writer, formData {
			Rsvp: &Rsvp{}, Errors: []string {},		//formData struct expects a pointer to a Rsvp value - created with &
		})
	} else if request.Method == http.MethodPost {
		request.ParseForm()		//Parses form data in an HTTP request and populates a map which can be accessed through the Form field
		responseData := Rsvp {		//Use form data to create an Rsvp value\
			Name: request.Form["name"][0],		//form data is presented as a slice
			Email: request.Form["email"][0],
			Phone: request.Form["phone"][0],
			WillAttend: request.Form["willattend"][0] == "true",
		}

		//validation
		errors := []string {}
		if responseData.Name == "" {
			errors = append(errors, "Please enter your name")
		}
		if responseData.Email == "" {
			errors = append(errors, "Please enter your email address")
		}
		if responseData.Phone == "" {
			errors = append(errors, "Please enter your phone number")
		}
		if len(errors) > 0 {
			//render the form again including the error messages
			templates["form"].Execute(writer, formData {
				Rsvp: &responseData, Errors: errors,
			})
		} else {
			responses = append(responses, &responseData)	//if you don't use the & pointer, the responseData object would be duplicated

			if responseData.WillAttend {
				templates["thanks"].Execute(writer, responseData.Name)
			} else {
			templates["sorry"].Execute(writer, responseData.Name)
			}
		}
	}
}

func main() {
	loadTemplates()

	//HandleFunc() from http package is used to specify a URL path and rhe handler that will receive matching requests
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/form", formHandler)

	//Create an HTTP server listening on port 5000 - second argument nil tells server to use handlers registered with HandleFunc()
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Println(err)
	}
}
 