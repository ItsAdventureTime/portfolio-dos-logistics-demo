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
}

export function DataTable<T>({ columns, data, caption, emptyMessage = 'No records found.' }: DataTableProps<T>) {
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

  return (
    <div className="w-full overflow-x-auto rounded-lg border border-[hsl(var(--border))]">
      <table className="w-full text-sm" role="table">
        {caption && <caption className="sr-only">{caption}</caption>}
        <thead className="bg-[hsl(var(--surface-panel))]">
          {table.getHeaderGroups().map((headerGroup) => (
            <tr key={headerGroup.id} className="border-b border-[hsl(var(--border))]">
              {headerGroup.headers.map((header) => {
                const sortDir = header.column.getIsSorted()
                const isSortable = header.column.getCanSort()
                return (
                  <th
                    key={header.id}
                    scope="col"
                    className="px-4 py-3 text-left font-medium text-[hsl(var(--content-muted))]"
                    aria-sort={sortDir === 'asc' ? 'ascending' : sortDir === 'desc' ? 'descending' : 'none'}
                  >
                    {header.isPlaceholder ? null : flexRender(
                      header.column.columnDef.header,
                      header.getContext(),
                    )}
                    {isSortable && (
                      <span className="ml-1 text-xs" aria-hidden="true">
                        {sortDir === 'asc' ? '▲' : sortDir === 'desc' ? '▼' : '↕'}
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
            table.getRowModel().rows.map((row) => (
              <tr
                key={row.id}
                className="border-b border-[hsl(var(--border))]/50 hover:bg-[hsl(var(--surface-panel))]/50 transition-colors"
              >
                {row.getVisibleCells().map((cell) => (
                  <td key={cell.id} className="px-4 py-3 text-[hsl(var(--content))]">
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </td>
                ))}
              </tr>
            ))
          ) : (
            <tr>
              <td colSpan={columns.length} className="px-4 py-8 text-center text-[hsl(var(--content-muted))]">
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