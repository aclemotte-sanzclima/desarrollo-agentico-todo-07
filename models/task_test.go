package models

import "strings"

func (ms *ModelSuite) Test_Task_Create_Valid() {
	task := &Task{Title: "Comprar café"}

	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.False(verrs.HasAny())

	ms.NotZero(task.ID)
	ms.False(task.Completed)

	count, err := ms.DB.Count(&Task{})
	ms.NoError(err)
	ms.Equal(1, count)
}

func (ms *ModelSuite) Test_Task_Title_Required() {
	task := &Task{Title: ""}

	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.True(verrs.HasAny())
	ms.Contains(verrs.Get("title"), "El título es obligatorio.")

	count, err := ms.DB.Count(&Task{})
	ms.NoError(err)
	ms.Equal(0, count)
}

func (ms *ModelSuite) Test_Task_Title_OnlySpaces_Invalid() {
	task := &Task{Title: "   \t  "}

	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.True(verrs.HasAny())
	ms.Contains(verrs.Get("title"), "El título es obligatorio.")
}

func (ms *ModelSuite) Test_Task_Title_MaxLength() {
	task := &Task{Title: strings.Repeat("a", TitleMaxLength)}

	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.False(verrs.HasAny())
}

func (ms *ModelSuite) Test_Task_Title_TooLong_Invalid() {
	task := &Task{Title: strings.Repeat("a", TitleMaxLength+1)}

	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.True(verrs.HasAny())
	ms.NotEmpty(verrs.Get("title"))
}

func (ms *ModelSuite) Test_Task_Toggle_Persists() {
	task := &Task{Title: "Revisar el correo"}
	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.False(verrs.HasAny())

	task.Completed = true
	ms.NoError(ms.DB.Update(task))

	reloaded := &Task{}
	ms.NoError(ms.DB.Find(reloaded, task.ID))
	ms.True(reloaded.Completed)
}

func (ms *ModelSuite) Test_Task_DueDate_OptionalByDefault() {
	task := &Task{Title: "Sin fecha"}
	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.False(verrs.HasAny())

	reloaded := &Task{}
	ms.NoError(ms.DB.Find(reloaded, task.ID))
	ms.Nil(reloaded.DueDate)
	ms.Equal("", reloaded.DueDateValue())
	ms.Equal("", reloaded.DueDateText())
}

func (ms *ModelSuite) Test_Task_SetDueDate_Persists() {
	task := &Task{Title: "Con fecha"}
	task.SetDueDate("2026-10-30")
	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.False(verrs.HasAny())

	reloaded := &Task{}
	ms.NoError(ms.DB.Find(reloaded, task.ID))
	ms.Equal("2026-10-30", reloaded.DueDateValue())
	ms.Equal("30/10/2026", reloaded.DueDateText())
}

func (ms *ModelSuite) Test_Task_SetDueDate_Empty_Removes() {
	for _, empty := range []string{"", "   "} {
		task := &Task{Title: "Con fecha"}
		task.SetDueDate("2026-10-30")
		verrs, err := ms.DB.ValidateAndCreate(task)
		ms.NoError(err)
		ms.False(verrs.HasAny())

		task.SetDueDate(empty)
		verrs, err = ms.DB.ValidateAndUpdate(task)
		ms.NoError(err)
		ms.False(verrs.HasAny())

		reloaded := &Task{}
		ms.NoError(ms.DB.Find(reloaded, task.ID))
		ms.Nil(reloaded.DueDate)
	}
}

func (ms *ModelSuite) Test_Task_SetDueDate_Invalid() {
	for _, value := range []string{"30/10/2026", "2026-02-30", "2026-10-30T10:00", "mañana"} {
		task := &Task{Title: "Fecha mala"}
		task.SetDueDate(value)
		verrs, err := ms.DB.ValidateAndCreate(task)
		ms.NoError(err)
		ms.Contains(verrs.Get("due_date"), "La fecha límite no es una fecha válida.", value)
	}

	count, err := ms.DB.Count(&Task{})
	ms.NoError(err)
	ms.Equal(0, count)
}

func (ms *ModelSuite) Test_Task_SetDueDate_InvalidThenValid() {
	task := &Task{Title: "Fecha corregida"}
	task.SetDueDate("mañana")
	task.SetDueDate("2026-10-30")

	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.False(verrs.HasAny())
	ms.Equal("2026-10-30", task.DueDateValue())
}
