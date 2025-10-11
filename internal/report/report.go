package report

import (
	"fmt"
	"os"
	"taskmanager/internal/item"
)

type Task = item.Task

func ExportTasksToText(user string, tasks []Task, path string) error {
	filename := fmt.Sprintf("%s/%s_tasks.txt", path, user)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, task := range tasks {
		fmt.Fprintf(f, "%d[%s] %s - %s\n", task.Id, task.Status.String(), task.Title, task.Category)
	}

	return nil
}
