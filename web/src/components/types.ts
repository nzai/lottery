// 刷选 pointer 事件处理器（由 useBrush 提供）
export interface BrushHandlers {
  handlePointerDown: (e: React.PointerEvent, issue: string) => void
  handlePointerMove: (e: React.PointerEvent) => void
  handlePointerUp: (e: React.PointerEvent) => void
  clear: () => void
}
