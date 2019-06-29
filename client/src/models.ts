export enum LIST_MODE {
    MODE_VIEW,
    MODE_VERIFY,
}

export interface Category {
    name: string;
    key: string;
    description?: string;
}

export interface Image {
  file: string;
  assignedCategories?: string[];
  proposedCategories?: string[];
  starredCategory?: string;
}

export interface Config {
  api: string;
}
