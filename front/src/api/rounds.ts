import type { Round } from '../types/api.ts'
import { get, post } from './client.ts'

export const getRounds = (gameId: string | number) => get<Round[]>(`/games/${gameId}/rounds`)
export const addRound = (gameId: string | number, number: number) =>
  post<Round>(`/games/${gameId}/rounds`, { number })
