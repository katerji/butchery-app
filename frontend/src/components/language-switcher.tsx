"use client";

import { useLocale } from "next-intl";
import { useRouter, usePathname } from "@/i18n/navigation";
import { Globe } from "lucide-react";
import { Button } from "@/components/ui/button";

export function LanguageSwitcher() {
  const locale = useLocale();
  const router = useRouter();
  const pathname = usePathname();

  const otherLocale = locale === "en" ? "ar" : "en";
  const label = locale === "en" ? "العربية" : "EN";

  function switchLocale() {
    router.replace(pathname, { locale: otherLocale });
  }

  return (
    <Button
      variant="ghost"
      size="sm"
      onClick={switchLocale}
      className="gap-1.5 cursor-pointer"
    >
      <Globe className="size-4" />
      {label}
    </Button>
  );
}
