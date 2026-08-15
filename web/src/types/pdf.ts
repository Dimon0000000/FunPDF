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

export interface PdfPoint {
  x: number
  y: number
}

interface AnnotationBase {
  id: string
  page: number
  color: string
  width: number
}

export interface InkAnnotation extends AnnotationBase {
  type: 'pen'
  points: PdfPoint[]
}

export interface LineAnnotation extends AnnotationBase {
  type: 'underline' | 'strike'
  start: PdfPoint
  end: PdfPoint
}

export interface HighlightAnnotation extends AnnotationBase {
  type: 'highlight'
  start: PdfPoint
  end: PdfPoint
}

export interface NoteAnnotation extends AnnotationBase {
  type: 'note'
  point: PdfPoint
  text: string
}

export type PdfAnnotation =
  | InkAnnotation
  | LineAnnotation
  | HighlightAnnotation
  | NoteAnnotation
