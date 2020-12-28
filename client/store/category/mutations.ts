import { MutationTree } from 'vuex'

import { Category, State } from './state'

const mutations: MutationTree<State> = {
  load(state: State) {
    state.loading = true
  },
  receive(state: State, { data = [] as Category[], error = undefined } = {}) {
    state.loading = false
    state.data = data
    state.error = error
  },
  update(state: State, { category }: { category: Category }) {
    const stateCategory = state.data.find(
      (stateCategory) => stateCategory.id === category.id
    )
    if (stateCategory) {
      stateCategory.name = category.name
      stateCategory.description = category.description
    }
  },
  add(state: State, { category }: { category: Category }) {
    state.data.push(category)
  },
}

export default mutations
