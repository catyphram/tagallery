import type { Category } from '../category/state'

export enum Status {
  Processed,
  Unprocessed,
}
export enum Mode {
  View,
  Verify,
}
export enum Categorization {
  Done,
  Open,
}

const state = () => ({
  status: Status.Processed,
  mode: Mode.View,
  categorization: Categorization.Done,
  categories: [] as Category[],
})

export type State = ReturnType<typeof state>

export default state
