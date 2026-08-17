interface ToolbarProps {
  range: number
  ranges: number[]
  fontScale: number
  onRangeChange: (range: number) => void
  onFontChange: (delta: number) => void
}

export function Toolbar({ range, ranges, fontScale, onRangeChange, onFontChange }: ToolbarProps) {
  return (
    <header className="toolbar">
      <h1 className="toolbar-title">双色球走势图</h1>
      <div className="toolbar-controls">
        <div className="range-group">
          {ranges.map((r) => (
            <button
              key={r}
              className={`range-btn${r === range ? ' active' : ''}`}
              onClick={() => onRangeChange(r)}
            >
              {r}
            </button>
          ))}
        </div>
        <div className="font-group">
          <button className="font-btn" onClick={() => onFontChange(-0.1)} aria-label="缩小">
            A−
          </button>
          <span className="font-scale">{Math.round(fontScale * 100)}%</span>
          <button className="font-btn" onClick={() => onFontChange(0.1)} aria-label="放大">
            A+
          </button>
        </div>
      </div>
    </header>
  )
}
