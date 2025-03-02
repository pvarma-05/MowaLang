import type { Metadata } from "next";
import { Outfit } from "next/font/google";
import "./globals.css";
import { myfont } from "./assets/fonts";


const outfit = Outfit({
  weight: ["100", "200", "300", "400", "500", "600", "700", "800", "900"],
  subsets: ["latin"],
  variable: "--font-outfit",
});

export const metadata: Metadata = {
  title: "MowaLang",
  description: "A Telugu Based Toy Programming Language",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="scrollbar-thin scrollbar-thumb-[#263238] scrollbar-track-[#001220]">
      <body
        className={`${outfit.variable} ${myfont.variable} antialiased`}
      >
        {children}
      </body>
    </html>
  );
}
