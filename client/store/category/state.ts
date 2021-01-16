export interface Category {
  id: string
  name: string
  description?: string
}

const state = () => ({
  data: [] as Category[],
  error: undefined,
  loading: false,
})

export type State = ReturnType<typeof state>

export default state
