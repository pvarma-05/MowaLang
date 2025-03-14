import './global.css';
import { RootProvider } from 'fumadocs-ui/provider';
import { Outfit } from 'next/font/google';
import type { ReactNode } from 'react';

import { myfont } from "./assets/fonts";

const outfit = Outfit({
  weight: ["100", "200", "300", "400", "500", "600", "700", "800", "900"],
  subsets: ["latin"],
  variable: "--font-outfit",
});

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <html lang="en" className={`${outfit.variable} ${myfont.variable} antialiased`} suppressHydrationWarning>
      <body className="flex flex-col min-h-screen scrollbar-thin scrollbar-thumb-[#263238] scrollbar-track-[#001220] ">
        <RootProvider>{children}</RootProvider>
      </body>
    </html>
  );
}
