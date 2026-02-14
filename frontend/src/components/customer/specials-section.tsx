import { useTranslations } from "next-intl";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";

const specials = [
  {
    key: "lambShoulder",
    weight: "1.5 kg",
    price: "AED 69.99",
    originalPrice: "AED 89.99",
  },
  {
    key: "chickenBreast",
    weight: "1 kg",
    price: "AED 44.99",
    originalPrice: "AED 54.99",
  },
  {
    key: "beefKebab",
    weight: "800 g",
    price: "AED 49.99",
    originalPrice: "AED 62.99",
  },
] as const;

export function SpecialsSection() {
  const t = useTranslations("specials");

  return (
    <section id="specials" className="bg-muted/50 py-24">
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
          {specials.map((item) => (
            <Card key={item.key} className="overflow-hidden border-border/50">
              <div className="flex h-48 items-center justify-center bg-muted">
                <span className="text-sm text-muted-foreground">
                  {t("productImage")}
                </span>
              </div>

              <CardContent className="p-6">
                <div className="flex items-start justify-between">
                  <div>
                    <h3 className="font-semibold text-foreground">
                      {t(`${item.key}.name`)}
                    </h3>
                    <p className="mt-1 text-sm text-muted-foreground">
                      {item.weight}
                    </p>
                  </div>
                  <Badge variant="secondary">{t(`${item.key}.tag`)}</Badge>
                </div>

                <div className="mt-4 flex items-center gap-2">
                  <span className="text-lg font-bold text-primary">
                    {item.price}
                  </span>
                  <span className="text-sm text-muted-foreground line-through">
                    {item.originalPrice}
                  </span>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </section>
  );
}
