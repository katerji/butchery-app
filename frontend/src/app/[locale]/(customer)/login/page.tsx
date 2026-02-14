import { useTranslations } from "next-intl";
import { Link } from "@/i18n/navigation";
import { Button } from "@/components/ui/button";
import { ArrowLeft, ArrowRight } from "lucide-react";
import { useLocale } from "next-intl";

export default function CustomerLoginPage() {
  const t = useTranslations();
  const locale = useLocale();
  const BackArrow = locale === "ar" ? ArrowRight : ArrowLeft;

  return (
    <div className="flex min-h-[60vh] flex-col items-center justify-center px-6">
      <div className="text-center">
        <h1 className="font-serif text-3xl font-bold text-foreground">
          {t("auth.loginHeading")}
        </h1>
        <p className="mt-4 text-muted-foreground">
          {t("common.comingSoon")}
        </p>
        <Button variant="outline" size="sm" className="mt-8" asChild>
          <Link href="/">
            <BackArrow className="size-4" />
            {t("common.backToHome")}
          </Link>
        </Button>
      </div>
    </div>
  );
}
