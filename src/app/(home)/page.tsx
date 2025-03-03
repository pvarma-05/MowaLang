"use client";
import Link from "next/link";
import { useState } from "react";
import { TextGenerateEffect } from "../components/ui/text-generate-effect";
import Image from "next/image";
import { FeaturesSectionDemo } from "../components/ui/features";
import { motion } from "framer-motion";
// import Footer from "./components/Footer";

const words = `MOWA-LANG`;
const translations = {
  en: {
    heading: (
      <>
        What is <span className="underline underline-offset-8">MOWA-LANG</span>?
      </>
    ),
    description:
      "MOWA-LANG is a fun and quirky Telugu-based programming language written in GO. Designed to make coding more entertaining and engaging, with a unique syntax. Whether you're a beginner looking for an enjoyable way to learn coding or an experienced developer seeking something new, MOWA-LANG is here to add excitement to your coding journey!",
    whyMowa: (
      <>
        Why <span className="underline underline-offset-8">MOWA-LANG</span>?
      </>
    ),
  },
  mowa: {
    heading: (
      <>
        <span className="underline underline-offset-8">MOWA-LANG</span> Antey?
      </>
    ),
    description:
      "MOWA-LANG oka masth fun telugu basha programming language ra! Idi GO lo rayabadindi. Kodtha syntax tho coding chaala enjoy ga untundi! Beginner aina sare, coding chadhavadaniki idi next level ra!",
    whyMowa: (
      <>
        Endhuku <span className="underline underline-offset-8">MOWA-LANG</span>?
      </>
    ),
  },
};

// const words2 = `Raa Mowa Coding Nerchukundham...`;

export default function Home() {
  const [language, setLanguage] = useState<"en" | "mowa">("mowa");
  return (
    <div className="flex flex-col">

      <motion.div
        initial={{ opacity: 0, x: -50 }}
        animate={{ opacity: 1, x: 0 }}
        transition={{ duration: 0.8 }}
        className="fixed top-5 left-5 z-50"
      >
        <button
          className="px-4 py-2 bg-[#F9CB43] cursor-pointer font-outfit text-black font-bold rounded-lg"
          onClick={() => setLanguage(language === "en" ? "mowa" : "en")}
        >
          {language === "en" ? "Switch to Mowa" : "Switch to English"}
        </button>
      </motion.div>

      <section className="relative w-full h-screen flex items-center justify-center text-white">

        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 1 }}
          className="absolute inset-0"
        >
          <Image
            src="/hero-bg.svg"
            alt="Hero Background"
            height={100}
            width={100}
            className="w-full h-full object-cover select-none"
            draggable={false}
          />
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 1 }}
          className="relative z-10 text-center"
        >
          <TextGenerateEffect className="font-gw text-9xl select-none text-[#F9CB43]" duration={2} words={words} />
          <p className="font-outfit select-none text-2xl text-gray-200 mt-2">Raa Mowa Coding Chedhaam 🧑🏻‍💻</p>

          <motion.div
            className="mt-6 flex justify-center gap-4 select-none"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.5, duration: 0.8 }}
          >
            <Link className='p-4 bg-[#F9CB43] w-40 text-black font-medium font-outfit text-xl rounded-lg transition' href={"/play"}>Try It</Link>
            <Link className='p-4 bg-[#F9CB43] w-40 text-black font-medium font-outfit text-xl rounded-lg transition' href={"/docs"}>Docs</Link>
          </motion.div>
        </motion.div>
      </section>
      <hr className="border-7 border-[#F9CB43]" />

      <section className="h-screen flex flex-col-reverse lg:flex-row items-center justify-center px-5 md:gap-10 lg:mx-20">
        <motion.div
          initial={{ opacity: 0, x: -50 }}
          whileInView={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.8 }}
          viewport={{ once: true }}
          className="w-full lg:w-2/4 flex flex-col gap-6 lg:gap-12"
        >
          <h1 className="text-3xl md:text-5xl lg:text-6xl font-gw select-none sm:text-center lg:text-left">
            {translations[language].heading}
          </h1>
          <p className="text-lg lg:text-xl text-gray-200 font-outfit select-none">
            <span className="text-[#F9CB43] font-semibold sm:text-left">MOWA-LANG</span> {translations[language].description}
          </p>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, x: 50 }}
          whileInView={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.8 }}
          viewport={{ once: true }}
          className="right w-full lg:w-2/4 flex justify-center"
        >
          <Image
            width={500}
            height={500}
            src={'/what.svg'}
            draggable={false}
            alt="what is mowa-lang"
            className="max-w-full h-auto select-none"
          />
        </motion.div>
      </section>

      <section className="h-auto min-h-screen flex flex-col items-center justify-center px-5 md:px-10 lg:px-20 py-10">
        <motion.h1
          initial={{ opacity: 0, y: -20 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
          viewport={{ once: true }}
          className="text-3xl md:text-5xl lg:text-6xl font-gw select-none sm:text-center"
        >
          {translations[language].whyMowa}
        </motion.h1>

        <motion.div
          initial="hidden"
          whileInView="visible"
          variants={{
            hidden: { opacity: 0, y: 20 },
            visible: { opacity: 1, y: 0, transition: { staggerChildren: 0.3 } }
          }}
          viewport={{ once: true }}
          className="w-full mt-8"
        >
          <FeaturesSectionDemo />
        </motion.div>
      </section>

    </div>
  );
}
