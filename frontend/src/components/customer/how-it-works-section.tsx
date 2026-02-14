import { Search, ClipboardList, Truck } from "lucide-react";

const steps = [
  {
    icon: Search,
    step: "1",
    title: "Browse",
    description: "Explore our full range of halal meats, from everyday staples to specialty cuts.",
  },
  {
    icon: ClipboardList,
    step: "2",
    title: "Order",
    description: "Place your order online. Choose your preferred cuts, weights, and quantities.",
  },
  {
    icon: Truck,
    step: "3",
    title: "Pickup or Delivery",
    description: "Collect in-store at your chosen time, or have it delivered fresh to your door.",
  },
];

export function HowItWorksSection() {
  return (
    <section className="py-24">
      <div className="mx-auto max-w-7xl px-6">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="font-serif text-3xl font-bold tracking-tight text-foreground sm:text-4xl">
            How It Works
          </h2>
          <p className="mt-4 text-lg text-muted-foreground">
            Ordering fresh meat has never been easier.
          </p>
        </div>

        <div className="mt-16 grid gap-12 sm:grid-cols-3">
          {steps.map((item) => (
            <div key={item.step} className="flex flex-col items-center text-center">
              <div className="relative flex size-16 items-center justify-center rounded-full bg-primary text-primary-foreground">
                <item.icon className="size-7" />
                <span className="absolute -right-1 -top-1 flex size-6 items-center justify-center rounded-full bg-secondary text-xs font-bold text-secondary-foreground">
                  {item.step}
                </span>
              </div>
              <h3 className="mt-6 font-serif text-xl font-semibold text-foreground">
                {item.title}
              </h3>
              <p className="mt-2 max-w-xs text-sm leading-relaxed text-muted-foreground">
                {item.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
