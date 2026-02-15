import { LoginForm } from "@/components/customer/login-form";
import { SuccessBanner } from "@/components/customer/success-banner";

interface Props {
  searchParams: Promise<{ registered?: string }>;
}

export default async function CustomerLoginPage({ searchParams }: Props) {
  const params = await searchParams;
  const registered = params.registered === "true";

  return (
    <div className="flex min-h-[60vh] flex-col items-center justify-center px-6">
      {registered && (
        <div className="mb-6 w-full max-w-md">
          <SuccessBanner messageKey="auth.registrationSuccess" />
        </div>
      )}
      <LoginForm />
    </div>
  );
}
