import Link from "next/link";
import { Separator } from "@/components/ui/separator";

const shopLinks = [
  { href: "#categories", label: "Our Meats" },
  { href: "#specials", label: "Weekly Specials" },
];

const infoLinks = [
  { href: "#about", label: "About Us" },
  { href: "#contact", label: "Contact" },
];

const accountLinks = [
  { href: "/login", label: "Log in" },
  { href: "/register", label: "Sign up" },
];

export function Footer() {
  return (
    <footer className="border-t border-border/50 bg-muted/30">
      <div className="mx-auto max-w-7xl px-6 py-12">
        <div className="grid gap-8 sm:grid-cols-2 lg:grid-cols-4">
          <div>
            <Link href="/" className="font-serif text-lg font-bold text-primary">
              Butchery
            </Link>
            <p className="mt-3 max-w-xs text-sm leading-relaxed text-muted-foreground">
              Premium halal meats, expertly prepared and delivered fresh to your
              table.
            </p>
          </div>

          <div>
            <h4 className="text-sm font-semibold text-foreground">Shop</h4>
            <ul className="mt-3 space-y-2">
              {shopLinks.map((link) => (
                <li key={link.href}>
                  <a
                    href={link.href}
                    className="text-sm text-muted-foreground transition-colors hover:text-foreground"
                  >
                    {link.label}
                  </a>
                </li>
              ))}
            </ul>
          </div>

          <div>
            <h4 className="text-sm font-semibold text-foreground">Info</h4>
            <ul className="mt-3 space-y-2">
              {infoLinks.map((link) => (
                <li key={link.href}>
                  <a
                    href={link.href}
                    className="text-sm text-muted-foreground transition-colors hover:text-foreground"
                  >
                    {link.label}
                  </a>
                </li>
              ))}
            </ul>
          </div>

          <div>
            <h4 className="text-sm font-semibold text-foreground">Account</h4>
            <ul className="mt-3 space-y-2">
              {accountLinks.map((link) => (
                <li key={link.href}>
                  <Link
                    href={link.href}
                    className="text-sm text-muted-foreground transition-colors hover:text-foreground"
                  >
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        </div>

        <Separator className="my-8" />

        <p className="text-center text-sm text-muted-foreground">
          &copy; {new Date().getFullYear()} Butchery. All rights reserved.
        </p>
      </div>
    </footer>
  );
}
