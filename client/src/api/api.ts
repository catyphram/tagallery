import { Category } from '@/models';

export const loadCategories = async (): Promise<Category[]> => {
    return await [
        { name: 'Category 1', key: 'cat1', description: 'Category 1' },
        { name: 'Category 2', key: 'cat2', description: 'Category 2' },
    ];
};
