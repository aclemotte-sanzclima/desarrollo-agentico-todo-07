import type { Page } from '@playwright/test';
import { test, expect } from '../fixtures';

// Comprueba la fecha límite de cada fila de la lista, en orden de creación.
async function expectDueDates(page: Page) {
  const rows = page.getByRole('list').getByRole('listitem');
  await expect(page.locator('.task-title')).toHaveText(['Revisar contrato', 'Pagar la cuota', 'Regar las plantas']);

  // Pendiente con fecha
  await expect(rows.nth(0).locator('.task-due-date')).toHaveText('Fecha límite: 05/11/2026');
  await expect(rows.nth(0).locator('.task-due-date')).toHaveAttribute('datetime', '2026-11-05');

  // Completada con fecha
  await expect(rows.nth(1).getByRole('button', { name: 'Reabrir tarea: Pagar la cuota' })).toBeVisible();
  await expect(rows.nth(1).locator('.task-due-date')).toBeVisible();
  await expect(rows.nth(1).locator('.task-due-date')).toHaveText('Fecha límite: 30/10/2026');
  await expect(rows.nth(1).locator('.task-due-date')).toHaveAttribute('datetime', '2026-10-30');

  // Sin fecha: ni la fecha ni una etiqueta vacía
  await expect(rows.nth(2).locator('.task-due-date')).toHaveCount(0);
  await expect(rows.nth(2).locator('time')).toHaveCount(0);

  await expect(page.locator('.task-due-date')).toHaveCount(2);
}

test.describe('Tareas', () => {
  test('la lista muestra la fecha límite de cada tarea', async ({ page, tasks }) => {
    // 1. Tres tareas: pendiente con fecha, completada con fecha y sin fecha
    await tasks([
      { title: 'Revisar contrato', dueDate: '2026-11-05' },
      { title: 'Pagar la cuota', dueDate: '2026-10-30', completed: true },
      { title: 'Regar las plantas' },
    ]);
    await page.reload();

    // 2. Cada tarea muestra su fecha límite, o ninguna
    await expectDueDates(page);

    // 3. Igual tras recargar
    await page.reload();
    await expectDueDates(page);
  });
});
