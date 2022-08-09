package taskcontroller

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"go-crud-modal-master/entities"
	"go-crud-modal-master/models/taskmodel"
)

var taskModel = taskmodel.New()

func Index(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"data": template.HTML(GetData()),
	}

	temp, _ := template.ParseFiles("views/task/index.html")
	temp.Execute(w, data)
}

func GetData() string {

	buffer := &bytes.Buffer{}

	temp, _ := template.New("data.html").Funcs(template.FuncMap{
		"increment": func(a, b int) int {
			return a + b
		},
	}).ParseFiles("views/task/data.html")

	var task []entities.Tugas
	err := taskModel.FindAll(&task)
	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{
		"task": task,
	}

	temp.ExecuteTemplate(buffer, "data.html", data)

	return buffer.String()
}

func GetForm(w http.ResponseWriter, r *http.Request) {

	queryString := r.URL.Query()
	id, err := strconv.ParseInt(queryString.Get("id"), 10, 64)

	var data map[string]interface{}

	if err != nil {
		data = map[string]interface{}{
			"title": "Tambah Data Task",
		}
	} else {
		var task entities.Tugas
		err := taskModel.Find(id, &task)
		if err != nil {
			panic(err)
		}

		data = map[string]interface{}{
			"title": "Edit Data Task",
			"task":  task,
		}
	}

	temp, _ := template.ParseFiles("views/task/form.html")
	temp.Execute(w, data)
}

func Store(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		r.ParseForm()
		var task entities.Tugas

		task.IsiTask = r.Form.Get("isi_task")
		task.Pegawai = r.Form.Get("pegawai")
		task.Deadline = r.Form.Get("deadline")
		task.Status = r.Form.Get("status")

		id, err := strconv.ParseInt(r.Form.Get("id"), 10, 64)

		var data map[string]interface{}

		if err != nil {
			// insert data
			err := taskModel.Create(&task)
			if err != nil {
				ResponseError(w, http.StatusInternalServerError, err.Error())
				return
			}
			data = map[string]interface{}{
				"message": "Data berhasil disimpan",
				"data":    template.HTML(GetData()),
			}
		} else {
			// mengupdate data
			task.Id = id
			err := taskModel.Update(task)
			if err != nil {
				ResponseError(w, http.StatusInternalServerError, err.Error())
				return
			}
			data = map[string]interface{}{
				"message": "Data berhasil diubah",
				"data":    template.HTML(GetData()),
			}
		}

		ResponseJson(w, http.StatusOK, data)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	id, err := strconv.ParseInt(r.Form.Get("id"), 10, 64)
	if err != nil {
		panic(err)
	}
	err = taskModel.Delete(id)
	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{
		"message": "Data berhasil dihapus",
		"data":    template.HTML(GetData()),
	}
	ResponseJson(w, http.StatusOK, data)
}

func ResponseError(w http.ResponseWriter, code int, message string) {
	ResponseJson(w, code, map[string]string{"error": message})
}

func ResponseJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
