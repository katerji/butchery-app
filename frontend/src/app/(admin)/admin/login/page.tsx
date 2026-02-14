import type { Metadata } from "next";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { ArrowLeft } from "lucide-react";

export const metadata: Metadata = {
  title: "Admin Login - Butchery",
};

export default function AdminLoginPage() {
  return (
    <div className="w-full max-w-sm px-6 text-center">
      <h1 className="font-serif text-2xl font-bold text-foreground">
        Admin Login
      </h1>
      <p className="mt-4 text-sm text-muted-foreground">
        This page is coming soon. Check back later.
      </p>
      <Button variant="outline" size="sm" className="mt-8" asChild>
        <Link href="/">
          <ArrowLeft className="size-4" />
          Back to home
        </Link>
      </Button>
    </div>
  );
}
