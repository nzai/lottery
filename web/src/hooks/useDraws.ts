import { useCallback, useEffect, useState } from 'react'
import type { Draw } from '../types'

const PAGE_SIZE = 100

// 加载走势图数据：首屏 limit 期，滚动到底追加更早数据
export function useDraws(limit: number) {
  const [draws, setDraws] = useState<Draw[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [exhausted, setExhausted] = useState(false)

  // 期数范围变化 → 重新拉取并重置
  useEffect(() => {
    let cancelled = false
    setDraws([])
    setError(null)
    setExhausted(false)
    setLoading(true)
    fetch(`/api/draws?limit=${limit}`)
      .then((r) => (r.ok ? r.json() : Promise.reject(new Error(`HTTP ${r.status}`))))
      .then((data) => {
        if (cancelled) return
        setDraws(data.draws ?? []) // 空库时后端返回 []，此处兜底防御
        setExhausted((data.draws?.length ?? 0) < limit)
      })
      .catch((e) => {
        if (!cancelled) setError(String(e))
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })
    return () => {
      cancelled = true
    }
  }, [limit])

  // 滚动到底：追加更早数据
  const loadMore = useCallback(async () => {
    if (loading || exhausted || draws.length === 0) return
    setLoading(true)
    const before = draws[draws.length - 1].issue
    try {
      const r = await fetch(`/api/draws?limit=${PAGE_SIZE}&before=${before}`)
      if (!r.ok) throw new Error(`HTTP ${r.status}`)
      const data = await r.json()
      setDraws((prev) => [...prev, ...data.draws])
      setExhausted(data.draws.length < PAGE_SIZE)
    } catch (e) {
      setError(String(e))
    } finally {
      setLoading(false)
    }
  }, [loading, exhausted, draws])

  return { draws, loading, error, exhausted, loadMore }
}
