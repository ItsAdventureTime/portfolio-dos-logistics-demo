import {
  useReactTable,
  getCoreRowModel,
  getSortedRowModel,
  flexRender,
  type ColumnDef,
  type SortingState,
} from '@tanstack/react-table'
import { useState } from 'react'

interface DataTableProps<T> {
  columns: ColumnDef<T, unknown>[]
  data: T[]
  caption?: string
  emptyMessage?: string
  zebra?: boolean
  density?: 'compact' | 'normal'
}

export function DataTable<T>({
  columns,
  data,
  caption,
  emptyMessage = 'No records found.',
  zebra = true,
  density = 'normal',
}: DataTableProps<T>) {
  const [sorting, setSorting] = useState<SortingState>([])

  const table = useReactTable({
    columns,
    data,
    state: { sorting },
    onSortingChange: setSorting,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
  })

  const hasData = data.length > 0
  const rowHeight = density === 'compact' ? 'py-2' : 'py-3.5'
  const headerHeight = density === 'compact' ? 'py-2' : 'py-3'

  return (
    <div className="w-full overflow-x-auto rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))]">
      <table className="w-full text-sm" role="table">
        {caption && <caption className="sr-only">{caption}</caption>}
        <thead className="sticky top-0 z-10 bg-[hsl(var(--ui-03))]">
          {table.getHeaderGroups().map((headerGroup) => (
            <tr key={headerGroup.id} className="border-b border-[hsl(var(--border-subtle))]">
              {headerGroup.headers.map((header) => {
                const sortDir = header.column.getIsSorted()
                const isSortable = header.column.getCanSort()
                return (
                  <th
                    key={header.id}
                    scope="col"
                    className={`px-4 ${headerHeight} text-left font-semibold text-[hsl(var(--text-02))] whitespace-nowrap`}
                    aria-sort={sortDir === 'asc' ? 'ascending' : sortDir === 'desc' ? 'descending' : 'none'}
                  >
                    {header.isPlaceholder ? null : flexRender(
                      header.column.columnDef.header,
                      header.getContext(),
                    )}
                    {isSortable && (
                      <span className="ml-1 inline-block text-xs text-[hsl(var(--text-03))]" aria-hidden="true">
                        {sortDir === 'asc' ? '\u2191' : sortDir === 'desc' ? '\u2193' : '\u2195'}
                      </span>
                    )}
                  </th>
                )
              })}
            </tr>
          ))}
        </thead>
        <tbody>
          {hasData ? (
            table.getRowModel().rows.map((row, i) => (
              <tr
                key={row.id}
                className={`border-b border-[hsl(var(--border-subtle))]/60 transition-colors hover:bg-[hsl(var(--hover-row))] ${zebra && i % 2 === 1 ? 'bg-[hsl(var(--ui-02))]/50' : ''}`}
              >
                {row.getVisibleCells().map((cell) => (
                  <td key={cell.id} className={`px-4 ${rowHeight} text-[hsl(var(--text-01))]`}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </td>
                ))}
              </tr>
            ))
          ) : (
            <tr>
              <td colSpan={columns.length} className="px-4 py-12 text-center text-sm text-[hsl(var(--text-03))]">
                {emptyMessage}
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  )
}

export { type ColumnDef }