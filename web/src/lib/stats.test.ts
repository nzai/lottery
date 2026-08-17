import { describe, expect, it } from 'vitest'
import { computeStats } from './stats'
import type { Draw } from '../types'

// 构造测试数据：3 期，红球 1-6 / 2-7 / 3-8，蓝球 1 / 2 / 1
const draws: Draw[] = [
  { issue: '2026091', date: '2026-08-06', red: [1, 2, 3, 4, 5, 6], blue: 1 },
  { issue: '2026092', date: '2026-08-09', red: [2, 3, 4, 5, 6, 7], blue: 2 },
  { issue: '2026093', date: '2026-08-13', red: [3, 4, 5, 6, 7, 8], blue: 1 },
]

const all = new Set(draws.map((d) => d.issue))

describe('computeStats 基础', () => {
  it('空选择返回 n=0 的空统计', () => {
    const s = computeStats(draws, new Set())
    expect(s.n).toBe(0)
    expect(s.redFreq).toEqual([])
    expect(s.redHot).toEqual([])
  })

  it('统计出现次数并按次数降序（并列按号码升序）', () => {
    const s = computeStats(draws, all)
    expect(s.n).toBe(3)
    // 红球 3,4,5,6 各出现 3 次 → 排最前
    expect(s.redFreq.slice(0, 4).map((e) => e.number)).toEqual([3, 4, 5, 6])
    expect(s.redFreq[0].count).toBe(3)
    expect(s.redFreq[0].rate).toBeCloseTo(1)
    // 红球 1 出现 1 次
    const one = s.redFreq.find((e) => e.number === 1)
    expect(one?.count).toBe(1)
    expect(one?.rate).toBeCloseTo(1 / 3)
    // 蓝球 1 出现 2 次 > 蓝球 2 出现 1 次
    expect(s.blueFreq[0].number).toBe(1)
    expect(s.blueFreq[0].count).toBe(2)
    expect(s.blueFreq[1].number).toBe(2)
  })

  it('遗漏值：最新一期开出的号码遗漏 0，未出现的为 n', () => {
    const s = computeStats(draws, all)
    // 最新 2026093 开出 3,4,5,6,7,8 → 遗漏 0
    expect(s.redOmission[3]).toBe(0)
    expect(s.redOmission[8]).toBe(0)
    // 号码 1 在 2026091 开出（距最新 2 期）→ 遗漏 2
    expect(s.redOmission[1]).toBe(2)
    // 号码 9 从未开出 → 遗漏 n=3
    expect(s.redOmission[9]).toBe(3)
    // 蓝球 1 在最新开出 → 0；蓝球 2 在中间开出 → 1
    expect(s.blueOmission[1]).toBe(0)
    expect(s.blueOmission[2]).toBe(1)
  })

  it('冷热分档：>均值热、<均值冷、=均值温', () => {
    const s = computeStats(draws, all)
    // 红球均值 = 6*3/33 = 0.545 → 出现 1 次即为热
    expect(s.redHot).toContain(3)
    expect(s.redCold).toContain(9)
    expect(s.redCold).toContain(33)
    // 蓝球均值 = 3/16 = 0.1875 → 出现 1 次即热
    expect(s.blueHot).toContain(2)
  })

  it('比例：红球奇偶、大小、蓝球奇偶', () => {
    const s = computeStats(draws, all)
    // 全部红球 18 个：1-8 各出现 3 次。奇数 1,3,5,7 → 3*3=9；偶数 2,4,6,8 → 3*3=9
    expect(s.redOddEven).toEqual({ odd: 9, even: 9 })
    // 大小：1-8 均 ≤16 → 全小
    expect(s.redBigSmall).toEqual({ big: 0, small: 18 })
    // 蓝球：1,2,1 → 奇 2 偶 1
    expect(s.blueOddEven).toEqual({ odd: 2, even: 1 })
  })
})

describe('computeStats 局部选择', () => {
  it('只统计选中的期', () => {
    const s = computeStats(draws, new Set(['2026091', '2026093']))
    expect(s.n).toBe(2)
    // 红球 3,4,5,6 各 2 次
    expect(s.redFreq[0].count).toBe(2)
    // 号码 1：只在 2026091 出现，距最新选中期（2026093）遗漏 1
    expect(s.redOmission[1]).toBe(1)
    // 蓝球 2 未选中 → 遗漏 n=2
    expect(s.blueOmission[2]).toBe(2)
  })
})
