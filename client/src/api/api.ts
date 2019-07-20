import { Category, Image } from '@/models';

import config from '../../config.json';

/**
 * loadCategories loads the categories from the API.
 */
export const loadCategories = async (): Promise<Category[]> => {
  return fetch(`${config.api}/categories`).then((response) => {
      return response.json() as Promise<Category[]>;
  });
};

// Additional parameter for the local stub only so the URL differs on reload
// and Vue updates the images in the HTML.
let counter = 0;

/**
 * loadImages loads a batch of images from the API.
 * @todo: Add parameters for categories, mode and lastImage filter
 * Also don't forget to cancel the fetch request when another request is sent.
 */
export const loadImages = async (): Promise<Image[]> => {
  const stubImages = [] as Image[];
  for (let i = 0; i < 30; i++) {
    stubImages.push({
      file: `https://picsum.photos/400/400/?random&key=${counter}-${i}`,
    } as Image);
  }

  counter++;

  return await stubImages;
};
