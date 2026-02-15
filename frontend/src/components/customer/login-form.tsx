"use client";

import { useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useTranslations, useLocale } from "next-intl";
import { useRouter } from "@/i18n/navigation";
import { Eye, EyeOff, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  CardFooter,
} from "@/components/ui/card";
import {
  Form,
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
} from "@/components/ui/form";
import { Link } from "@/i18n/navigation";
import { loginCustomer } from "@/lib/api/auth";
import { ApiError } from "@/lib/api/client";
import { useAuth } from "@/lib/auth-context";

interface LoginFormValues {
  email: string;
  password: string;
}

export function LoginForm() {
  const t = useTranslations("auth.login");
  const tAuth = useTranslations("auth");
  const locale = useLocale();
  const router = useRouter();
  const { isAuthenticated, login } = useAuth();
  const [showPassword, setShowPassword] = useState(false);
  const [serverError, setServerError] = useState<string | null>(null);

  useEffect(() => {
    if (isAuthenticated) {
      router.replace("/dashboard");
    }
  }, [isAuthenticated, router]);

  const schema = z.object({
    email: z.string().min(1, t("emailRequired")).email(t("emailInvalid")),
    password: z.string().min(1, t("passwordRequired")),
  });

  const form = useForm<LoginFormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  async function onSubmit(values: LoginFormValues) {
    setServerError(null);
    try {
      const tokens = await loginCustomer(values);
      login(tokens);
      router.push("/dashboard");
    } catch (error) {
      if (error instanceof ApiError && error.status === 401) {
        setServerError(t("invalidCredentials"));
      } else {
        setServerError(t("serverError"));
      }
    }
  }

  return (
    <Card className="w-full max-w-md">
      <CardHeader>
        <CardTitle className="font-serif text-2xl">{tAuth("loginHeading")}</CardTitle>
      </CardHeader>
      <CardContent>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="grid gap-4">
            {serverError && (
              <div className="rounded-md bg-destructive/10 px-4 py-3 text-sm text-destructive">
                {serverError}
              </div>
            )}

            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("emailLabel")}</FormLabel>
                  <FormControl>
                    <Input
                      type="email"
                      placeholder={t("emailPlaceholder")}
                      dir="ltr"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <div className="flex items-center justify-between">
                    <FormLabel>{t("passwordLabel")}</FormLabel>
                    <Link
                      href="/forgot-password"
                      className="text-xs text-muted-foreground underline-offset-4 hover:underline"
                      tabIndex={-1}
                    >
                      {t("forgotPassword")}
                    </Link>
                  </div>
                  <FormControl>
                    <div className="relative">
                      <Input
                        type={showPassword ? "text" : "password"}
                        placeholder={t("passwordPlaceholder")}
                        className={locale === "ar" ? "pl-9" : "pr-9"}
                        dir="ltr"
                        {...field}
                      />
                      <Button
                        type="button"
                        variant="ghost"
                        size="icon-xs"
                        className={`absolute top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground ${
                          locale === "ar" ? "left-2" : "right-2"
                        }`}
                        onClick={() => setShowPassword(!showPassword)}
                        tabIndex={-1}
                      >
                        {showPassword ? (
                          <EyeOff className="size-4" />
                        ) : (
                          <Eye className="size-4" />
                        )}
                      </Button>
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <Button
              type="submit"
              className="w-full"
              disabled={form.formState.isSubmitting}
            >
              {form.formState.isSubmitting ? (
                <>
                  <Loader2 className="size-4 animate-spin" />
                  {t("submitting")}
                </>
              ) : (
                t("submit")
              )}
            </Button>
          </form>
        </Form>
      </CardContent>
      <CardFooter className="justify-center">
        <p className="text-sm text-muted-foreground">
          {t("noAccount")}{" "}
          <Link href="/register" className="text-primary underline-offset-4 hover:underline">
            {t("signUp")}
          </Link>
        </p>
      </CardFooter>
    </Card>
  );
}
