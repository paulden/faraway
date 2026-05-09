import type { Game } from '../types/api.ts'
import { get, post, put } from './client.ts'

export const getGames = () => get<Game[]>('/games')
export const getGame = (id: string | number) => get<Game>(`/games/${id}`)
export const createGame = (title: string) => post<Game>('/games', { title, is_finished: false })
export const updateGame = (id: string | number, title: string, isFinished: boolean) =>
  put<Game>(`/games/${id}`, { title, is_finished: isFinished })
