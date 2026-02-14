import Link from "next/link";
import { Button } from "@/components/ui/button";
import { ArrowLeft } from "lucide-react";

export default function CustomerRegisterPage() {
  return (
    <div className="flex min-h-[60vh] flex-col items-center justify-center px-6">
      <div className="text-center">
        <h1 className="font-serif text-3xl font-bold text-foreground">
          Create an account
        </h1>
        <p className="mt-4 text-muted-foreground">
          This page is coming soon. Check back later.
        </p>
        <Button variant="outline" size="sm" className="mt-8" asChild>
          <Link href="/">
            <ArrowLeft className="size-4" />
            Back to home
          </Link>
        </Button>
      </div>
    </div>
  );
}
