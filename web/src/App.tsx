import { useEffect, useMemo, useState } from 'react'
import { StatsPanel } from './components/StatsPanel'
import { Toolbar } from './components/Toolbar'
import { TrendChart } from './components/TrendChart'
import type { BrushHandlers } from './components/types'
import { useBrush } from './hooks/useBrush'
import { useDraws } from './hooks/useDraws'
import { computeStats, type Stats } from './lib/stats'

const RANGES = [20, 50, 100, 200, 500]
const FONT_MIN = 0.8
const FONT_MAX = 1.5

export default function App() {
  const [range, setRange] = useState(100)
  const [fontScale, setFontScale] = useState(() => {
    const saved = localStorage.getItem('font-scale')
    const n = saved ? Number(saved) : 1
    return Number.isFinite(n) ? Math.min(FONT_MAX, Math.max(FONT_MIN, n)) : 1
  })
  const { draws, loading, error, loadMore } = useDraws(range)

  // 服务端编译版本（编译时间），用于观察更新是否生效
  const [version, setVersion] = useState<string | null>(null)
  useEffect(() => {
    let cancelled = false
    fetch('/api/version')
      .then((r) => (r.ok ? r.json() : null))
      .then((data) => {
        if (!cancelled && data) setVersion(data.version)
      })
      .catch(() => {})
    return () => {
      cancelled = true
    }
  }, [])

  // 选择集合：刷选即时更新（行高亮），统计 300ms debounce 后计算
  const [selected, setSelected] = useState<Set<string>>(new Set())
  const [stats, setStats] = useState<Stats | null>(null)

  const brush: BrushHandlers = useBrush((s) => setSelected(s))

  useEffect(() => {
    if (selected.size === 0) {
      setStats(null)
      return
    }
    const t = setTimeout(() => setStats(computeStats(draws, selected)), 300)
    return () => clearTimeout(t)
  }, [selected, draws])

  const changeFont = (delta: number) => {
    setFontScale((s) => {
      const next = Math.min(FONT_MAX, Math.max(FONT_MIN, Math.round((s + delta) * 10) / 10))
      localStorage.setItem('font-scale', String(next))
      return next
    })
  }

  const selectedCount = useMemo(() => selected.size, [selected])

  return (
    <div className="app">
      <Toolbar
        range={range}
        ranges={RANGES}
        fontScale={fontScale}
        onRangeChange={(r) => {
          setRange(r)
          brush.clear()
        }}
        onFontChange={changeFont}
      />

      {error && (
        <div className="error-banner">
          <span>数据加载失败：{error}</span>
          <button onClick={() => window.location.reload()}>重试</button>
        </div>
      )}

      <TrendChart
        draws={draws}
        selected={selected}
        fontScale={fontScale}
        brush={brush}
        onLoadMore={loadMore}
      />
      {loading && draws.length === 0 && <div className="chart-empty">加载中…</div>}

      {selectedCount > 0 && (
        <StatsPanel stats={stats} selectedCount={selectedCount} onClear={brush.clear} />
      )}

      <footer className="app-footer">
        {version ? `版本 ${version}` : ''}
      </footer>
    </div>
  )
}
