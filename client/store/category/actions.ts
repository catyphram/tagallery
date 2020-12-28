import { ActionTree } from 'vuex'

import { Category, State } from './state'

const mockCategories = [
  {
    id: 'category1',
    name: 'Category 1',
    description: 'Category 1 Description',
  },
  {
    id: 'category2',
    name: 'Category 2',
    description: 'Category 2 Description',
  },
  {
    id: 'category3',
    name: 'Category 3',
    description: 'Category 3 Description',
  },
  {
    id: 'category4',
    name: 'Category 4',
    description: 'Category 4 Description',
  },
  {
    id: 'category5',
    name: 'Category 5',
    description: 'Category 5 Description',
  },
]

const actions: ActionTree<State, State> = {
  load({ commit }) {
    commit('load')
    commit('receive', {
      data: mockCategories,
    })
  },
  update({ commit }, { category }: { category: Category }) {
    commit('update', { category })
  },
  add({ commit }, { category }: { category: Category }) {
    commit('add', { category })
  },
}

export default actions
