import { ActionTree } from 'vuex'
import { Category } from '../category/state'

import { Categorization, Mode, State as Filter, Status } from '../filter/state'
import { Image, State } from './state'

type ImageParams = {
  status: 'unprocessed' | 'uncategorized' | 'autocategorized' | 'categorized'
  count: number
  lastImage?: string
  categories?: Category[]
}

let counter = 0

const getNextImage = () => {
  return `https://picsum.photos/id/${++counter}/200/300`
}

const getNextImageSet = (amount: number) => {
  const images = [] as Image[]
  for (let i = 0; i < amount; i++) {
    images.push({
      file: getNextImage(),
      assignedCategories: [],
      proposedCategories: [],
    })
  }
  return images
}
const actions: ActionTree<State, State> = {
  load(
    { state, commit },
    {
      append = true,
      filter = {
        status: Status.Processed,
        mode: Mode.View,
        categorization: Categorization.Done,
        categories: [],
      },
      count = 15,
      lastImage,
    }: {
      append?: boolean
      filter?: Filter
      count?: number
      lastImage?: string
    } = {}
  ) {
    // TODO: Cancel current request and send a new one when already loading.
    if (!state.loading) {
      // TODO: Remove dummy counter
      if (state.data.length > counter) counter = state.data.length

      const params: ImageParams = {
        count,
        status: 'unprocessed',
      }

      if (filter.status === Status.Processed) {
        if (filter.categorization === Categorization.Done) {
          if (filter.mode === Mode.View) {
            params.status = 'categorized'
          } else {
            params.status = 'autocategorized'
          }
        } else {
          params.status = 'uncategorized'
        }
      }

      if (lastImage) params.lastImage = lastImage
      if (filter.categories.length) params.categories = filter.categories

      commit('load')
      commit('receive', {
        data: append
          ? [...state.data, ...getNextImageSet(15)]
          : getNextImageSet(15),
      })
    }
  },
}

export default actions
