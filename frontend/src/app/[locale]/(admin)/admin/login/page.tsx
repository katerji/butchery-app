import { useTranslations } from "next-intl";
import { Link } from "@/i18n/navigation";
import { Button } from "@/components/ui/button";
import { ArrowLeft, ArrowRight } from "lucide-react";
import { useLocale } from "next-intl";

export default function AdminLoginPage() {
  const t = useTranslations();
  const locale = useLocale();
  const BackArrow = locale === "ar" ? ArrowRight : ArrowLeft;

  return (
    <div className="w-full max-w-sm px-6 text-center">
      <h1 className="font-serif text-2xl font-bold text-foreground">
        {t("auth.adminLoginHeading")}
      </h1>
      <p className="mt-4 text-sm text-muted-foreground">
        {t("common.comingSoon")}
      </p>
      <Button variant="outline" size="sm" className="mt-8" asChild>
        <Link href="/">
          <BackArrow className="size-4" />
          {t("common.backToHome")}
        </Link>
      </Button>
    </div>
  );
}
