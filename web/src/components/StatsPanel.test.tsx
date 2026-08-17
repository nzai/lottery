import { render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { StatsPanel } from './StatsPanel'
import type { Stats } from '../lib/stats'

const stats: Stats = {
  n: 3,
  redFreq: [
    { number: 3, count: 3, rate: 1 },
    { number: 1, count: 1, rate: 1 / 3 },
  ],
  blueFreq: [
    { number: 1, count: 2, rate: 2 / 3 },
    { number: 2, count: 1, rate: 1 / 3 },
  ],
  redOmission: { 1: 2, 2: 1, 3: 0 },
  blueOmission: { 1: 0, 2: 1 },
  redHot: [3],
  redCold: [1],
  blueHot: [1, 2],
  blueCold: [],
  redOddEven: { odd: 9, even: 9 },
  redBigSmall: { big: 0, small: 18 },
  blueOddEven: { odd: 2, even: 1 },
}

describe('StatsPanel', () => {
  it('显示选中期数与清除按钮', () => {
    render(<StatsPanel stats={stats} selectedCount={3} onClear={() => {}} />)
    expect(screen.getByText(/已选 3 期/)).toBeTruthy()
    expect(screen.getByText('清除选择')).toBeTruthy()
  })

  it('渲染红球频次、遗漏与比例', () => {
    render(<StatsPanel stats={stats} selectedCount={3} onClear={() => {}} />)
    expect(screen.getByText('红球出现次数')).toBeTruthy()
    expect(screen.getByText('遗漏')).toBeTruthy()
    expect(screen.getByText('奇 9 : 偶 9')).toBeTruthy()
    expect(screen.getByText('大 0 : 小 18')).toBeTruthy()
  })

  it('清除按钮触发 onClear', () => {
    const onClear = vi.fn()
    render(<StatsPanel stats={stats} selectedCount={3} onClear={onClear} />)
    screen.getByText('清除选择').click()
    expect(onClear).toHaveBeenCalled()
  })

  it('stats 为 null 时隐藏内容', () => {
    render(<StatsPanel stats={null} selectedCount={3} onClear={() => {}} />)
    expect(screen.queryByText('红球出现次数')).toBeNull()
    expect(screen.getByText(/已选 3 期/)).toBeTruthy()
  })
})
