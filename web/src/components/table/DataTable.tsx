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
  const rowHeight = density === 'compact' ? 'py-2' : 'py-3'
  const headerHeight = density === 'compact' ? 'py-2' : 'py-2.5'

  return (
    <div className="w-full overflow-x-auto rounded-[var(--radius-md)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] shadow-[var(--shadow-01)]">
      <table className="w-full text-sm" role="table">
        {caption && <caption className="sr-only">{caption}</caption>}
        <thead className="sticky top-0 z-10 bg-[hsl(var(--ui-03))]">
          {table.getHeaderGroups().map((headerGroup) => (
            <tr key={headerGroup.id} className="border-b border-[hsl(var(--border-subtle))]">
              {renderDetail && (
                <th className="w-10 px-3" scope="col" aria-label="Expand row" />
              )}
              {headerGroup.headers.map((header) => {
                const sortDir = header.column.getIsSorted()
                const isSortable = header.column.getCanSort()
                return (
                  <th
                    key={header.id}
                    scope="col"
                    className={`px-4 ${headerHeight} text-left font-semibold text-[hsl(var(--text-02))] whitespace-nowrap transition-colors var(--duration-fast-01) var(--ease-productive) ${isSortable ? 'cursor-pointer hover:text-[hsl(var(--text-01))]' : ''}`}
                    aria-sort={sortDir === 'asc' ? 'ascending' : sortDir === 'desc' ? 'descending' : 'none'}
                    onClick={isSortable ? header.column.getToggleSortingHandler() : undefined}
                  >
                    {header.isPlaceholder ? null : flexRender(header.column.columnDef.header, header.getContext())}
                    {isSortable && (
                      <span className="ml-1.5 inline-flex flex-col text-[10px] leading-none text-[hsl(var(--text-03))]" aria-hidden="true">
                        <span className={sortDir === 'asc' ? 'text-[hsl(var(--interactive-04))]' : ''}>{'\u25B2'}</span>
                        <span className={sortDir === 'desc' ? 'text-[hsl(var(--interactive-04))]' : ''}>{'\u25BC'}</span>
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
                  className={`border-b border-[hsl(var(--border-subtle))]/60 transition-colors var(--duration-fast-01) var(--ease-productive) hover:bg-[hsl(var(--hover-row))] ${zebra && i % 2 === 1 ? 'bg-[hsl(var(--ui-02))]/40' : ''} ${renderDetail ? 'cursor-pointer' : ''}`}
                  onClick={renderDetail ? () => row.toggleExpanded() : undefined}
                >
                  {renderDetail && (
                    <td className="w-10 px-3 text-center" onClick={(e) => { e.stopPropagation(); row.toggleExpanded() }}>
                      <span
                        aria-hidden="true"
                        className="inline-flex h-5 w-5 items-center justify-center rounded-[var(--radius-sm)] text-[hsl(var(--text-03))] transition-transform var(--duration-moderate-01) var(--ease-productive)"
                        style={{ transform: row.getIsExpanded() ? 'rotate(90deg)' : 'none' }}
                      >
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
                      <div
                        className="rounded-[var(--radius-md)] border border-[hsl(var(--border-subtle))] bg-[hsl(var(--ui-01))] p-5 shadow-[var(--shadow-02)]"
                        style={{ animation: 'var(--duration-moderate-02) var(--ease-entrance) fadeIn' }}
                      >
                        {renderDetail(row.original)}
                      </div>
                    </td>
                  </tr>
                )}
              </Fragment>
            ))
          ) : (
            <tr>
              <td colSpan={columns.length + (renderDetail ? 1 : 0)} className="px-4 py-16 text-center text-sm text-[hsl(var(--text-03))]">
                {emptyMessage}
              </td>
            </tr>
          )}
        </tbody>
      </table>

      <style>{`
        @keyframes fadeIn {
          from { opacity: 0; transform: translateY(-4px); }
          to { opacity: 1; transform: translateY(0); }
        }
      `}</style>
    </div>
  )
}

export { type ColumnDef }