export interface Image {
  file: string
  assignedCategories: string[]
  proposedCategories: string[]
  starredCategory?: string
}

const state = () => ({
  data: [] as string[],
  error: undefined,
  loading: false,
})

export type State = ReturnType<typeof state>

export default state
