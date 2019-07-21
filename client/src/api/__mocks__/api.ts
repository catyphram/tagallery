import { Category, Image } from '@/models';

export const loadCategories = jest.fn(async (): Promise<Category[]> => {
  return await [
    { name: 'Category 1', key: 'cat1', description: 'Category 1' },
    { name: 'Category 2', key: 'cat2', description: 'Category 2' },
    { name: 'Category 3', key: 'cat3', description: 'Category 3' },
    { name: 'Category 4', key: 'cat4', description: 'Category 4' },
    { name: 'Category 5', key: 'cat5', description: 'Category 5' },
    { name: 'Category 6', key: 'cat6', description: 'Category 6' },
  ];
});

export const loadImages = jest.fn(async (): Promise<Image[]> => {
  return await [
    { file: 'image0.jpg' },
    { file: 'image1.jpg' },
    { file: 'image2.jpg' },
    { file: 'image3.jpg' },
    { file: 'image4.jpg' },
    { file: 'image5.jpg' },
    { file: 'image6.jpg' },
    { file: 'image7.jpg' },
    { file: 'image8.jpg' },
    { file: 'image9.jpg' },
  ] as Image[];
});

export const updateImage = jest.fn(async (image): Promise<boolean> => {
  return await true;
});
