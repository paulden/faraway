import type { Player } from '../types/api.ts'
import { get, post, del } from './client.ts'

export const getPlayers = (gameId: string | number) => get<Player[]>(`/games/${gameId}/players`)
export const addPlayer = (gameId: string | number, name: string) =>
  post<Player>(`/games/${gameId}/players`, { name })
export const removePlayer = (gameId: string | number, playerId: number) =>
  del(`/games/${gameId}/players/${playerId}`)
