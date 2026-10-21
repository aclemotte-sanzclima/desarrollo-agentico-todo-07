# Poder asignar una fecha límite a las tareas

**Issue:** #1 · **Tipo:** feature

## Resumen
Las tareas pueden tener una fecha límite opcional, solo fecha y sin hora. Se indica al crear, se añade, cambia o quita al editar, y la lista la muestra en cada tarea que la tiene. Las tareas existentes quedan sin fecha y funcionan igual.

## Qué se construyó
- Campo opcional «Fecha límite» en el formulario de añadir y en la pantalla de edición.
- En la edición, dejar el campo vacío quita la fecha.
- En la lista, bajo el título: «Fecha límite: DD/MM/AAAA», también en tareas completadas. Sin fecha, no se muestra nada.
- Una fecha no válida se rechaza en el servidor con `422` y el mensaje «La fecha límite no es una fecha válida.».

## Implementación técnica
| Fichero | Papel |
|---|---|
| `migrations/20261021000000_add_due_date_to_tasks.{up,down}.fizz` | Columna `due_date` de tipo `date`, nullable, en `tasks` |
| `models/task.go` | `DueDate *time.Time`; `SetDueDate` (parseo `AAAA-MM-DD`), validación, `DueDateValue` (formulario) y `DueDateText` (lista) |
| `actions/tasks.go` | `taskForm` compartido por crear y actualizar, con `title` y `due_date` |
| `templates/tasks/index.plush.html`, `edit.plush.html` | Campo `type="date"`, `<time datetime>` en la lista y errores por campo |
| `public/assets/app.css` | Estilos de la fecha en la lista |
| `e2e/fixtures.ts` | La fixture `tasks` acepta `dueDate` (`AAAA-MM-DD`) |

Cambios que importan:
- `due_date` es `date` nullable: sin hora, y las filas previas quedan en `NULL`.
- Un valor vacío es `nil`; uno no parseable marca `invalidDueDate` (campo no persistido) y `Validate` lo convierte en error del campo `due_date`.
- Los errores van por campo en las plantillas (`errors.Get("title")` / `errors.Get("due_date")`), para que un error de fecha no marque el título.
- Actualizar el título no toca el estado, y el formulario de edición siempre envía la fecha: vacía la quita.
- Documentación viva actualizada por el propio cambio: `docs/knowledge/features/tareas.md` y `docs/knowledge/overview.md`.

## Cómo se usa
1. En la lista, escribe el título y, si quieres, elige una fecha en el selector junto al botón «Añadir».
2. Para cambiarla, abre la edición de la tarea, modifica «Fecha límite» y guarda.
3. Para quitarla, vacía el campo en la edición y guarda.

## Configuración
Requiere aplicar la migración `20261021000000_add_due_date_to_tasks` (`soda migrate`; la aplicación aplica las pendientes al arrancar). No añade variables de entorno.

## Cómo se probó
Acta de la verificación, validada en el intento 2 de 3: [acta](https://github.com/aclemotte-sanzclima/desarrollo-agentico-todo-07/issues/1#issuecomment-6322496588).

| Criterio | Resultado | Prueba |
|---|---|---|
| 1. Opcional al crear | Verificado | `e2e/tareas/anadir-con-fecha-limite.spec.ts`, `Test_TasksCreate_EmptyDueDate` |
| 2. Añadir, cambiar, quitar | Verificado | `e2e/tareas/cambiar-y-quitar-fecha-limite.spec.ts`, `Test_TasksUpdate_*` |
| 3. Solo fecha | Verificado | `Test_Task_SetDueDate_Persists`, `Test_TasksCreate_InvalidDueDate`; columna `date` medida en la base |
| 4. La lista la muestra | Verificado | `e2e/tareas/lista-muestra-fechas-limite.spec.ts`: rojo sobre `main`, verde sobre el cambio |
| 5. Tareas existentes | Verificado | Suite existente, `Test_TasksIndex_ExistingTasksWithoutDueDate`; migración sobre 3 tareas previas, las 3 sin fecha y sin otros cambios |

Además: fechas inválidas (`2026-02-30`, `mañana`) dan `422` al crear y al editar sin guardar nada; 8 tests de navegador y los paquetes `actions` y `models` en verde.

Para repetir:
```bash
soda migrate -e test
GO_ENV=test go test -p 1 ./actions/... ./models/...
npx playwright test
```

## Notas
- El intento 1 de la verificación se rechazó porque el criterio 4 no tenía un test con varias tareas ([acta](https://github.com/aclemotte-sanzclima/desarrollo-agentico-todo-07/issues/1#issuecomment-6322109820)). La [revisión 2 del plan](https://github.com/aclemotte-sanzclima/desarrollo-agentico-todo-07/issues/1#issuecomment-6322179782) añadió `e2e/tareas/lista-muestra-fechas-limite.spec.ts` y la opción `dueDate` de la fixture `tasks`.
- Fuera de alcance: ordenar o filtrar por fecha, avisos y destacar tareas vencidas.
- Pendientes que señaló el acta, sin bloquear:
  - La ayuda de la edición («Los demás datos de la tarea no cambian.») puede leerse como que la fecha no cambia.
  - `docs/knowledge/features/tareas.md` dice que, al rechazar el formulario, el campo conserva la fecha escrita; si falla la propia fecha, no la conserva. Solo ocurre con peticiones hechas a mano.
  - `docs/procedures/development.md` no menciona la opción `dueDate` de la fixture `tasks`.
  - `DueDateLayout` está exportada y solo se usa en el paquete.
  - En Windows con `core.autocrlf=true`, `gofmt` y los tests del lanzador `adws/` salen en rojo por el entorno, no por el cambio.
  - Ya en `main`: `public/assets/app.js` completa tareas con GET a una ruta que espera PUT.
- El formato de la lista (`Fecha límite: 30/10/2026`) quedó pendiente de la aprobación del humano.
