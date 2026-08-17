import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { Toolbar } from './Toolbar'

const RANGES = [20, 50, 100, 200, 500]

describe('Toolbar', () => {
  it('渲染标题、期数切换按钮与缩放按钮', () => {
    render(<Toolbar range={100} ranges={RANGES} fontScale={1} onRangeChange={() => {}} onFontChange={() => {}} />)
    expect(screen.getByText('双色球走势图')).toBeTruthy()
    expect(screen.getByText('100')).toBeTruthy()
    expect(screen.getByText('A−')).toBeTruthy()
    expect(screen.getByText('A+')).toBeTruthy()
  })

  it('点击期数按钮触发 onRangeChange', () => {
    const onRangeChange = vi.fn()
    render(<Toolbar range={100} ranges={RANGES} fontScale={1} onRangeChange={onRangeChange} onFontChange={() => {}} />)
    fireEvent.click(screen.getByText('50'))
    expect(onRangeChange).toHaveBeenCalledWith(50)
  })

  it('点击缩放按钮触发 onFontChange', () => {
    const onFontChange = vi.fn()
    render(<Toolbar range={100} ranges={RANGES} fontScale={1} onRangeChange={() => {}} onFontChange={onFontChange} />)
    fireEvent.click(screen.getByText('A+'))
    expect(onFontChange).toHaveBeenCalledWith(0.1)
    fireEvent.click(screen.getByText('A−'))
    expect(onFontChange).toHaveBeenCalledWith(-0.1)
  })
})
