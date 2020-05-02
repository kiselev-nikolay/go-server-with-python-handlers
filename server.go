package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	py3 "github.com/DataDog/go-python3"
)

func myHandler(writer http.ResponseWriter, request *http.Request) {
	someDict := py3.PyDict_New()
	py3.PySys_SetObject("scope", someDict)
	defer someDict.DecRef()

	requestQuery := request.URL.Query()
	for key, values := range requestQuery {
		if len(values) > 1 {
			args := py3.PyTuple_New(len(values))
			defer args.DecRef()
			for i, v := range values {
				index := int(i)
				py3.PyTuple_SetItem(args, index, py3.PyUnicode_FromString(v))
			}
			someDict.SetItem(py3.PyUnicode_FromString(key), args)
		} else {
			someDict.SetItem(py3.PyUnicode_FromString(key), py3.PyUnicode_FromString(values[0]))
		}
	}
	someDict.SetItem(py3.PyUnicode_FromString("query"), py3.PyUnicode_FromString(""))
	py3.PySys_SetObject("scope", someDict)

	py3.PyRun_SimpleString("sys.response = (handler(sys.scope))")

	response := py3.PySys_GetObject("response")
	respUTF8 := py3.PyUnicode_AsUTF8(response)
	writer.Write([]byte(respUTF8))
}

func main() {
	someChan := make(chan os.Signal, 1)
	signal.Notify(someChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-someChan
		log.Print("...")
		py3.Py_Finalize()
		os.Exit(0)
	}()

	py3.Py_Initialize()
	// py3.PyRun_SimpleString("from importlib import reload; import sys; reload(sys)")
	py3.PyRun_SimpleString("import sys")
	py3.PyRun_AnyFile("handler.py")

	http.HandleFunc("/", myHandler)

	log.Printf("runin: http://localhost:8080")

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
