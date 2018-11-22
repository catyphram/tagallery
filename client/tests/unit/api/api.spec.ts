import { loadCategories, loadImages } from '@/api/api';

describe('api.ts', () => {
  it('loadCategories should load and return the categories', async () => {
    expect(await loadCategories()).toMatchSnapshot();
  });

  it('loadImages should load and return a batch of images', async () => {
    expect(await loadImages()).toMatchSnapshot();
  });
});
