import { useTranslations } from "next-intl";
import { Link } from "@/i18n/navigation";
import { Button } from "@/components/ui/button";
import { ArrowLeft, ArrowRight, CheckCircle2 } from "lucide-react";
import { useLocale } from "next-intl";

interface Props {
  searchParams: Promise<{ registered?: string }>;
}

export default async function CustomerLoginPage({ searchParams }: Props) {
  const params = await searchParams;
  const registered = params.registered === "true";

  return <LoginContent registered={registered} />;
}

function LoginContent({ registered }: { registered: boolean }) {
  const t = useTranslations();
  const locale = useLocale();
  const BackArrow = locale === "ar" ? ArrowRight : ArrowLeft;

  return (
    <div className="flex min-h-[60vh] flex-col items-center justify-center px-6">
      {registered && (
        <div className="mb-6 flex w-full max-w-md items-center gap-3 rounded-md border border-green-200 bg-green-50 px-4 py-3 text-sm text-green-800 dark:border-green-800 dark:bg-green-950/50 dark:text-green-200">
          <CheckCircle2 className="size-4 shrink-0" />
          {t("auth.registrationSuccess")}
        </div>
      )}
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
