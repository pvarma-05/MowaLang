import './global.css';
import { RootProvider } from 'fumadocs-ui/provider';
import { Outfit,Poppins,Fira_Code } from 'next/font/google';
import type { ReactNode } from 'react';
import { myfont,myfont2,myfont3 } from "./assets/fonts";

const outfit = Outfit({
  weight: ["100", "200", "300", "400", "500", "600", "700", "800", "900"],
  subsets: ["latin"],
  variable: "--font-outfit",
});
const ak = Poppins({
  weight: ["100", "200", "300", "400", "500", "600", "700", "800", "900"],
  subsets: ["latin"],
  variable: "--font-at",
});
const firaCode = Fira_Code({
  weight: ["400", "500", "700"],
  subsets: ["latin"],
  variable: "--font-firacode",
});

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <html lang="en" className={`${outfit.variable} ${myfont.variable} ${myfont2.variable} ${myfont3.variable} ${ak.variable} ${firaCode.variable} antialiased`} suppressHydrationWarning>
      <body className="flex flex-col min-h-screen scrollbar-thin scrollbar-thumb-[#263238] scrollbar-track-[#001220] ">
        <RootProvider theme={{enabled:false}}>{children}</RootProvider>
      </body>
    </html>
  );
}
