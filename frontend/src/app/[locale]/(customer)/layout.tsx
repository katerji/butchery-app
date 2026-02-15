import { Header } from "@/components/customer/header";
import { Footer } from "@/components/customer/footer";
import { AuthProvider } from "@/lib/auth-context";

export default function CustomerLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <AuthProvider>
      <div className="flex min-h-screen flex-col">
        <Header />
        <main className="flex-1">{children}</main>
        <Footer />
      </div>
    </AuthProvider>
  );
}
