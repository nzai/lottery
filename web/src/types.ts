// 一期开奖结果（与后端 /api/draws 返回结构一致）
export interface Draw {
  issue: string
  date: string
  red: number[]
  blue: number
}
