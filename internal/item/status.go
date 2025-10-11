package item

type Status int

const (
	StatusInQueue    Status = 0
	StatusInProgress Status = 1
	StatusDone       Status = 2
	StatusFailed     Status = -1
	StatusUnkown     Status = -2
)

func (s Status) String() string {
	switch s {
	case StatusInQueue:
		return "В очереди"
	case StatusInProgress:
		return "В процессе"
	case StatusDone:
		return "Готово"
	case StatusFailed:
		return "Провалено"
	default:
		return "Неизвестно"
	}
}

func StatusFromInt(num int) Status {
	switch num {
	case 0:
		return StatusInQueue
	case 1:
		return StatusInProgress
	case 2:
		return StatusDone
	case 3:
		return StatusFailed
	default:
		return StatusUnkown
	}
}
