import {
  Modal,
  ModalOverlay,
  Dialog as RACDialog,
  Button as RACButton,
  Heading,
} from 'react-aria-components'
import { type ReactNode } from 'react'
import { Button } from './Button'

interface DialogProps {
  title: string
  children: ReactNode
  isOpen: boolean
  onOpenChange: (open: boolean) => void
  primaryActionLabel?: string
  onPrimaryAction?: () => void
  isPrimaryDisabled?: boolean
  cancelLabel?: string
}

export function Dialog({
  title,
  children,
  isOpen,
  onOpenChange,
  primaryActionLabel,
  onPrimaryAction,
  isPrimaryDisabled,
  cancelLabel = 'Cancel',
}: DialogProps) {
  return (
    <ModalOverlay
      isOpen={isOpen}
      onOpenChange={onOpenChange}
      isDismissable
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
    >
      <Modal className="w-full max-w-lg rounded-xl bg-[hsl(var(--surface-elevated))] border border-[hsl(var(--border))] shadow-xl">
        <RACDialog className="flex flex-col gap-4 p-6 outline-none">
          <Heading slot="title" className="text-lg font-semibold text-[hsl(var(--content))]">
            {title}
          </Heading>
          <div className="text-sm text-[hsl(var(--content-muted))]">{children}</div>
          {(primaryActionLabel || cancelLabel) && (
            <div className="flex justify-end gap-3 pt-2">
              <RACButton
                onPress={() => onOpenChange(false)}
                className="rounded-lg px-4 py-2 text-sm text-[hsl(var(--content-muted))] hover:text-[hsl(var(--content))] focus:outline-none focus-visible:ring-2 focus-visible:ring-[hsl(var(--focus))]"
              >
                {cancelLabel}
              </RACButton>
              {primaryActionLabel && (
                <Button onPress={onPrimaryAction} isDisabled={isPrimaryDisabled} size="sm">
                  {primaryActionLabel}
                </Button>
              )}
            </div>
          )}
        </RACDialog>
      </Modal>
    </ModalOverlay>
  )
}

interface AlertDialogProps {
  title: string
  children: ReactNode
  isOpen: boolean
  onOpenChange: (open: boolean) => void
  confirmLabel: string
  onConfirm: () => void
  cancelLabel?: string
}

export function AlertDialog({
  title,
  children,
  isOpen,
  onOpenChange,
  confirmLabel,
  onConfirm,
  cancelLabel = 'Cancel',
}: AlertDialogProps) {
  return (
    <ModalOverlay
      isOpen={isOpen}
      onOpenChange={onOpenChange}
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
    >
      <Modal className="w-full max-w-md rounded-xl bg-[hsl(var(--surface-elevated))] border border-[hsl(var(--border))] shadow-xl">
        <RACDialog role="alertdialog" className="flex flex-col gap-4 p-6 outline-none">
          <Heading slot="title" className="text-lg font-semibold text-[hsl(var(--content))]">
            {title}
          </Heading>
          <div className="text-sm text-[hsl(var(--content-muted))]">{children}</div>
          <div className="flex justify-end gap-3 pt-2">
            <Button variant="ghost" size="sm" onPress={() => onOpenChange(false)}>
              {cancelLabel}
            </Button>
            <Button variant="destructive" size="sm" onPress={onConfirm}>
              {confirmLabel}
            </Button>
          </div>
        </RACDialog>
      </Modal>
    </ModalOverlay>
  )
}