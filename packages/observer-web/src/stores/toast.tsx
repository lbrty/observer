import { createContext, useCallback, useContext, useState } from "react";
import type { ReactNode } from "react";

type ToastVariant = "success" | "error" | "info";

interface ToastItem {
  id: string;
  message: string;
  variant: ToastVariant;
}

interface ToastActions {
  success: (message: string) => void;
  error: (message: string) => void;
  info: (message: string) => void;
  dismiss: (id: string) => void;
}

interface ToastContextValue {
  toasts: ToastItem[];
  toast: ToastActions;
}

const ToastContext = createContext<ToastContextValue | null>(null);

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<ToastItem[]>([]);

  const addToast = useCallback((message: string, variant: ToastVariant) => {
    const id = crypto.randomUUID();
    setToasts((prev) => [...prev, { id, message, variant }]);
    if (variant !== "error") {
      setTimeout(() => setToasts((prev) => prev.filter((t) => t.id !== id)), 4000);
    }
  }, []);

  const dismiss = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const toast: ToastActions = {
    success: (msg) => addToast(msg, "success"),
    error: (msg) => addToast(msg, "error"),
    info: (msg) => addToast(msg, "info"),
    dismiss,
  };

  return <ToastContext value={{ toasts, toast }}>{children}</ToastContext>;
}

export function useToast(): ToastActions {
  const ctx = useContext(ToastContext);
  if (!ctx) throw new Error("useToast must be used within ToastProvider");
  return ctx.toast;
}

export function useToasts(): ToastItem[] {
  const ctx = useContext(ToastContext);
  if (!ctx) throw new Error("useToasts must be used within ToastProvider");
  return ctx.toasts;
}
