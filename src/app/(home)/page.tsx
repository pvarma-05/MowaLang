"use client";
import Link from "next/link";
import { useEffect, useRef } from "react";
import gsap from "gsap";
import SplitType from "split-type";
import { TextPlugin } from "gsap/TextPlugin";

gsap.registerPlugin(TextPlugin);

export default function Home() {
  const titleRef = useRef<HTMLHeadingElement>(null);
  const subtitleRef = useRef<HTMLParagraphElement>(null);
  const button1Ref = useRef<HTMLAnchorElement>(null);
  const button2Ref = useRef<HTMLAnchorElement>(null);

  useEffect(() => {
    if (
      !titleRef.current ||
      !subtitleRef.current ||
      !button1Ref.current ||
      !button2Ref.current
    )
      return;

    document.fonts.ready.then(() => {
      const split = new SplitType(titleRef.current!, { types: "chars" });

      // Initial state
      gsap.set(titleRef.current, { visibility: "visible" });
      gsap.set(split.chars, {
        autoAlpha: 0,
        y: 50,
        scale: 0.6,
        transformOrigin: "50% 50%",
      });
      gsap.set(subtitleRef.current, { autoAlpha: 0, visibility: "visible" });
      gsap.set([button1Ref.current, button2Ref.current], {
        autoAlpha: 0,
        y: 30,
        visibility: "visible",
      });

      const tl = gsap.timeline();

      tl.to(split.chars, {
        autoAlpha: 1,
        y: 0,
        scale: 1.2,
        ease: "back.out(3)",
        duration: 0.4,
        stagger: 0.08,
      })
        .to(
          split.chars,
          {
            scale: 1,
            duration: 0.3,
            ease: "elastic.out(1, 0.5)",
            stagger: 0.08,
          },
          "-=0.6"
        )
        .to(
          titleRef.current,
          {
            scale: 1.05,
            duration: 0.4,
            ease: "power1.out",
          },
          "0"
        )
        .to(
          titleRef.current,
          {
            scale: 1,
            duration: 0.4,
            ease: "elastic.out(1, 0.4)",
          },
          "-=0.2"
        )
        .to(
          subtitleRef.current,
          {
            autoAlpha: 1,
            text: {
              value:
                "A programming language with Telugu syntax and Movie Dialogues.",
            },
            duration: 2.5,
            ease: "none",
          },
          "+=0.2"
        )
        .to(
          button1Ref.current,
          {
            autoAlpha: 1,
            y: 0,
            duration: 0.8,
            ease: "power3.out",
          },
          "-=1.5"
        )
        .to(
          button2Ref.current,
          {
            autoAlpha: 1,
            y: 0,
            duration: 0.8,
            ease: "power3.out",
          },
          "-=1.5"
        );
    });
  }, []);






  return (
    <div className="flex flex-col min-h-screen">
      <section className="w-full h-screen flex items-center justify-center text-white overflow-hidden">
        <div className="text-center flex flex-col gap-5">
          <h1
            ref={titleRef}
            className="font-wg text-[#A93E39] text-7xl md:text-[100px] lg:text-[155px] select-none invisible tracking-wide">MOwa-Lang
          </h1>

          <p
            ref={subtitleRef}
            className="font-outfit font-serif select-none text-lg md:text-xl lg:text-2xl text-black min-h-[2.5rem] invisible"
          >
          </p>

          <div className="mt-5 flex justify-center gap-6 select-none">
            <Link
              ref={button1Ref}
              className="px-10 py-4 bg-[#A93E39]  text-white font-medium font-outfit text-md md:text-lg lg:text-xl rounded-md transition-all duration-300 hover:scale-105 invisible"
              href={"/play"}
            >
              Try It
            </Link>
            <Link
              ref={button2Ref}
              className="px-10 py-4 bg-[#A93E39]  text-white font-medium font-outfit text-md md:text-lg lg:text-xl rounded-md transition-all duration-300 hover:scale-105 invisible"
              href={"/docs"}
            >
              Docs
            </Link>
          </div>
        </div>
      </section>
    </div>
  );
}
