import { useTranslations, useLocale } from "next-intl";
import { Link } from "@/i18n/navigation";
import { Button } from "@/components/ui/button";
import { ArrowRight, ArrowLeft } from "lucide-react";

export function HeroSection() {
  const t = useTranslations("hero");
  const locale = useLocale();
  const Arrow = locale === "ar" ? ArrowLeft : ArrowRight;

  return (
    <section className="relative overflow-hidden">
      <div className="absolute inset-0 bg-gradient-to-b from-primary/5 to-transparent" />

      <div className="relative mx-auto flex min-h-[calc(100vh-4rem)] max-w-7xl flex-col items-center justify-center px-6 py-24 text-center">
        <p className="text-sm font-medium uppercase tracking-widest text-primary">
          {t("tagline")}
        </p>

        <h1 className="mt-4 max-w-3xl font-serif text-4xl font-bold tracking-tight text-foreground sm:text-5xl lg:text-6xl">
          {t("heading")}
        </h1>

        <p className="mt-6 max-w-2xl text-lg leading-relaxed text-muted-foreground">
          {t("description")}
        </p>

        <div className="mt-10 flex flex-col gap-4 sm:flex-row">
          <Button size="lg" asChild>
            <a href="#categories">
              {t("browseMeats")}
              <Arrow className="size-4" />
            </a>
          </Button>
          <Button variant="outline" size="lg" asChild>
            <Link href="/register">{t("createAccount")}</Link>
          </Button>
        </div>
      </div>
    </section>
  );
}
