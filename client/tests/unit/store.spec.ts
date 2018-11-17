import { getters, mutations, actions } from '@/store';
import { Category, LIST_MODE } from '@/models';

jest.mock('@/api/api');

describe('store.ts', () => {
  describe('getters', () => {
    it('isCategorySelected should return true if a category is found in selectedCategories', () => {
      const state = {
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
        categories: {
          data: [] as Category[],
          loading: false,
        },
      };
      mutations.updateCategories(state, {
        categories,
        loading: false,
      });

      expect(state.categories.data).toBe(categories);
      expect(state.categories.loading).toBe(false);
    });

    it('setMode should set the mode in the state', () => {
      const state = {
        listMode: LIST_MODE.MODE_VIEW,
      };
      mutations.setMode(state, { mode: LIST_MODE.MODE_VERIFY });

      expect(state.listMode).toBe(LIST_MODE.MODE_VERIFY);
    });

    it('selectCategory should add the category in the state and unset #listUncategorized', () => {
      const category = { name: 'Category 1' } as Category;
      const state = {
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
      const state = { selectedCategories: categories };

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
      };
    });

    it('loadCategories should load and set the categories', async () => {
      await actions.loadCategories(context);
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
