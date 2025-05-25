import './global.css';
import { RootProvider } from 'fumadocs-ui/provider';
import { Outfit, Fira_Code } from 'next/font/google';
import type { ReactNode } from 'react';
import { myfont } from "./assets/fonts"
import '@/app/monaco-config'

const outfit = Outfit({
  weight: ["100", "200", "300", "400", "500", "600", "700", "800", "900"],
  subsets: ["latin"],
  variable: "--font-outfit",
});

const firaCode = Fira_Code({
  weight: ["400", "500", "700"],
  subsets: ["latin"],
  variable: "--font-firacode",
});

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <html lang="en" className={`${outfit.variable} ${myfont.variable}   ${firaCode.variable} antialiased scroll-smooth`} suppressHydrationWarning>
      <body className="flex flex-col min-h-screen scrollbar-thin scrollbar-thumb-[#263238] scrollbar-track-[#001220] selection:bg-red-200/35 selection:text-[#A93E39] dark:selection:bg-red-400/8">
        <RootProvider>
            {children}
        </RootProvider>
      </body>
    </html>
  );
}
