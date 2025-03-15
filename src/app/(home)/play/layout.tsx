import type { Metadata } from 'next';
import type { ReactNode } from 'react';

export const metadata: Metadata = {
  title: "MowaLang Playground",
  description: "Playgronund for MowaLang - A Telugu Based Toy Programming Language",
};

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <>
    {children}
    </>
  )
}
