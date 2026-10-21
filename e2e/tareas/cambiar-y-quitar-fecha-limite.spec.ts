import { test, expect } from '../fixtures';

test.describe('Tareas', () => {
  test('cambiar y quitar la fecha límite', async ({ page, tasks }) => {
    await tasks([{ title: 'Llamar al proveedor' }]);
    await page.reload();
    await expect(page.locator('.task-due-date')).toHaveCount(0);

    // 1. Añade una fecha límite desde la edición
    await page.getByRole('link', { name: 'Editar tarea: Llamar al proveedor' }).click();
    await expect(page.getByLabel('Fecha límite')).toHaveValue('');
    await page.getByLabel('Fecha límite').fill('2026-11-15');
    await page.getByRole('button', { name: 'Guardar cambios' }).click();
    await expect(page.locator('.task-due-date')).toHaveText('Fecha límite: 15/11/2026');

    // 2. Cambia la fecha límite
    await page.getByRole('link', { name: 'Editar tarea: Llamar al proveedor' }).click();
    await expect(page.getByLabel('Fecha límite')).toHaveValue('2026-11-15');
    await page.getByLabel('Fecha límite').fill('2026-11-20');
    await page.getByRole('button', { name: 'Guardar cambios' }).click();
    await expect(page.locator('.task-due-date')).toHaveText('Fecha límite: 20/11/2026');

    // 3. Quita la fecha límite vaciando el campo
    await page.getByRole('link', { name: 'Editar tarea: Llamar al proveedor' }).click();
    await page.getByLabel('Fecha límite').fill('');
    await page.getByRole('button', { name: 'Guardar cambios' }).click();
    await expect(page.locator('.task-due-date')).toHaveCount(0);
    await expect(page.locator('.task-title')).toHaveText(['Llamar al proveedor']);
  });
});
