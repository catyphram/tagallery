import { MutationTree } from 'vuex'

import { Image, State } from './state'

const mutations: MutationTree<State> = {
  load: (state: State) => {
    state.loading = true
  },
  receive: (state: State, { data = [] as Image[], error = undefined } = {}) => {
    state.loading = false
    state.data = data
    state.error = error
  },
}

export default mutations
