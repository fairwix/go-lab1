package model

import (
	"fmt"
	"time"
)

// TimeSlot — интервал времени бронирования (например, "с 10:00 до 11:00").
//
// Это value type: маленькая структура из двух полей, у неё нет собственной
// идентичности и её не нужно "делить" между разными частями программы —
// каждому, кому нужен TimeSlot, достаточно отдать копию.
type TimeSlot struct {
	Start time.Time
	End   time.Time
}

// NewTimeSlot — функция-конструктор: валидирует входные данные и
// возвращает готовый TimeSlot по значению.
func NewTimeSlot(start, end time.Time) (TimeSlot, error) {
	if !end.After(start) {
		return TimeSlot{}, ValidationError{
			Message: "time slot end must be after start",
		}
	}

	return TimeSlot{Start: start, End: end}, nil
}

// Duration только читает поля TimeSlot, поэтому объявлен с value receiver:
// вызывающему не нужно давать доступ к оригиналу — копии достаточно,
// и метод гарантированно не может случайно изменить чужой TimeSlot.
func (t TimeSlot) Duration() time.Duration {
	return t.End.Sub(t.Start)
}

func (t TimeSlot) String() string {
	return fmt.Sprintf(
		"%s → %s (%s)",
		t.Start.Format(time.RFC3339),
		t.End.Format(time.RFC3339),
		t.Duration(),
	)
}

// ShiftSlot переносит интервал на d вперёд.
//
// Сравните с Duration()/String() выше — вот та самая демонстрация разницы
// между value и pointer, которую требует ТЗ:
//   - Duration() и String() берут TimeSlot ПО ЗНАЧЕНИЮ, потому что только
//     читают его и не должны иметь возможность что-то в нём поменять;
//   - ShiftSlot берёт *TimeSlot (по указателю), потому что должен
//     ИЗМЕНИТЬ существующий слот у вызывающего кода (перенести бронь),
//     а не вернуть новую копию, которую тот может забыть присвоить обратно.
//
// Правило простое: value — когда только читаем, pointer — когда мутируем
// и хотим, чтобы изменение было видно вызывающей стороне.
func ShiftSlot(t *TimeSlot, d time.Duration) {
	t.Start = t.Start.Add(d)
	t.End = t.End.Add(d)
}
