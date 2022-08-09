package taskmodel

import (
	"database/sql"

	"go-crud-modal-master/config"
	"go-crud-modal-master/entities"
)

type TaskModel struct {
	db *sql.DB
}

func New() *TaskModel {
	db, err := config.DBConnection()
	if err != nil {
		panic(err)
	}
	return &TaskModel{db: db}
}

func (m *TaskModel) FindAll(task *[]entities.Tugas) error {
	rows, err := m.db.Query("select * from tugas")
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var data entities.Tugas
		rows.Scan(
			&data.Id,
			&data.IsiTask,
			&data.Pegawai,
			&data.Deadline,
			&data.Status)
		*task = append(*task, data)
	}

	return nil
}

func (m *TaskModel) Create(task *entities.Tugas) error {
	result, err := m.db.Exec("insert into tugas (isi_task, pegawai, deadline, status) values(?,?,?,?)",
		task.IsiTask, task.Pegawai, task.Deadline, task.Status)

	if err != nil {
		return err
	}

	lastInsertId, _ := result.LastInsertId()
	task.Id = lastInsertId
	return nil
}

func (m *TaskModel) Find(id int64, task *entities.Tugas) error {
	return m.db.QueryRow("select * from tugas where id = ?", id).Scan(
		&task.Id,
		&task.IsiTask,
		&task.Pegawai,
		&task.Deadline,
		&task.Status)
}

func (m *TaskModel) Update(task entities.Tugas) error {

	_, err := m.db.Exec("update tugas set isi_task = ?, pegawai = ?, deadline = ?, status = ? where id = ?",
		task.IsiTask, task.Pegawai, task.Deadline, task.Status, task.Id)

	if err != nil {
		return err
	}

	return nil
}

func (m *TaskModel) Delete(id int64) error {
	_, err := m.db.Exec("delete from tugas where id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
