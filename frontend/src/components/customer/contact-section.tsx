import { MapPin, Clock, Phone } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";

const contactDetails = [
  {
    icon: MapPin,
    title: "Location",
    lines: ["Al Barsha, Dubai", "United Arab Emirates"],
  },
  {
    icon: Clock,
    title: "Opening Hours",
    lines: ["Mon–Fri: 7:00 AM – 6:00 PM", "Sat: 7:00 AM – 4:00 PM", "Sun: Closed"],
  },
  {
    icon: Phone,
    title: "Get in Touch",
    lines: ["+971 4 123 4567", "hello@butchery.ae"],
  },
];

export function ContactSection() {
  return (
    <section id="contact" className="py-24">
      <div className="mx-auto max-w-7xl px-6">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="font-serif text-3xl font-bold tracking-tight text-foreground sm:text-4xl">
            Visit Us
          </h2>
          <p className="mt-4 text-lg text-muted-foreground">
            Come in and see what&apos;s fresh today, or give us a call to place
            an order.
          </p>
        </div>

        <div className="mt-16 grid gap-6 sm:grid-cols-3">
          {contactDetails.map((detail) => (
            <Card key={detail.title} className="border-border/50">
              <CardContent className="flex flex-col items-center p-8 text-center">
                <div className="flex size-12 items-center justify-center rounded-full bg-primary/10">
                  <detail.icon className="size-6 text-primary" />
                </div>
                <h3 className="mt-4 font-semibold text-foreground">
                  {detail.title}
                </h3>
                <div className="mt-2 space-y-1">
                  {detail.lines.map((line) => (
                    <p key={line} className="text-sm text-muted-foreground">
                      {line}
                    </p>
                  ))}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </section>
  );
}
