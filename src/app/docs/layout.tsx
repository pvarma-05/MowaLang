import type { Metadata } from 'next';
import { DocsLayout } from 'fumadocs-ui/layouts/docs';
import type { ReactNode } from 'react';
import { baseOptions } from '@/app/layout.config';
import { source } from '@/lib/source';

export const metadata: Metadata = {
  title: {
    template: 'MowaLang Docs | %s',
    default: 'MowaLang Docs'
  },
  description: "Documentation for MowaLang - A Telugu Based Toy Programming Language",
}

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <div className='font-outfit'>
      <DocsLayout tree={source.pageTree} {...baseOptions}>
        {children}
      </DocsLayout>
    </div>
  );
}
