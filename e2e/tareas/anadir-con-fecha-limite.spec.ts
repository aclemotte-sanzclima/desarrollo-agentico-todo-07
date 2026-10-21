import { test, expect } from '../fixtures';

test.describe('Tareas', () => {
  test('añadir una tarea con fecha límite', async ({ page }) => {
    // 1. Rellena el título y la fecha límite
    await page.getByRole('textbox', { name: 'Nueva tarea' }).fill('Pagar el alquiler');
    await page.getByLabel('Fecha límite').fill('2026-10-30');

    // 2. Añade la tarea
    await page.getByRole('button', { name: 'Añadir' }).click();

    // 3. La lista la muestra con su fecha límite
    await expect(page.getByRole('status')).toHaveText('Tarea creada.');
    await expect(page.locator('.task-title')).toHaveText(['Pagar el alquiler']);
    await expect(page.locator('.task-due-date')).toHaveText('Fecha límite: 30/10/2026');

    // 4. La fecha se conserva al recargar
    await page.reload();
    await expect(page.locator('.task-due-date')).toHaveText('Fecha límite: 30/10/2026');
  });
});
