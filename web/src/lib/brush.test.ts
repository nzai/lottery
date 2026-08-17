import { describe, expect, it } from 'vitest'
import { brushEnd, brushMove, brushStart } from './brush'

describe('brush 刷选状态机', () => {
  it('从未选中开始 → 勾选模式并选中该期', () => {
    const s = brushStart(new Set(), '2026094')
    expect(s.mode).toBe('select')
    expect(s.selected.has('2026094')).toBe(true)
  })

  it('从已选中开始 → 取消模式并取消该期', () => {
    const s = brushStart(new Set(['2026094']), '2026094')
    expect(s.mode).toBe('deselect')
    expect(s.selected.has('2026094')).toBe(false)
  })

  it('勾选模式下滑过 → 依次勾选', () => {
    let s = brushStart(new Set(), '2026094')
    s = brushMove(s, '2026093')
    s = brushMove(s, '2026092')
    expect(s.selected).toEqual(new Set(['2026094', '2026093', '2026092']))
  })

  it('取消模式下滑过 → 依次取消', () => {
    let s = brushStart(new Set(['2026094', '2026093', '2026092']), '2026094')
    s = brushMove(s, '2026093')
    s = brushMove(s, '2026092')
    expect(s.selected).toEqual(new Set())
  })

  it('重复划过同一期不产生变化', () => {
    let s = brushStart(new Set(), '2026094')
    s = brushMove(s, '2026094')
    expect(s.selected).toEqual(new Set(['2026094']))
  })

  it('brushEnd 清除手势状态但保留选择', () => {
    const s = brushEnd(brushStart(new Set(), '2026094'))
    expect(s.lastIssue).toBeNull()
    expect(s.selected.has('2026094')).toBe(true)
  })

  it('不修改原集合（不可变）', () => {
    const original = new Set(['2026094'])
    const s = brushStart(original, '2026093')
    expect(original.has('2026093')).toBe(false)
    expect(s.selected.has('2026093')).toBe(true)
  })
})
