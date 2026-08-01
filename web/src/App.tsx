import { Button } from 'react-aria-components'

export default function App() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center gap-4 p-8">
      <h1 className="text-2xl font-semibold text-[hsl(var(--content))]">
        DOS FreightFlow Control
      </h1>
      <p className="text-[hsl(var(--content-muted))]">
        Build verified. Ready for development.
      </p>
      <Button className="rounded-lg bg-[hsl(var(--action))] px-4 py-2 text-[hsl(var(--action-foreground))]">
        Primary action
      </Button>
    </main>
  )
}