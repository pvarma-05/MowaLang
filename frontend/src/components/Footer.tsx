"use client";

import { useEffect, useRef } from "react";
import { Github, Mail, BookOpen, Code } from "lucide-react";
import gsap from "gsap";

export default function Footer() {
  const footerRef = useRef<HTMLElement>(null);
  const contentRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!footerRef.current || !contentRef.current) return;

    const reduceMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;

    gsap.set(contentRef.current, { autoAlpha: 0, y: 30 });

    gsap.to(contentRef.current, {
      autoAlpha: 1,
      y: 0,
      duration: reduceMotion ? 0 : 1,
      ease: "power2.out",
    });
  }, []);

  return (
    <footer
      ref={footerRef}
      className="w-full bg-gradient-to-t bg-[#A93E39] py-12 select-none"
      aria-labelledby="footer-title"
    >
      <div className="container mx-auto max-w-7xl px-4">
        <div ref={contentRef} className="flex flex-col gap-10 md:flex-row justify-between">
          <div className="flex flex-col gap-4 max-w-2xl">
            <h3 className="font-wg text-5xl  text-[#FDF2EC] select-none">
              MOwa-Lang
            </h3>
            <p className="font-outfit text-white text-md">
              MowaLang is a dynamically typed, Telugu-based programming language built with Go. It lets you write expressive programs using Telugu syntax and shows iconic Telugu movie dialogues when your code succeeds or fails.
            </p>
          </div>
          <div className="flex flex-col gap-4">
            <h3 className="font-outfit text-2xl font-semibold text-[#FDF2EC] select-none">
              Quick Links
            </h3>
            <ul className="flex flex-col gap-2">
              <li>
                <a
                  href="#playground"
                  className="font-outfit text-md text-white hover:text-gray-300 transition-colors duration-300 flex items-center gap-2 focus:outline-none focus:ring-2 focus:ring-[#7B2D2A]"
                  aria-label="MowaLang playground"
                >
                  <Code size={18} aria-hidden="true" />
                  Playground
                </a>
              </li>
              <li>
                <a
                  href="/docs"
                  className="font-outfit text-md text-white hover:text-gray-300 transition-colors duration-300 flex items-center gap-2 focus:outline-none focus:ring-2 focus:ring-[#7B2D2A]"
                  aria-label="MowaLang documentation"
                >
                  <BookOpen size={18} aria-hidden="true" />
                  Docs
                </a>
              </li>
              <li>
                <a
                  href="https://github.com/pvarma-05/MowaLang"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="font-outfit text-md text-white hover:text-gray-300 transition-colors duration-300 flex items-center gap-2 focus:outline-none focus:ring-2 focus:ring-[#7B2D2A]"
                  aria-label="MowaLang GitHub repository"
                >
                  <Github size={18} aria-hidden="true" />
                  GitHub
                </a>
              </li>
            </ul>
          </div>
        </div>

        <hr className="my-8 border-gray-300" />

        <div className="text-center">
          <p className="font-outfit text-md text-white">
            &copy; {new Date().getFullYear()} MowaLang. Mowa, all rights reserved!
          </p>
        </div>
      </div>
    </footer>
  );
}
