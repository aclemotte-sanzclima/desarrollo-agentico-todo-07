# Lista de tareas

Lo que ve y puede hacer el usuario con sus tareas. Todo pasa en dos pantallas: la lista (`/`) y la edición de una tarea (`/tasks/{task_id}/edit`). El código está en `actions/tasks.go`, `models/task.go`, `templates/tasks/` y `public/assets/app.js`.

## La lista

- La pestaña y el enlace de la navegación lateral dicen «Tareas». El título de la página también.
- Debajo del título, el número de tareas: «3 tareas en tu lista», y en singular «1 tarea en tu lista».
- Las tareas salen por orden de creación, de la más antigua a la más reciente.
- Cada tarea muestra:
  - un botón redondo para completarla o reabrirla;
  - su título;
  - su fecha límite, si tiene, debajo del título: «Fecha límite: 30/10/2026». Las tareas sin fecha no muestran nada;
  - la etiqueta «Completada», si lo está;
  - un enlace para editarla y un botón para eliminarla.
- Los botones solo tienen icono, y su nombre accesible incluye el título: «Completar tarea: Comprar pan», «Reabrir tarea: …», «Editar tarea: …», «Eliminar tarea: …».
- Sin tareas, la lista muestra «No hay tareas todavía. Añade la primera con el formulario.», y el número dice «0 tareas en tu lista».

## Reglas del título

Las comprueba `Task.Validate` al crear y al editar:

| Regla | Mensaje si no se cumple |
| --- | --- |
| Es obligatorio. Un título solo con espacios cuenta como vacío | «El título es obligatorio.» |
| Tiene como mucho 200 caracteres (`models.TitleMaxLength`), contados sin los espacios del principio y del final | «El título no puede superar los 200 caracteres.» |

Antes de guardar, la acción quita los espacios del principio y del final. Los campos de título del navegador tienen también `maxlength="200"`.

## Reglas de la fecha límite

Es opcional y es solo una fecha, sin hora. En los formularios se escribe con un campo de fecha del navegador y el servidor solo acepta `AAAA-MM-DD`. Se admite cualquier fecha válida, también una pasada. La comprueba `Task.Validate` al crear y al editar:

| Regla | Mensaje si no se cumple |
| --- | --- |
| Si se indica, es una fecha válida en formato `AAAA-MM-DD` (no valen una hora, `30/10/2026` ni `2026-02-30`) | «La fecha límite no es una fecha válida.» |

Un campo vacío, o solo con espacios, significa sin fecha límite.

## Añadir

- El formulario está debajo de la lista: el campo «Nueva tarea», con el texto de ayuda «Añadir una tarea…», el campo de fecha «Fecha límite (opcional)» y el botón «Añadir».
- Si el título y la fecha son válidos, la tarea se crea pendiente, con su fecha límite si se indicó, aparece al final de la lista y sale el aviso «Tarea creada.».
- Si no, la lista se vuelve a mostrar con `422`: los campos conservan lo escrito (título y fecha), el campo que falla queda marcado como inválido y el mensaje sale debajo del formulario.

## Editar

- Se llega desde el enlace de editar de cada tarea. La pantalla tiene el enlace «← Volver a la lista», el título «Editar tarea» y el campo «Título» con el título actual. Debajo, la ayuda: «Hasta 200 caracteres. Los demás datos de la tarea no cambian.». Después, el campo «Fecha límite» con la fecha actual (vacío si no tiene) y la ayuda «Opcional. Déjala vacía para quitarla.».
- «Guardar cambios» cambia el título y la fecha límite: se puede añadir, cambiar o quitar la fecha (vaciando el campo). Si la tarea estaba completada, lo sigue estando. Después vuelve a la lista con el aviso «Tarea actualizada.».
- Con un título o una fecha inválidos, la pantalla se vuelve a mostrar con `422` y el mensaje. «Cancelar» vuelve a la lista sin cambiar nada.
- Si la tarea no existe, `404`.

## Completar y reabrir

- El botón de cada tarea la alterna entre pendiente y completada. Una tarea completada lleva la etiqueta «Completada», y su botón pasa a reabrirla.
- Sin JavaScript, el formulario envía `PUT /tasks/{task_id}/toggle`, la acción cambia el estado y vuelve a la lista. No hay aviso.
- Con JavaScript, `app.js` cambia la tarea en la página sin recargarla (la etiqueta, el estado del botón y su nombre accesible) y después llama a `fetch(form.action)`.

## Eliminar

- El botón de eliminar pide confirmación: «¿Seguro que quieres eliminar esta tarea?». Si se cancela, no pasa nada. Sin JavaScript, elimina sin preguntar.
- Al eliminar, la tarea desaparece de la lista y sale el aviso «Tarea eliminada.».
- Si la tarea no existe, `404`.

## Avisos

Los avisos salen arriba, en la primera pantalla que se muestra después de la acción, y se anuncian a los lectores de pantalla (`role="status"`):

| Acción | Aviso |
| --- | --- |
| Añadir | «Tarea creada.» |
| Editar | «Tarea actualizada.» |
| Eliminar | «Tarea eliminada.» |
| Completar o reabrir | Ninguno |
