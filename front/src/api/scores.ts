import type { RoundScore } from '../types/api.ts'
import { get, post, put } from './client.ts'

export const getScores = (gameId: string | number, roundId: number) =>
  get<RoundScore[]>(`/games/${gameId}/rounds/${roundId}/scores`)
export const createScore = (gameId: string | number, roundId: number, playerId: number, score: number) =>
  post<RoundScore>(`/games/${gameId}/rounds/${roundId}/scores`, { player_id: playerId, score })
export const updateScore = (gameId: string | number, roundId: number, scoreId: number, score: number) =>
  put<RoundScore>(`/games/${gameId}/rounds/${roundId}/scores/${scoreId}`, { score })
