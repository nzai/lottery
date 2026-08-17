import type { Stats } from '../lib/stats'

interface StatsPanelProps {
  stats: Stats | null
  selectedCount: number
  onClear: () => void
}

// 号码 + 次数 + 条形图一行
function FreqBar({ number, count, rate, max, color }: { number: number; count: number; rate: number; max: number; color: 'red' | 'blue' }) {
  return (
    <div className="freq-row">
      <span className={`freq-ball ${color}`}>{String(number).padStart(2, '0')}</span>
      <div className="freq-bar-track">
        <div className={`freq-bar ${color}`} style={{ width: `${(count / max) * 100}%` }} />
      </div>
      <span className="freq-count">
        {count} 次 ({Math.round(rate * 100)}%)
      </span>
    </div>
  )
}

// 号码徽章组（冷热 / 遗漏）
function Chips({ items, color }: { items: { number: number; text: string }[]; color: string }) {
  if (items.length === 0) return <span className="chips-empty">无</span>
  return (
    <div className="chips">
      {items.map(({ number, text }) => (
        <span key={number} className={`chip ${color}`}>
          {String(number).padStart(2, '0')}
          <em>{text}</em>
        </span>
      ))}
    </div>
  )
}

export function StatsPanel({ stats, selectedCount, onClear }: StatsPanelProps) {
  return (
    <div className="stats-bar">
      <div className="stats-bar-head">
        <span>已选 {selectedCount} 期</span>
        <button className="clear-btn" onClick={onClear}>
          清除选择
        </button>
      </div>
      {stats && (
        <div className="stats-body">
          <section className="stats-section">
            <h3>红球出现次数</h3>
            {stats.redFreq.map((e) => (
              <FreqBar key={e.number} number={e.number} count={e.count} rate={e.rate} max={stats.redFreq[0].count} color="red" />
            ))}
          </section>

          <section className="stats-section">
            <h3>蓝球出现次数</h3>
            {stats.blueFreq.map((e) => (
              <FreqBar key={e.number} number={e.number} count={e.count} rate={e.rate} max={stats.blueFreq[0].count} color="blue" />
            ))}
          </section>

          <section className="stats-section">
            <h3>冷热号</h3>
            <div className="stats-line">
              <span className="stat-label">红球热</span>
              <Chips items={stats.redHot.map((n) => ({ number: n, text: '热' }))} color="hot" />
            </div>
            <div className="stats-line">
              <span className="stat-label">红球冷</span>
              <Chips items={stats.redCold.map((n) => ({ number: n, text: '冷' }))} color="cold" />
            </div>
            <div className="stats-line">
              <span className="stat-label">蓝球热</span>
              <Chips items={stats.blueHot.map((n) => ({ number: n, text: '热' }))} color="hot" />
            </div>
            <div className="stats-line">
              <span className="stat-label">蓝球冷</span>
              <Chips items={stats.blueCold.map((n) => ({ number: n, text: '冷' }))} color="cold" />
            </div>
          </section>

          <section className="stats-section">
            <h3>遗漏</h3>
            <div className="stats-line">
              <span className="stat-label">红球</span>
              <Chips
                items={Array.from({ length: 33 }, (_, i) => i + 1)
                  .sort((a, b) => stats.redOmission[b] - stats.redOmission[a])
                  .map((n) => ({ number: n, text: `${stats.redOmission[n]}期` }))}
                color="plain"
              />
            </div>
            <div className="stats-line">
              <span className="stat-label">蓝球</span>
              <Chips
                items={Array.from({ length: 16 }, (_, i) => i + 1)
                  .sort((a, b) => stats.blueOmission[b] - stats.blueOmission[a])
                  .map((n) => ({ number: n, text: `${stats.blueOmission[n]}期` }))}
                color="plain"
              />
            </div>
          </section>

          <section className="stats-section">
            <h3>比例</h3>
            <div className="stats-line">
              <span className="stat-label">红球奇偶</span>
              <span className="ratio-text">
                奇 {stats.redOddEven.odd} : 偶 {stats.redOddEven.even}
              </span>
            </div>
            <div className="stats-line">
              <span className="stat-label">红球大小</span>
              <span className="ratio-text">
                大 {stats.redBigSmall.big} : 小 {stats.redBigSmall.small}
              </span>
            </div>
            <div className="stats-line">
              <span className="stat-label">蓝球奇偶</span>
              <span className="ratio-text">
                奇 {stats.blueOddEven.odd} : 偶 {stats.blueOddEven.even}
              </span>
            </div>
          </section>
        </div>
      )}
    </div>
  )
}
