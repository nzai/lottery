// 复选框刷选状态机（纯函数，便于单测）
export type BrushMode = 'select' | 'deselect'

export interface BrushState {
  selected: Set<string> // 已选期号集合
  mode: BrushMode // 当前刷选手势的模式
  lastIssue: string | null // 上一次划过的期号
}

// 手势起点：从未选中开始=勾选模式，从已选中开始=取消模式
export function brushStart(selected: Set<string>, issue: string): BrushState {
  const mode: BrushMode = selected.has(issue) ? 'deselect' : 'select'
  return { selected: apply(selected, issue, mode), mode, lastIssue: issue }
}

// 划过某期：按当前模式勾选/取消
export function brushMove(state: BrushState, issue: string): BrushState {
  if (state.lastIssue === issue) return state
  return { ...state, selected: apply(state.selected, issue, state.mode), lastIssue: issue }
}

function apply(selected: Set<string>, issue: string, mode: BrushMode): Set<string> {
  const next = new Set(selected)
  if (mode === 'select') next.add(issue)
  else next.delete(issue)
  return next
}
