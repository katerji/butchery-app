import { ShieldCheck, Leaf, Award } from "lucide-react";

const values = [
  {
    icon: ShieldCheck,
    title: "100% Halal Certified",
    description: "Every product meets strict halal standards, certified and traceable.",
  },
  {
    icon: Leaf,
    title: "Responsibly Sourced",
    description: "We partner with trusted farms that prioritise animal welfare and sustainability.",
  },
  {
    icon: Award,
    title: "Expert Butchers",
    description: "Our team brings decades of experience to every cut, every order.",
  },
];

export function AboutSection() {
  return (
    <section id="about" className="bg-muted/50 py-24">
      <div className="mx-auto max-w-7xl px-6">
        <div className="grid items-center gap-16 lg:grid-cols-2">
          <div>
            <h2 className="font-serif text-3xl font-bold tracking-tight text-foreground sm:text-4xl">
              Quality meat, honest service
            </h2>
            <p className="mt-6 text-lg leading-relaxed text-muted-foreground">
              We believe great food starts with great ingredients. That is why we
              source only the finest halal meats, prepare them with care, and
              deliver them fresh to your family&apos;s table.
            </p>
            <p className="mt-4 text-lg leading-relaxed text-muted-foreground">
              Whether you are cooking a weeknight dinner or hosting a special
              gathering, we are here to provide the cuts you need with the
              quality you expect.
            </p>
          </div>

          <div className="space-y-6">
            {values.map((value) => (
              <div key={value.title} className="flex gap-4">
                <div className="flex size-12 shrink-0 items-center justify-center rounded-lg bg-primary/10">
                  <value.icon className="size-6 text-primary" />
                </div>
                <div>
                  <h3 className="font-semibold text-foreground">
                    {value.title}
                  </h3>
                  <p className="mt-1 text-sm leading-relaxed text-muted-foreground">
                    {value.description}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
