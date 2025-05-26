'use client';

import { useEffect, useRef } from 'react';
import { Button } from '@/components/ui/button';
import { FaApple, FaWindows, FaLinux } from 'react-icons/fa';
import { IconType } from 'react-icons';
import gsap from 'gsap';

interface DownloadItem {
  name: string;
  link: string;
  os: string;
  description: string;
  icon: IconType;
}

const downloads: DownloadItem[] = [
  {
    name: 'Windows Installer',
    link: 'https://mowalang.vercel.app/scripts/windows_install.bat',
    os: 'Windows',
    description: 'Install MowaLang on Windows 10/11 (64-bit).',
    icon: FaWindows,
  },
  {
    name: 'Linux Installer',
    link: 'https://mowalang.vercel.app/scripts/linux_install.sh',
    os: 'Linux',
    description: 'Compatible with Ubuntu, Debian, and other Linux distributions.',
    icon: FaLinux,
  },
  {
    name: 'macOS Installer (M1/M2)',
    link: 'https://mowalang.vercel.app/scripts/mac_install_arm.sh',
    os: 'macOS ARM',
    description: 'Optimized for macOS on Apple Silicon (M1/M2).',
    icon: FaApple,
  },
];

export default function DownloadsPage() {
  const sectionRef = useRef<HTMLElement>(null);
  const titleRef = useRef<HTMLHeadingElement>(null);
  const introRef = useRef<HTMLParagraphElement>(null);
  const noteRef = useRef<HTMLParagraphElement>(null);
  const cardsRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!titleRef.current || !introRef.current || !noteRef.current || !cardsRef.current) return;

    const cards = cardsRef.current.querySelectorAll('.download-card');
    const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;

    gsap.set([titleRef.current, introRef.current, noteRef.current, ...cards], {
      autoAlpha: 0,
      y: 20,
    });

    const tl = gsap.timeline({
      defaults: { ease: 'power2.out', duration: reduceMotion ? 0 : 0.5 },
    });

    tl.to(titleRef.current, { autoAlpha: 1, y: 0 })
      .to(introRef.current, { autoAlpha: 1, y: 0 }, '-=0.2')
      .to(noteRef.current, { autoAlpha: 1, y: 0 }, '-=0.2')
      .to(cards, { autoAlpha: 1, y: 0, stagger: 0.15 }, '-=0.2');
  }, []);

  return (
    <section
      ref={sectionRef}
      className="w-full px-4 py-12 md:py-16 min-h-screen flex flex-col items-center"
      aria-labelledby="downloads-title"
    >
      <div className="container mx-auto max-w-7xl">
        <h1
          ref={titleRef}
          id="downloads-title"
          className="text-4xl md:text-5xl lg:text-6xl font-wg text-[#A93E39] mb-4 text-center select-none"
        >
          DOwnlOad MOwa-Lang
        </h1>
        <p
          ref={introRef}
          className="text-lg font-outfit text-gray-600 mb-2 text-center max-w-3xl mx-auto"
        >
          Get started with MowaLang by downloading the installer for your operating system. Run the script to set up the language and start coding with Telugu syntax and epic movie dialogues!
        </p>
        <br />
        <p
          ref={noteRef}
          className="text-lg font-outfit text-gray-600 mb-12 text-center max-w-3xl mx-auto"
        >
          Make sure to run the scripts in admin or sudo mode to work.
        </p>
        <div
          ref={cardsRef}
          className="grid gap-6 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3"
        >
          {downloads.map((item) => (
            <div
              key={item.link}
              className="download-card backdrop-blur-md bg-white/10 rounded-xl p-6 flex flex-col justify-between gap-4 shadow-lg hover:scale-105 transition-transform duration-300"
              role="region"
              aria-label={`Download ${item.name}`}
            >
              <div>
                <item.icon className="text-3xl text-[#A93E39] mb-2" aria-hidden="true" />
                <h2 className="text-xl font-outfit text-[#A93E39] mb-2 font-semibold">{item.name}</h2>
                <p className="text-sm font-outfit text-gray-600 mb-2">{item.description}</p>
                <p className="text-sm font-outfit text-gray-500">For {item.os}</p>
              </div>
              <Button
                className="select-none w-full bg-[#A93E39] text-white font-outfit font-medium rounded-md py-3 hover:bg-[#7B2D2A] transition-all duration-300 gap-2 focus:ring-2 focus:ring-offset-2 focus:ring-[#7B2D2A]"
              >
                <a
                  href={item.link}
                  download
                  aria-label={`Download MowaLang for ${item.os}`}
                >
                  Download
                </a>
              </Button>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

