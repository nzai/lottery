import type { Draw } from '../types'

// 单个号码的频次统计
export interface FreqEntry {
  number: number
  count: number
  rate: number // 出现率 = count / n
}

export interface Ratio {
  odd: number
  even: number
}

export interface Stats {
  n: number // 选中期数
  redFreq: FreqEntry[] // 红球 1-33，按 count 降序（并列按号码升序）
  blueFreq: FreqEntry[] // 蓝球 1-16，同上
  redOmission: Record<number, number> // 号码 -> 遗漏期数（距选中范围最后一期）
  blueOmission: Record<number, number>
  redHot: number[]
  redWarm: number[]
  redCold: number[]
  blueHot: number[]
  blueWarm: number[]
  blueCold: number[]
  redOddEven: Ratio // 选中期内所有红球的奇偶计数
  redBigSmall: Ratio // 17-33 为大，1-16 为小
  blueOddEven: Ratio // 选中期内所有蓝球的奇偶计数
}

export function computeStats(draws: Draw[], selected: Set<string>): Stats {
  const list = draws
    .filter((d) => selected.has(d.issue))
    .sort((a, b) => (a.issue < b.issue ? 1 : -1)) // 最新在前
  const n = list.length

  const empty: Stats = {
    n: 0,
    redFreq: [], blueFreq: [],
    redOmission: {}, blueOmission: {},
    redHot: [], redWarm: [], redCold: [],
    blueHot: [], blueWarm: [], blueCold: [],
    redOddEven: { odd: 0, even: 0 },
    redBigSmall: { big: 0, small: 0 },
    blueOddEven: { odd: 0, even: 0 },
  }
  if (n === 0) return empty

  // 1. 出现次数
  const redCount = new Map<number, number>()
  const blueCount = new Map<number, number>()
  for (let i = 1; i <= 33; i++) redCount.set(i, 0)
  for (let i = 1; i <= 16; i++) blueCount.set(i, 0)
  for (const d of list) {
    for (const r of d.red) redCount.set(r, (redCount.get(r) ?? 0) + 1)
    blueCount.set(d.blue, (blueCount.get(d.blue) ?? 0) + 1)
  }
  const freq = (count: Map<number, number>, max: number): FreqEntry[] =>
    Array.from({ length: max }, (_, i) => i + 1)
      .map((number) => ({ number, count: count.get(number)!, rate: count.get(number)! / n }))
      .sort((a, b) => b.count - a.count || a.number - b.number)
  const redFreq = freq(redCount, 33)
  const blueFreq = freq(blueCount, 16)

  // 2. 遗漏值：距选中范围最新一期（list[0]）的未出现期数
  const redOmission: Record<number, number> = {}
  const blueOmission: Record<number, number> = {}
  for (let i = 1; i <= 33; i++) redOmission[i] = omissionOf(list, i, false)
  for (let i = 1; i <= 16; i++) blueOmission[i] = omissionOf(list, i, true)

  // 3. 冷热分档：count > 均值 为热，< 均值为冷，= 均值为温
  const buckets = (freqList: FreqEntry[], mean: number) => {
    const hot: number[] = []
    const warm: number[] = []
    const cold: number[] = []
    for (const e of freqList) {
      if (e.count > mean) hot.push(e.number)
      else if (e.count < mean) cold.push(e.number)
      else warm.push(e.number)
    }
    return { hot, warm, cold }
  }
  const red = buckets(redFreq, (6 * n) / 33)
  const blue = buckets(blueFreq, n / 16)

  // 4. 比例
  const ratio = (numbers: number[]): Ratio => {
    let odd = 0
    for (const v of numbers) if (v % 2 === 1) odd++
    return { odd, even: numbers.length - odd }
  }
  const allRed = list.flatMap((d) => d.red)
  let big = 0
  for (const v of allRed) if (v > 16) big++

  return {
    n,
    redFreq, blueFreq,
    redOmission, blueOmission,
    redHot: red.hot, redWarm: red.warm, redCold: red.cold,
    blueHot: blue.hot, blueWarm: blue.warm, blueCold: blue.cold,
    redOddEven: ratio(allRed),
    redBigSmall: { big, small: allRed.length - big },
    blueOddEven: ratio(list.map((d) => d.blue)),
  }
}

// 某号码在 list（最新在前）中的遗漏期数：第一个包含它的位置索引；从未出现则为 n
function omissionOf(list: Draw[], number: number, isBlue: boolean): number {
  for (let i = 0; i < list.length; i++) {
    const d = list[i]
    if (isBlue ? d.blue === number : d.red.includes(number)) return i
  }
  return list.length
}
