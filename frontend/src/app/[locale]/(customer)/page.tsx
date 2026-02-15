import { HeroSection } from "@/components/customer/hero-section";
import { CategoriesSection } from "@/components/customer/categories-section";
import { SpecialsSection } from "@/components/customer/specials-section";
import { HowItWorksSection } from "@/components/customer/how-it-works-section";
import { AboutSection } from "@/components/customer/about-section";
import { ContactSection } from "@/components/customer/contact-section";

interface Props {
  searchParams: Promise<{ logged_out?: string }>;
}

export default async function CustomerHomePage({ searchParams }: Props) {
  const params = await searchParams;
  const loggedOut = params.logged_out === "true";

  return (
    <>
      <HeroSection showLogoutBanner={loggedOut} />
      <CategoriesSection />
      <SpecialsSection />
      <HowItWorksSection />
      <AboutSection />
      <ContactSection />
    </>
  );
}
