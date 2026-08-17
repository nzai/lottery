import { act, renderHook, waitFor } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { useDraws } from './useDraws'
import type { Draw } from '../types'

const mkDraw = (issue: string): Draw => ({
  issue,
  date: '2026-08-13',
  red: [1, 2, 3, 4, 5, 6],
  blue: 1,
})

describe('useDraws', () => {
  let fetchMock: ReturnType<typeof vi.fn>

  beforeEach(() => {
    fetchMock = vi.fn()
    vi.stubGlobal('fetch', fetchMock)
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('首次加载返回数据', async () => {
    fetchMock.mockResolvedValue({
      ok: true,
      json: async () => ({ draws: [mkDraw('2026094'), mkDraw('2026093')] }),
    })
    const { result } = renderHook(() => useDraws(100))
    await waitFor(() => expect(result.current.draws.length).toBe(2))
    expect(fetchMock).toHaveBeenCalledWith('/api/draws?limit=100')
  })

  it('加载失败时记录 error', async () => {
    fetchMock.mockResolvedValue({ ok: false, status: 500 })
    const { result } = renderHook(() => useDraws(100))
    await waitFor(() => expect(result.current.error).toBeTruthy())
  })

  it('loadMore 追加更早数据，不足一页时标记 exhausted', async () => {
    const page1 = Array.from({ length: 100 }, (_, i) => mkDraw(`2026${String(100 - i).padStart(3, '0')}`))
    const page2 = Array.from({ length: 30 }, (_, i) => mkDraw(`2025${String(70 - i).padStart(3, '0')}`))
    fetchMock
      .mockResolvedValueOnce({ ok: true, json: async () => ({ draws: page1 }) })
      .mockResolvedValueOnce({ ok: true, json: async () => ({ draws: page2 }) })

    const { result } = renderHook(() => useDraws(100))
    await waitFor(() => expect(result.current.draws.length).toBe(100))

    await act(async () => {
      await result.current.loadMore()
    })
    expect(result.current.draws.length).toBe(130)
    expect(result.current.exhausted).toBe(true)
    // before 参数 = 当前最早期号（page1 末位 2026001）
    expect(fetchMock).toHaveBeenLastCalledWith('/api/draws?limit=100&before=2026001')
  })
})
