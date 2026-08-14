export type ToolType =
  | 'cursor'
  | 'pen'
  | 'highlight'
  | 'eraser'
  | 'underline'
  | 'strike'
  | 'note'

export interface ToolItem {
  id: ToolType
  label: string
  icon: string
}

export interface SidebarItem {
  id: string
  label: string
  icon: string
  badge?: string
}
