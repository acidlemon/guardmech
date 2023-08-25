export type BSelectItem = {
  label: string
  value: string
  tips?: string
  disabled?: boolean
}

export type BTableColumn = {
  key: string
  label: string
}

export type BTableRow = {
  [index: string]: any
}