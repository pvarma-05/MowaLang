"use client";
import Link from "next/link";
import { useEffect, useRef, useState } from "react";
import gsap from "gsap";
import SplitType from "split-type";
import { TextPlugin } from "gsap/TextPlugin";
import { ScrollTrigger } from "gsap/ScrollTrigger";
import Playground, { PlaygroundHandle } from "@/components/Playground";
import Footer from "@/components/Footer";

gsap.registerPlugin(TextPlugin, ScrollTrigger);

export default function Home() {
  const titleRef = useRef<HTMLHeadingElement>(null);
  const subtitleRef = useRef<HTMLParagraphElement>(null);
  const buttonsRef = useRef<HTMLDivElement>(null);
  const playgroundRef = useRef<PlaygroundHandle>(null);
  const playgroundSectionRef = useRef<HTMLDivElement>(null);
  const [isPlaygroundVisible, setIsPlaygroundVisible] = useState(false);

  useEffect(() => {
    if (!titleRef.current || !subtitleRef.current) return;

    document.fonts.ready.then(() => {
      const split = new SplitType(titleRef.current!, { types: "chars" });

      // Initial state
      gsap.set(titleRef.current, { visibility: "visible" });
      gsap.set(split.chars, {
        autoAlpha: 0,
        y: 50,
        rotationX: -90,
        transformOrigin: "50% 50%",
      });
      gsap.set(subtitleRef.current, { autoAlpha: 0, visibility: "visible" });
      gsap.set(buttonsRef.current, { autoAlpha: 0, y: 30 });

      const tl = gsap.timeline();

      tl.to(split.chars, {
        autoAlpha: 1,
        y: 0,
        rotationX: 0,
        ease: "back.out(1.7)",
        duration: 0.8,
        stagger: 0.05,
      })
        .to(
          subtitleRef.current,
          {
            autoAlpha: 1,
            duration: 0.8,
            ease: "power2.out",
          },
          "-=0.4"
        )
        .to(
          subtitleRef.current,
          {
            text: {
              value:
                "A programming language with Telugu syntax and Movie Dialogues.",
            },
            duration: 1.8,
            ease: "none",
          },
          "-=0.6"
        )
        .to(
          buttonsRef.current,
          {
            autoAlpha: 1,
            y: 0,
            duration: 0.6,
            ease: "power2.out",
          },
          "-=0.8"
        );
    });

    // Lazy load playground on scroll
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting && !isPlaygroundVisible) {
            setIsPlaygroundVisible(true);
          }
        });
      },
      { threshold: 0.1, rootMargin: "100px" }
    );

    if (playgroundSectionRef.current) {
      observer.observe(playgroundSectionRef.current);
    }

    return () => observer.disconnect();
  }, [isPlaygroundVisible]);

  return (
    <div className="flex flex-col min-h-screen">
      {/* Hero Section */}
      <section className="relative w-full h-screen flex items-center justify-center overflow-hidden">
        {/* Animated background elements */}
        <div className="absolute inset-0 overflow-hidden pointer-events-none">
          <div className="absolute top-20 left-10 w-64 h-64 bg-[#A93E39] opacity-5 rounded-full blur-3xl animate-float"></div>
          <div className="absolute bottom-20 right-10 w-96 h-96 bg-[#A93E39] opacity-5 rounded-full blur-3xl animate-float-delayed"></div>
        </div>

        <div className="relative text-center flex flex-col gap-6 px-4 z-10">
          <h1
            ref={titleRef}
            className="font-wg text-[#A93E39] text-7xl md:text-[100px] lg:text-[155px] select-none invisible tracking-wide drop-shadow-sm"
          >
            MOwa-Lang
          </h1>
          <p
            ref={subtitleRef}
            className="font-outfit select-none text-lg md:text-xl lg:text-2xl text-black/80 min-h-[2.5rem] invisible max-w-3xl mx-auto"
          ></p>
          <div
            ref={buttonsRef}
            className="mt-8 flex flex-wrap justify-center items-center gap-4 md:gap-6 select-none"
          >
            <Link
              className="group relative px-8 md:px-12 py-4 bg-[#A93E39] text-white font-semibold font-outfit text-md md:text-lg rounded-lg transition-all duration-300 hover:shadow-xl hover:-translate-y-1 overflow-hidden"
              href="#playground"
            >
              <span className="relative z-10">Try It</span>
              <div className="absolute inset-0 bg-gradient-to-r from-[#C94C47] to-[#A93E39] opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
            </Link>
            <Link
              className="group px-8 md:px-12 py-4 bg-white text-[#A93E39] font-semibold font-outfit text-md md:text-lg rounded-lg border-2 border-[#A93E39] transition-all duration-300 hover:bg-[#A93E39] hover:text-white hover:shadow-xl hover:-translate-y-1"
              href="/docs"
            >
              Docs
            </Link>
            <Link
              className="group px-8 md:px-12 py-4 bg-white text-[#A93E39] font-semibold font-outfit text-md md:text-lg rounded-lg border-2 border-[#A93E39] transition-all duration-300 hover:bg-[#A93E39] hover:text-white hover:shadow-xl hover:-translate-y-1"
              href="https://github.com/pvarma-05/MowaLang"
              target="_blank"
            >
              Source
            </Link>
          </div>
        </div>

        {/* Scroll indicator */}
        <div className="absolute bottom-8 left-1/2 transform -translate-x-1/2 animate-bounce">
          <svg
            className="w-6 h-6 text-[#A93E39]"
            fill="none"
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="2"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path d="M19 14l-7 7m0 0l-7-7m7 7V3"></path>
          </svg>
        </div>
      </section>

      {/* Playground Section */}
      <div
        id="playground"
        ref={playgroundSectionRef}
        className="min-h-screen relative"
      >
        {isPlaygroundVisible ? (
          <Playground ref={playgroundRef} />
        ) : (
          <div className="w-full h-screen flex items-center justify-center">
            <div className="animate-pulse text-[#A93E39] text-xl font-outfit">
              Loading playground...
            </div>
          </div>
        )}
      </div>

      <Footer />
    </div>
  );
}
