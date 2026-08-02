import {
  useReactTable,
  getCoreRowModel,
  getSortedRowModel,
  getExpandedRowModel,
  flexRender,
  type ColumnDef,
  type SortingState,
  type ExpandedState,
  type Row,
} from '@tanstack/react-table'
import { useState, Fragment, type ReactNode } from 'react'

interface DataTableProps<T> {
  columns: ColumnDef<T, unknown>[]
  data: T[]
  caption?: string
  emptyMessage?: string
  zebra?: boolean
  density?: 'compact' | 'normal'
  renderDetail?: (row: T) => ReactNode
}

export function DataTable<T>({
  columns,
  data,
  caption,
  emptyMessage = 'No records found.',
  zebra = true,
  density = 'normal',
  renderDetail,
}: DataTableProps<T>) {
  const [sorting, setSorting] = useState<SortingState>([])
  const [expanded, setExpanded] = useState<ExpandedState>({})

  const table = useReactTable({
    columns,
    data,
    state: { sorting, expanded },
    onSortingChange: setSorting,
    onExpandedChange: setExpanded,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getExpandedRowModel: getExpandedRowModel(),
    getRowCanExpand: () => !!renderDetail,
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
              {renderDetail && (
                <th className="w-8 px-2" scope="col" aria-label="Expand" />
              )}
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
                    {header.isPlaceholder ? null : flexRender(header.column.columnDef.header, header.getContext())}
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
            table.getRowModel().rows.map((row: Row<T>, i) => (
              <Fragment key={row.id}>
                <tr
                  className={`border-b border-[hsl(var(--border-subtle))]/60 transition-colors hover:bg-[hsl(var(--hover-row))] ${zebra && i % 2 === 1 ? 'bg-[hsl(var(--ui-02))]/50' : ''} ${renderDetail ? 'cursor-pointer' : ''}`}
                  onClick={renderDetail ? () => row.toggleExpanded() : undefined}
                >
                  {renderDetail && (
                    <td className="w-8 px-2 text-center text-[hsl(var(--text-03))]" onClick={(e) => { e.stopPropagation(); row.toggleExpanded() }}>
                      <span aria-hidden="true" className="inline-block transition-transform" style={{ transform: row.getIsExpanded() ? 'rotate(90deg)' : 'none' }}>
                        {'\u203A'}
                      </span>
                    </td>
                  )}
                  {row.getVisibleCells().map((cell) => (
                    <td key={cell.id} className={`px-4 ${rowHeight} text-[hsl(var(--text-01))]`}>
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </td>
                  ))}
                </tr>
                {renderDetail && row.getIsExpanded() && (
                  <tr className="bg-[hsl(var(--ui-02))]/30">
                    <td colSpan={row.getVisibleCells().length + 1} className="px-4 py-4">
                      <div className="rounded-[var(--radius-sm)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] p-4">
                        {renderDetail(row.original)}
                      </div>
                    </td>
                  </tr>
                )}
              </Fragment>
            ))
          ) : (
            <tr>
              <td colSpan={columns.length + (renderDetail ? 1 : 0)} className="px-4 py-12 text-center text-sm text-[hsl(var(--text-03))]">
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