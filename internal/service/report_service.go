package service

import (
	"context"
	"fmt"
	"os"
	"taskmanager/internal/item"
	"time"
)

const REPORTS_DIR = "reports/"
const REPORT_FILENAME_PATTERN = REPORTS_DIR + "report_%s.txt"

type ReportService struct {
	TaskService *TaskService
	UserService *UserService
}

func NewReportService() *ReportService {
	return &ReportService{
		NewTaskService(),
		NewUserService(),
	}
}

func (s *ReportService) makeReportFileName() string {
	return fmt.Sprintf(REPORT_FILENAME_PATTERN, time.Now().Format("2006-01-02_15-04-05"))
}

func (s *ReportService) makeReportUserRow(user item.User) []byte {
	return fmt.Appendf(nil, "Пользователь %d %s\n", user.Id, user.Login)
}

func (s *ReportService) makeReportTaskRow(task item.Task) []byte {
	return fmt.Appendf(nil, "Задача %d %s %s %s\n", task.Id, task.Title, task.Status, task.Category)
}

func (s *ReportService) GenerateReport(ctx context.Context) error {
	s.TaskService.LoadFromFile(ctx)
	s.UserService.LoadFromFile(ctx)

	if err := os.MkdirAll(REPORTS_DIR, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	reportFile, err := os.Create(s.makeReportFileName())
	if err != nil {
		return err
	}

	defer reportFile.Close()

	users := s.UserService.ListUsers(ctx)
	for _, user := range users {
		_, err = reportFile.Write(s.makeReportUserRow(user))
		if err != nil {
			return err
		}

		userTasks := s.TaskService.ListTasksByUserId(ctx, user.Id)

		for _, task := range userTasks {
			_, err = reportFile.Write(s.makeReportTaskRow(task))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
