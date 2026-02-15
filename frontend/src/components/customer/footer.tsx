"use client";

import { useTranslations } from "next-intl";
import { Link } from "@/i18n/navigation";
import { Separator } from "@/components/ui/separator";
import { useAuth } from "@/lib/auth-context";

export function Footer() {
  const t = useTranslations("footer");
  const tNav = useTranslations("nav");
  const { isAuthenticated } = useAuth();

  const shopLinks = [
    { href: "#categories", label: t("shopOurMeats") },
    { href: "#specials", label: t("shopWeeklySpecials") },
  ];

  const infoLinks = [
    { href: "#about", label: t("infoAboutUs") },
    { href: "#contact", label: t("infoContact") },
  ];

  const accountLinks = isAuthenticated
    ? [
        { href: "/dashboard" as const, label: tNav("dashboard") },
      ]
    : [
        { href: "/login" as const, label: tNav("logIn") },
        { href: "/register" as const, label: tNav("signUp") },
      ];

  return (
    <footer className="border-t border-border/50 bg-muted/30">
      <div className="mx-auto max-w-7xl px-6 py-12">
        <div className="grid gap-8 sm:grid-cols-2 lg:grid-cols-4">
          <div>
            <Link href="/" className="font-serif text-lg font-bold text-primary">
              Butchery
            </Link>
            <p className="mt-3 max-w-xs text-sm leading-relaxed text-muted-foreground">
              {t("description")}
            </p>
          </div>

          <div>
            <h4 className="text-sm font-semibold text-foreground">{t("shop")}</h4>
            <ul className="mt-3 space-y-2">
              {shopLinks.map((link) => (
                <li key={link.href}>
                  <a
                    href={link.href}
                    className="text-sm text-muted-foreground transition-colors hover:text-foreground"
                  >
                    {link.label}
                  </a>
                </li>
              ))}
            </ul>
          </div>

          <div>
            <h4 className="text-sm font-semibold text-foreground">{t("info")}</h4>
            <ul className="mt-3 space-y-2">
              {infoLinks.map((link) => (
                <li key={link.href}>
                  <a
                    href={link.href}
                    className="text-sm text-muted-foreground transition-colors hover:text-foreground"
                  >
                    {link.label}
                  </a>
                </li>
              ))}
            </ul>
          </div>

          <div>
            <h4 className="text-sm font-semibold text-foreground">{t("account")}</h4>
            <ul className="mt-3 space-y-2">
              {accountLinks.map((link) => (
                <li key={link.href}>
                  <Link
                    href={link.href}
                    className="text-sm text-muted-foreground transition-colors hover:text-foreground"
                  >
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        </div>

        <Separator className="my-8" />

        <p className="text-center text-sm text-muted-foreground">
          {t("copyright", { year: new Date().getFullYear() })}
        </p>
      </div>
    </footer>
  );
}
