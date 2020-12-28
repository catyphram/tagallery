import { MutationTree } from 'vuex'
import { Category } from '../category/state'

import { Categorization, Mode, State, Status } from './state'

const mutations: MutationTree<State> = {
  set(
    state: State,
    {
      status,
      mode,
      categorization,
      categories,
    }: {
      status?: Status
      mode?: Mode
      categorization?: Categorization
      categories?: Category[]
    }
  ) {
    state.status = status ?? state.status
    state.mode = mode ?? state.mode
    state.categorization = categorization ?? state.categorization
    state.categories = categories ?? state.categories
  },
}

export default mutations
