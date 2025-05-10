'use client';

import React from 'react';
import { Button } from '@/app/components/ui/button';
import { DownloadIcon } from 'lucide-react';

const downloads = [
  {
    name: 'Windows Installer',
    link: 'https://mowalang.vercel.app/scripts/windows_install.bat',
    os: 'Windows',
  },
  {
    name: 'Linux Installer',
    link: 'https://mowalang.vercel.app/scripts/linux_install.sh',
    os: 'Linux',
  },
  {
    name: 'macOS Installer (Intel)',
    link: 'https://mowalang.vercel.app/scripts/mac_install.sh',
    os: 'macOS Intel',
  },
  // {
  //   name: 'macOS Installer (M1/M2)',
  //   link: 'https://your-site.com/bin/install-mac-arm.sh',
  //   os: 'macOS ARM',
  // },
];

export default function DownloadsPage() {
  return (
    <div className="max-w-4xl mx-auto px-4 py-12">
      <h1 className="text-3xl font-semibold mb-8 text-center">Download MowaLang</h1>
      <div className="grid gap-6 sm:grid-cols-1">
        {downloads.map((item) => (
          <div
            key={item.link}
            className="border rounded-2xl p-6 flex flex-col items-start justify-between gap-4 shadow-sm"
          >
            <div>
              <h2 className="text-xl font-medium mb-1">{item.name}</h2>
              <p className="text-sm text-muted-foreground">For {item.os} users</p>
            </div>
            <Button className="gap-2">
              <a href={item.link} download>
                <DownloadIcon className="w-4 h-4" /> Download
              </a>
            </Button>
          </div>
        ))}
      </div>
    </div>
  );
}
