import Link from "next/link";
import { Button } from "@/components/ui/button";
import { ArrowRight } from "lucide-react";

export function HeroSection() {
  return (
    <section className="relative overflow-hidden">
      <div className="absolute inset-0 bg-gradient-to-b from-primary/5 to-transparent" />

      <div className="relative mx-auto flex min-h-[calc(100vh-4rem)] max-w-7xl flex-col items-center justify-center px-6 py-24 text-center">
        <p className="text-sm font-medium uppercase tracking-widest text-primary">
          Premium Halal Meats
        </p>

        <h1 className="mt-4 max-w-3xl font-serif text-4xl font-bold tracking-tight text-foreground sm:text-5xl lg:text-6xl">
          Fresh cuts, expertly prepared for your table
        </h1>

        <p className="mt-6 max-w-2xl text-lg leading-relaxed text-muted-foreground">
          Hand-selected halal beef, lamb, poultry, and more — cut to order and
          delivered fresh. Quality you can taste in every bite.
        </p>

        <div className="mt-10 flex flex-col gap-4 sm:flex-row">
          <Button size="lg" asChild>
            <a href="#categories">
              Browse our meats
              <ArrowRight className="size-4" />
            </a>
          </Button>
          <Button variant="outline" size="lg" asChild>
            <Link href="/register">Create an account</Link>
          </Button>
        </div>
      </div>
    </section>
  );
}
