package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/sebps/template-engine/internal/rendering"
)

const TEMPLATE_DIR = "../templates"

type RequestHandler func(w http.ResponseWriter, r *http.Request)

type HttpHandler struct {
	Pattern string
	Method  string
	Handler RequestHandler
}

func (handler *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.Handler(w, r)
}

func Serve(
	address string,
	port int,
	leftDelimiter string,
	rightDelimiter string,
	leftLoopVariableDelimiter string,
	rightLoopVariableDelimiter string,
	leftLoopBlockDelimiter string,
	rightLoopBlockDelimiter string,
) {

	handlers := []*HttpHandler{
		{
			Pattern: "/",
			Method:  "POST",
			Handler: rootHandler,
		},
		{
			Pattern: "/Render",
			Method:  "POST",
			Handler: getRenderHandler(
				leftDelimiter,
				rightDelimiter,
				leftLoopVariableDelimiter,
				rightLoopVariableDelimiter,
				leftLoopBlockDelimiter,
				rightLoopBlockDelimiter,
			),
		},
		{
			Pattern: "/Register",
			Method:  "POST",
			Handler: registerHandler,
		},
	}

	mux := http.NewServeMux()
	for _, h := range handlers {
		mux.Handle(h.Pattern, h)
	}

	log.Println("Template engine server listening at ", address, ":", port)
	http.ListenAndServe(address+":"+strconv.Itoa(port), mux)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Template engine server listening..."))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	uploadFile(w, r)
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
}

func getRenderHandler(
	leftDelimiter string,
	rightDelimiter string,
	leftLoopVariableDelimiter string,
	rightLoopVariableDelimiter string,
	leftLoopBlockDelimiter string,
	rightLoopBlockDelimiter string,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type Params struct {
			Variables map[string]interface{}
			Template  string
		}
		params := &Params{}

		err := json.NewDecoder(r.Body).Decode(params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		content, err := os.ReadFile(TEMPLATE_DIR + "/" + params.Template)
		if err != nil {
			panic(err)
		}

		rendered := rendering.Render(
			string(content),
			params.Variables,
			leftDelimiter,
			rightDelimiter,
			leftLoopVariableDelimiter,
			rightLoopVariableDelimiter,
			leftLoopBlockDelimiter,
			rightLoopBlockDelimiter,
			false,
		)

		w.Write([]byte(rendered))
	}
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(TEMPLATE_DIR); os.IsNotExist(err) {
		os.Mkdir(TEMPLATE_DIR, 0777)
	}

	// Create file
	dst, err := os.Create(TEMPLATE_DIR + "/" + handler.Filename)
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	defer file.Close()
}
