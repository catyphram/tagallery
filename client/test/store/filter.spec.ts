import { Category } from '~/store/category/state'
import mutations from '~/store/filter/mutations'
import { Categorization, Mode, State, Status } from '~/store/filter/state'

const getDefaultState = (): State => ({
  status: Status.Processed,
  mode: Mode.View,
  categorization: Categorization.Done,
  categories: [] as Category[],
})

describe('filter', () => {
  describe('mutation', () => {
    describe('set', () => {
      test('should set new options', () => {
        const state = getDefaultState()
        const updates = {
          status: Status.Unprocessed,
          mode: Mode.Verify,
          categorization: Categorization.Open,
          categories: [
            {
              id: 'cat1',
              name: 'Category 1',
              description: 'Category 1 Description',
            },
          ],
        }

        mutations.set(state, updates)

        expect(state).toMatchSnapshot()
      })
    })
  })
})
