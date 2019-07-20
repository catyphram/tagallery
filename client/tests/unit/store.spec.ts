import { getters, mutations, actions, state as defaultState } from '@/store';
import { Category, LIST_MODE, Image } from '@/models';
import * as api from '@/api/api';

jest.mock('@/api/api');

describe('store.ts', () => {
  describe('getters', () => {
    it('isCategorySelected should return true if a category is found in selectedCategories', () => {
      const state = {
        ...defaultState,
        selectedCategories: [
          { name: 'Category 1' }, { name: 'Category 2' }, { name: 'Category 3' },
        ] as Category[],
      };
      expect(getters.isCategorySelected(state)(state.selectedCategories[0])).toBe(true);
    });
  });

  describe('mutations', () => {
    it('updateCategories should set the categories and loading status in the state', () => {
      const categories = [{ name: 'Category 1' }, { name: 'Category 2' } ] as Category[];
      const state = {
        ...defaultState,
        categories: {
          data: [] as Category[],
          loading: true,
          error: undefined,
        },
      };
      mutations.updateCategories(state, {
        categories,
      });

      expect(state.categories.data).toBe(categories);
      expect(state.categories.loading).toBe(false);
      expect(state.categories.error).toBeUndefined();

      mutations.updateCategories(state, {
        categories: [],
        error: 'error',
      });

      expect(state.categories.error).toBe('error');
    });

    it('updateImages should set the images, completed and loading status in the state', () => {
      const images = [{ file: 'test1.jpg' }, { file: 'test2.jpg' } ] as Image[];
      const state = {
        ...defaultState,
        images: {
          data: [] as Image[],
          loading: true,
          completed: false,
        },
      };
      mutations.updateImages(state, {
        images,
        loading: false,
        completed: true,
      });

      expect(state.images.data).toBe(images);
      expect(state.images.loading).toBe(false);
      expect(state.images.completed).toBe(true);
    });

    it('updateImages should append the images if #append is true', () => {
      const images = [{ file: 'test1.jpg' }, { file: 'test2.jpg' } ] as Image[];
      const newImages = [{ file: 'test3.jpg' }, { file: 'test4.jpg' } ] as Image[];
      const state = {
        ...defaultState,
        images: {
          completed: false,
          data: [...images],
          loading: true,
        },
      };
      mutations.updateImages(state, {
        images: newImages,
        append: true,
      });

      expect(state.images.data).toEqual(images.concat(newImages));
    });

    it('selectImage should select the image', () => {
      const state = {
        ...defaultState,
      };
      mutations.selectImage(state, { index: 1 });

      expect(state.selectedImage).toBe(1);
    });

    it('setMode should set the mode in the state', () => {
      const state = {
        ...defaultState,
        listMode: LIST_MODE.MODE_VIEW,
      };
      mutations.setMode(state, { mode: LIST_MODE.MODE_VERIFY });

      expect(state.listMode).toBe(LIST_MODE.MODE_VERIFY);
    });

    it('selectCategory should add the category in the state and unset #listUncategorized', () => {
      const category = { name: 'Category 1' } as Category;
      const state = {
        ...defaultState,
        selectedCategories: [] as Category[],
        listUncategorized: true,
      };

      mutations.selectCategory(state, { category });
      expect(state.selectedCategories).toEqual([category ]);
      expect(state.listUncategorized).toBe(false);

      mutations.selectCategory(state, { category });
      expect(state.selectedCategories).toEqual([category ]);
    });

    it('unselectCategory should remove a category if it is selected', () => {
      const cat1 = { name: 'Category 1' } as Category;
      const cat2 = { name: 'Category 2' } as Category;
      const cat3 = { name: 'Category 3' } as Category;
      const categories = [cat1, cat2, cat3 ];
      const state = {
        ...defaultState,
        selectedCategories: categories,
      };

      mutations.unselectCategory(state, { category: { name: 'Category 0' } as Category });
      expect(state.selectedCategories).toEqual([cat1, cat2, cat3 ]);

      mutations.unselectCategory(state, { category: cat2 });
      expect(state.selectedCategories).toEqual([cat1, cat3 ]);
    });

    it('setListUncategorized should change the flag in the state and remove selected categories if true', () => {
      const state = {
        selectedCategories: [
          { name: 'Category 1' }, { name: 'Category 2' }, { name: 'Category 3' },
        ] as Category[],
        listUncategorized: false,
      };

      mutations.setListUncategorized(state, { flag: true });
      expect(state.selectedCategories).toEqual([]);
      expect(state.listUncategorized).toBe(true);

      mutations.setListUncategorized(state, { flag: false });
      expect(state.listUncategorized).toBe(false);
    });
  });

  describe('actions', () => {
    let context;

    beforeEach(() => {
      context = {
        commit: jest.fn(),
        dispatch: jest.fn(),
      };
    });

    it('loadCategories should load the categories', async () => {
      await actions.loadCategories(context);
      expect(context.commit).toMatchSnapshot();
      (api.loadCategories as jest.Mock).mockImplementationOnce(() => {
        return Promise.reject('error');
      });
      await actions.loadCategories(context);
      expect(context.commit).toMatchSnapshot();

    });

    it('loadImages should load a batch of images', async () => {
      await actions.loadImages(context, true);
      expect(context.commit).toMatchSnapshot();
      await actions.loadImages(context);
      expect(context.commit).toMatchSnapshot();
    });

    it('setMode should set the viewing mode', () => {
      actions.setMode(context, LIST_MODE.MODE_VERIFY);
      expect(context.commit).toHaveBeenCalledWith('setMode', { mode: LIST_MODE.MODE_VERIFY });
    });

    it('toggleCategory should select a category if it is not selected yet', () => {
      const category = { name: 'Category 1', key: 'cat1' } as Category;
      context.state = {
        selectedCategories: [],
      };
      actions.toggleCategory(context, category);
      expect(context.commit).toHaveBeenCalledWith('selectCategory', { category });
    });

    it('toggleCategory should unselect a category if it is selected already', () => {
      const category = { name: 'Category 1', key: 'cat1' } as Category;
      context.state = {
        selectedCategories: [category ],
      };
      actions.toggleCategory(context, category);
      expect(context.commit).toHaveBeenCalledWith('unselectCategory', { category });
    });

    it('toggleListUncategorized should toggle the flag', () => {
      context.state = {
        listUncategorized: false,
      };
      actions.toggleListUncategorized(context);
      expect(context.commit).toHaveBeenCalledWith('setListUncategorized', { flag: true });
    });
  });
});
