import { Beef, Bird, Drumstick, Link2, Sparkles } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";

const categories = [
  {
    icon: Beef,
    name: "Beef",
    description: "Steaks, roasts, mince, and premium cuts from grass-fed cattle.",
  },
  {
    icon: Drumstick,
    name: "Lamb",
    description: "Chops, shanks, racks, and whole legs — tender and full of flavour.",
  },
  {
    icon: Bird,
    name: "Poultry",
    description: "Whole chickens, breasts, thighs, and wings — fresh and free-range.",
  },
  {
    icon: Link2,
    name: "Sausages",
    description: "House-made sausages in classic and seasonal varieties.",
  },
  {
    icon: Sparkles,
    name: "Specialty Cuts",
    description: "Marinated meats, kebab packs, and custom cuts prepared to order.",
  },
];

export function CategoriesSection() {
  return (
    <section id="categories" className="py-24">
      <div className="mx-auto max-w-7xl px-6">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="font-serif text-3xl font-bold tracking-tight text-foreground sm:text-4xl">
            Our Meats
          </h2>
          <p className="mt-4 text-lg text-muted-foreground">
            All our products are 100% halal certified, sourced from trusted
            farms, and prepared fresh daily.
          </p>
        </div>

        <div className="mt-16 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {categories.map((category) => (
            <Card
              key={category.name}
              className="group cursor-pointer border-border/50 transition-all hover:border-primary/30 hover:shadow-md"
            >
              <CardContent className="flex flex-col items-center p-8 text-center">
                <div className="flex size-14 items-center justify-center rounded-full bg-primary/10 transition-colors group-hover:bg-primary/20">
                  <category.icon className="size-7 text-primary" />
                </div>
                <h3 className="mt-4 font-serif text-xl font-semibold text-foreground">
                  {category.name}
                </h3>
                <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
                  {category.description}
                </p>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </section>
  );
}
