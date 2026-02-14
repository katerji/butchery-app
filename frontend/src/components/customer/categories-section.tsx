import { useTranslations } from "next-intl";
import { Beef, Bird, Drumstick, Link2, Sparkles } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";

const categoryKeys = [
  { key: "beef", icon: Beef },
  { key: "lamb", icon: Drumstick },
  { key: "poultry", icon: Bird },
  { key: "sausages", icon: Link2 },
  { key: "specialty", icon: Sparkles },
] as const;

export function CategoriesSection() {
  const t = useTranslations("categories");

  return (
    <section id="categories" className="py-24">
      <div className="mx-auto max-w-7xl px-6">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="font-serif text-3xl font-bold tracking-tight text-foreground sm:text-4xl">
            {t("heading")}
          </h2>
          <p className="mt-4 text-lg text-muted-foreground">
            {t("description")}
          </p>
        </div>

        <div className="mt-16 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {categoryKeys.map(({ key, icon: Icon }) => (
            <Card
              key={key}
              className="border-border/50 transition-shadow hover:shadow-md"
            >
              <CardContent className="flex flex-col items-center p-8 text-center">
                <div className="flex size-14 items-center justify-center rounded-full bg-primary/10">
                  <Icon className="size-7 text-primary" />
                </div>
                <h3 className="mt-4 font-serif text-xl font-semibold text-foreground">
                  {t(`${key}.name`)}
                </h3>
                <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
                  {t(`${key}.description`)}
                </p>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </section>
  );
}
