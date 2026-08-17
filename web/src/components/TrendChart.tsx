import { useCallback, useLayoutEffect, useMemo, useRef } from 'react'
import type { Draw } from '../types'
import type { BrushHandlers } from './types'

const RED_NUMBERS = Array.from({ length: 33 }, (_, i) => i + 1)
const BLUE_NUMBERS = Array.from({ length: 16 }, (_, i) => i + 1)

interface TrendChartProps {
  draws: Draw[]
  selected: Set<string>
  fontScale: number
  brush: BrushHandlers
  onLoadMore: () => void
}

// 走势图：左侧固定列（复选框+期号+日期），中间 33 红球 + 16 蓝球，可横向滚动。
// 数据层为倒序（最新在前），渲染时反转使最新一期在底部（时间正序，从上往下看历史）。
export function TrendChart({ draws, selected, fontScale, brush, onLoadMore }: TrendChartProps) {
  const scrollRef = useRef<HTMLDivElement>(null)
  const loadMoreRef = useRef(onLoadMore)
  loadMoreRef.current = onLoadMore
  const initializedRef = useRef(false)
  const lastScrollHeightRef = useRef(0)

  // 最新一期在底部（时间正序）
  const displayedDraws = useMemo(() => [...draws].reverse(), [draws])

  // 数据变化时的滚动位置处理：
  // - 首次加载/范围切换：定位到底部（最新一期）
  // - 顶部追加更早数据：补偿 scrollTop，保持视口内容不跳动
  useLayoutEffect(() => {
    const el = scrollRef.current
    if (!el) return
    if (draws.length === 0) {
      initializedRef.current = false
      lastScrollHeightRef.current = 0
      return
    }
    if (!initializedRef.current) {
      el.scrollTop = el.scrollHeight
      initializedRef.current = true
    } else if (lastScrollHeightRef.current > 0 && el.scrollHeight > lastScrollHeightRef.current) {
      el.scrollTop += el.scrollHeight - lastScrollHeightRef.current
    }
    lastScrollHeightRef.current = el.scrollHeight
  }, [draws])

  // 纵向滚动接近顶部时加载更早数据（正序下更早的在上面）
  const handleScroll = useCallback(() => {
    const el = scrollRef.current
    if (!el) return
    if (el.scrollTop <= 200) {
      loadMoreRef.current()
    }
  }, [])

  return (
    <div
      className="chart-scroll"
      ref={scrollRef}
      onScroll={handleScroll}
      style={{ '--font': `${16 * fontScale}px` } as React.CSSProperties}
    >
      {/* 表头：与行同列结构，sticky top */}
      <div className="chart-header">
        <span className="col-check">选</span>
        <span className="col-issue">期号</span>
        <span className="col-date">日期</span>
        {RED_NUMBERS.map((n) => (
          <span key={`r${n}`} className="col-num red">{String(n).padStart(2, '0')}</span>
        ))}
        {BLUE_NUMBERS.map((n) => (
          <span key={`b${n}`} className="col-num blue">{String(n).padStart(2, '0')}</span>
        ))}
      </div>

      {/* 数据行（时间正序：最早在上，最新在底部） */}
      <div className="chart-body">
        {displayedDraws.map((d) => {
          const isSelected = selected.has(d.issue)
          return (
            <div
              key={d.issue}
              data-issue={d.issue}
              className={`trend-row${isSelected ? ' selected' : ''}`}
            >
              <span
                className="col-check"
                onPointerDown={(e) => brush.handlePointerDown(e, d.issue)}
                onPointerMove={brush.handlePointerMove}
                onPointerUp={brush.handlePointerUp}
              >
                <span className={`checkbox${isSelected ? ' checked' : ''}`} />
              </span>
              <span className="col-issue">{d.issue}</span>
              <span className="col-date">{d.date.slice(5)}</span>
              {RED_NUMBERS.map((n) => (
                <span key={n} className="cell">
                  {d.red.includes(n) && <span className="ball red">{String(n).padStart(2, '0')}</span>}
                </span>
              ))}
              {BLUE_NUMBERS.map((n) => (
                <span key={n} className="cell">
                  {d.blue === n && <span className="ball blue">{String(n).padStart(2, '0')}</span>}
                </span>
              ))}
            </div>
          )
        })}
        {draws.length === 0 && <div className="chart-empty">加载中…</div>}
      </div>
    </div>
  )
}
