import { useCallback, useRef } from 'react'
import { brushMove, brushStart, type BrushState } from '../lib/brush'

// 将复选框列的 pointer 事件转换为刷选状态更新。
// 选择集合跨手势保留在 ref 中，onChange 在每次变化时通知。
export function useBrush(onChange: (selected: Set<string>) => void) {
  const selectionRef = useRef<Set<string>>(new Set())
  const activeRef = useRef<BrushState | null>(null)
  const onChangeRef = useRef(onChange)
  onChangeRef.current = onChange

  const emit = useCallback(() => {
    onChangeRef.current(new Set(selectionRef.current))
  }, [])

  const handlePointerDown = useCallback(
    (e: React.PointerEvent, issue: string) => {
      if (e.pointerType === 'mouse' && e.button !== 0) return
      e.preventDefault()
      activeRef.current = brushStart(selectionRef.current, issue)
      selectionRef.current = activeRef.current.selected
      emit()
      ;(e.currentTarget as HTMLElement).setPointerCapture(e.pointerId)
    },
    [emit],
  )

  const handlePointerMove = useCallback(
    (e: React.PointerEvent) => {
      const active = activeRef.current
      if (!active) return
      // 手指/光标当前所在的行（几何命中，跨单元格）
      const el = document.elementFromPoint(e.clientX, e.clientY)?.closest('[data-issue]')
      const issue = el?.getAttribute('data-issue')
      if (issue) {
        activeRef.current = brushMove(active, issue)
        selectionRef.current = activeRef.current.selected
        emit()
      }
    },
    [emit],
  )

  const handlePointerUp = useCallback((e: React.PointerEvent) => {
    activeRef.current = null
    ;(e.currentTarget as HTMLElement).releasePointerCapture?.(e.pointerId)
  }, [])

  const clear = useCallback(() => {
    selectionRef.current = new Set()
    emit()
  }, [emit])

  return { handlePointerDown, handlePointerMove, handlePointerUp, clear }
}
