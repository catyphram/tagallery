import { loadCategories } from '@/api/api';

describe('api.ts', () => {
  it('loadCategories should load and return the categories', async () => {
    expect(await loadCategories()).toMatchSnapshot();
  });
});
