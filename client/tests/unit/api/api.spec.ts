import { loadCategories, loadImages } from '@/api/api';
import { loadCategories as mockloadCategories } from '@/api/__mocks__/api';
import config from '../../../config.json';

describe('api.ts', () => {
  it('loadCategories should load and return the categories', async () => {
    window.fetch = jest.fn(() => {
      return Promise.resolve({
        json: mockloadCategories,
      });
    });

    expect(await loadCategories()).toEqual(await mockloadCategories());
    expect(window.fetch).toHaveBeenCalledWith(`${config.api}/categories`);
  });

  it('loadImages should load and return a batch of images', async () => {
    expect(await loadImages()).toMatchSnapshot();
  });
});
