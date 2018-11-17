import { Category } from '@/models';

export const loadCategories = async (): Promise<Category[]> => {
    return await [
        { name: 'Category 1', key: 'cat1', description: 'Category 1' },
        { name: 'Category 2', key: 'cat2', description: 'Category 2' },
        { name: 'Category 3', key: 'cat3', description: 'Category 3' },
        { name: 'Category 4', key: 'cat4', description: 'Category 4' },
        { name: 'Category 5', key: 'cat5', description: 'Category 5' },
        { name: 'Category 6', key: 'cat6', description: 'Category 6' },
    ];
};
