import { useTranslations } from "next-intl";
import { Search, ClipboardList, Truck } from "lucide-react";

const steps = [
  { key: "step1", icon: Search },
  { key: "step2", icon: ClipboardList },
  { key: "step3", icon: Truck },
] as const;

export function HowItWorksSection() {
  const t = useTranslations("howItWorks");

  return (
    <section className="py-24">
      <div className="mx-auto max-w-7xl px-6">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="font-serif text-3xl font-bold tracking-tight text-foreground sm:text-4xl">
            {t("heading")}
          </h2>
          <p className="mt-4 text-lg text-muted-foreground">
            {t("description")}
          </p>
        </div>

        <div className="mt-16 grid gap-12 sm:grid-cols-3">
          {steps.map(({ key, icon: Icon }, index) => (
            <div key={key} className="flex flex-col items-center text-center">
              <div className="relative flex size-16 items-center justify-center rounded-full bg-primary text-primary-foreground">
                <Icon className="size-7" />
                <span className="absolute -right-1 -top-1 flex size-6 items-center justify-center rounded-full bg-secondary text-xs font-bold text-secondary-foreground">
                  {index + 1}
                </span>
              </div>
              <h3 className="mt-6 font-serif text-xl font-semibold text-foreground">
                {t(`${key}.title`)}
              </h3>
              <p className="mt-2 max-w-xs text-sm leading-relaxed text-muted-foreground">
                {t(`${key}.description`)}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
