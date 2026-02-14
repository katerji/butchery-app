import { HeroSection } from "@/components/customer/hero-section";
import { CategoriesSection } from "@/components/customer/categories-section";
import { SpecialsSection } from "@/components/customer/specials-section";
import { HowItWorksSection } from "@/components/customer/how-it-works-section";
import { AboutSection } from "@/components/customer/about-section";
import { ContactSection } from "@/components/customer/contact-section";

export default function CustomerHomePage() {
  return (
    <>
      <HeroSection />
      <CategoriesSection />
      <SpecialsSection />
      <HowItWorksSection />
      <AboutSection />
      <ContactSection />
    </>
  );
}
