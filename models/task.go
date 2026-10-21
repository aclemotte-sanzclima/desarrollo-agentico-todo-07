package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// TitleMaxLength es el número máximo de caracteres permitido en el título.
const TitleMaxLength = 200

// DueDateLayout es el formato de la fecha límite en los formularios (el de <input type="date">).
const DueDateLayout = "2006-01-02"

// dueDateDisplayLayout es el formato de la fecha límite en la lista.
const dueDateDisplayLayout = "02/01/2006"

// Task representa una tarea de la lista.
type Task struct {
	ID        uuid.UUID `json:"id" db:"id" form:"-"`
	Title     string    `json:"title" db:"title" form:"title"`
	Completed bool      `json:"completed" db:"completed" form:"-"`
	CreatedAt time.Time `json:"created_at" db:"created_at" form:"-"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" form:"-"`

	// DueDate es la fecha límite de la tarea, solo la fecha. nil significa sin fecha límite.
	DueDate *time.Time `json:"due_date" db:"due_date" form:"-"`

	// invalidDueDate indica que el último valor recibido en SetDueDate no era una fecha.
	invalidDueDate bool `json:"-" db:"-" form:"-"`
}

// Tasks es una colección de Task.
type Tasks []Task

// Validate se ejecuta en cada llamada a los métodos "pop.Validate*".
func (t *Task) Validate(tx *pop.Connection) (*validate.Errors, error) {
	verrs := validate.NewErrors()

	title := strings.TrimSpace(t.Title)
	if title == "" {
		verrs.Add("title", "El título es obligatorio.")
	} else if len([]rune(title)) > TitleMaxLength {
		verrs.Add("title", fmt.Sprintf("El título no puede superar los %d caracteres.", TitleMaxLength))
	}

	if t.invalidDueDate {
		verrs.Add("due_date", "La fecha límite no es una fecha válida.")
	}

	return verrs, nil
}

// SetDueDate asigna la fecha límite desde el texto de un formulario (AAAA-MM-DD).
// Un valor vacío quita la fecha; un valor que no es una fecha se rechaza en Validate.
func (t *Task) SetDueDate(value string) {
	t.invalidDueDate = false

	value = strings.TrimSpace(value)
	if value == "" {
		t.DueDate = nil
		return
	}

	d, err := time.Parse(DueDateLayout, value)
	if err != nil {
		t.invalidDueDate = true
		return
	}
	t.DueDate = &d
}

// DueDateValue devuelve la fecha límite como AAAA-MM-DD, o vacío si no tiene.
func (t Task) DueDateValue() string {
	if t.DueDate == nil {
		return ""
	}
	return t.DueDate.Format(DueDateLayout)
}

// DueDateText devuelve la fecha límite como DD/MM/AAAA, o vacío si no tiene.
func (t Task) DueDateText() string {
	if t.DueDate == nil {
		return ""
	}
	return t.DueDate.Format(dueDateDisplayLayout)
}
