"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { CheckCircle2, X } from "lucide-react";

interface SuccessBannerProps {
  messageKey: string;
}

export function SuccessBanner({ messageKey }: SuccessBannerProps) {
  const [dismissed, setDismissed] = useState(false);
  const t = useTranslations();

  if (dismissed) {
    return null;
  }

  return (
    <div className="flex items-center gap-3 rounded-md border border-green-200 bg-green-50 px-4 py-3 text-sm text-green-800 dark:border-green-800 dark:bg-green-950/50 dark:text-green-200">
      <CheckCircle2 className="size-4 shrink-0" />
      <span className="flex-1">{t(messageKey)}</span>
      <button
        type="button"
        onClick={() => setDismissed(true)}
        className="shrink-0 cursor-pointer text-green-800/60 transition-colors hover:text-green-800 dark:text-green-200/60 dark:hover:text-green-200"
        aria-label="Dismiss"
      >
        <X className="size-4" />
      </button>
    </div>
  );
}
