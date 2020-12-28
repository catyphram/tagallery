import { ActionTree } from 'vuex'
import { Category } from '../category/state'

import { Categorization, Mode, State, Status } from './state'

const actions: ActionTree<State, State> = {
  set(
    { commit },
    {
      status,
      mode,
      categorization,
      categories,
    }: {
      status: Status
      mode: Mode
      categorization: Categorization
      categories: Category[]
    }
  ) {
    commit('set', {
      status,
      mode,
      categorization,
      categories,
    })
  },
}

export default actions
