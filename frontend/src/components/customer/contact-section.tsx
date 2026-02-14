import { useTranslations } from "next-intl";
import { MapPin, Clock, Phone } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";

const contactKeys = [
  { key: "location", icon: MapPin, lineCount: 2 },
  { key: "hours", icon: Clock, lineCount: 3 },
  { key: "phone", icon: Phone, lineCount: 2 },
] as const;

export function ContactSection() {
  const t = useTranslations("contact");

  return (
    <section id="contact" className="py-24">
      <div className="mx-auto max-w-7xl px-6">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="font-serif text-3xl font-bold tracking-tight text-foreground sm:text-4xl">
            {t("heading")}
          </h2>
          <p className="mt-4 text-lg text-muted-foreground">
            {t("description")}
          </p>
        </div>

        <div className="mt-16 grid gap-6 sm:grid-cols-3">
          {contactKeys.map(({ key, icon: Icon, lineCount }) => (
            <Card key={key} className="border-border/50">
              <CardContent className="flex flex-col items-center p-8 text-center">
                <div className="flex size-12 items-center justify-center rounded-full bg-primary/10">
                  <Icon className="size-6 text-primary" />
                </div>
                <h3 className="mt-4 font-semibold text-foreground">
                  {t(`${key}.title`)}
                </h3>
                <div className="mt-2 space-y-1">
                  {Array.from({ length: lineCount }, (_, i) => (
                    <p key={i} className="text-sm text-muted-foreground">
                      {t(`${key}.line${i + 1}`)}
                    </p>
                  ))}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </section>
  );
}
