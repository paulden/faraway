export interface Game {
  id: number
  title: string
  is_finished: boolean
  created_at: string
  updated_at: string
}

export interface Player {
  id: number
  name: string
  game_id: number
  created_at: string
  updated_at: string
}

export interface Round {
  id: number
  number: number
  game_id: number
  created_at: string
  updated_at: string
}

export interface RoundScore {
  id: number
  score: number
  round_id: number
  player_id: number
  created_at: string
  updated_at: string
}
