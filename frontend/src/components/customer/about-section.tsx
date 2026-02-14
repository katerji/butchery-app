import { useTranslations } from "next-intl";
import { ShieldCheck, Leaf, Award } from "lucide-react";

const valueKeys = [
  { key: "halalCertified", icon: ShieldCheck },
  { key: "responsiblySourced", icon: Leaf },
  { key: "expertButchers", icon: Award },
] as const;

export function AboutSection() {
  const t = useTranslations("about");

  return (
    <section id="about" className="bg-muted/50 py-24">
      <div className="mx-auto max-w-7xl px-6">
        <div className="grid items-center gap-16 lg:grid-cols-2">
          <div>
            <h2 className="font-serif text-3xl font-bold tracking-tight text-foreground sm:text-4xl">
              {t("heading")}
            </h2>
            <p className="mt-6 text-lg leading-relaxed text-muted-foreground">
              {t("paragraph1")}
            </p>
            <p className="mt-4 text-lg leading-relaxed text-muted-foreground">
              {t("paragraph2")}
            </p>
          </div>

          <div className="space-y-6">
            {valueKeys.map(({ key, icon: Icon }) => (
              <div key={key} className="flex gap-4">
                <div className="flex size-12 shrink-0 items-center justify-center rounded-lg bg-primary/10">
                  <Icon className="size-6 text-primary" />
                </div>
                <div>
                  <h3 className="font-semibold text-foreground">
                    {t(`${key}.title`)}
                  </h3>
                  <p className="mt-1 text-sm leading-relaxed text-muted-foreground">
                    {t(`${key}.description`)}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
