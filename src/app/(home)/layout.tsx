import type { Metadata } from 'next';
import { HomeLayout } from 'fumadocs-ui/layouts/home';
import { baseOptions } from '@/app/layout.config';
import type { ReactNode } from 'react';

export const metadata: Metadata = {
  title: "MowaLang",
  description: "A Telugu Based Toy Programming Language",
};

export default function Layout({ children }: { children: ReactNode }) {
  return <HomeLayout {...baseOptions} nav={{ enabled: false }} className='bg-[#001220] !pt-0  '>{children}</HomeLayout>;
}
