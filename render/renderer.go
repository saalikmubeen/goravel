package render

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/justinas/nosurf"
)

type Render struct {
	Renderer   string // Renderer is the name of the rendering engine that we want to use
	RootPath   string // path to the folder that holds the views
	Secure     bool   // true if we want to use HTTPS
	Port       string
	ServerName string
	JetViews   *jet.Set
	Session    *scs.SessionManager
}

// TemplateData is a struct that holds the data that we want to pass to the templates
type TemplateData struct {
	IsAuthenticated bool
	IntMap          map[string]int
	StringMap       map[string]string
	FloatMap        map[string]float64
	Data            map[string]interface{}
	CSRFToken       string
	Port            string
	ServerName      string
	Secure          bool
	Error           string
	Flash           string
}

func (r *Render) defaultData(td *TemplateData, req *http.Request) *TemplateData {
	if td == nil {
		td = &TemplateData{}
	}

	// check if the userId key exists in the session related to the current request/user
	if r.Session.Exists(req.Context(), "userID") {
		td.IsAuthenticated = true
	}

	td.CSRFToken = nosurf.Token(req) // add the CSRF token to the template data

	td.Port = r.Port
	td.ServerName = r.ServerName
	td.Secure = r.Secure

	// get the error and flash messages from the session and add them to the template data
	// pop deletes the value from the session after it has been retrieved
	td.Error = r.Session.PopString(req.Context(), "error")
	td.Flash = r.Session.PopString(req.Context(), "flash")

	return td
}

// Page Function will render a page
func (r *Render) Page(w http.ResponseWriter, req *http.Request, view string, variables, data interface{}) error {
	// view is the name of the view (or template) that we want to render

	switch strings.ToLower(r.Renderer) {
	case "go":
		// render the page using the Go template engine
		r.GoPage(w, req, view, data)
	case "jet":
		// render the page using the Jet template engine
		r.JetPage(w, req, view, variables, data)
	}
	return nil
}

// GoPage renders a standard Go template
func (r *Render) GoPage(w http.ResponseWriter, req *http.Request, view string, data interface{}) error {
	// render the page using the Go template engine

	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/views/%s.page.tmpl", r.RootPath, view))
	if err != nil {
		return err
	}

	td := &TemplateData{}
	if data != nil {
		templateData, ok := data.(*TemplateData)
		if !ok {
			return fmt.Errorf("data is not of type TemplateData")
		} else {
			td = templateData
		}
	}

	// add the default data to the template data
	td = r.defaultData(td, req)

	err = tmpl.Execute(w, td)
	if err != nil {
		return err
	}

	return nil
}

// JetPage renders a template using the Jet templating engine
func (r *Render) JetPage(w http.ResponseWriter, req *http.Request, templateName string, variables, data interface{}) error {
	var vars jet.VarMap

	if variables == nil {
		vars = make(jet.VarMap)
	} else {
		vars = variables.(jet.VarMap)
	}

	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}

	// add the default data to the template data
	td = r.defaultData(td, req)

	t, err := r.JetViews.GetTemplate(fmt.Sprintf("%s.jet", templateName))
	if err != nil {
		log.Println(err)
		return err
	}

	if err = t.Execute(w, vars, td); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
