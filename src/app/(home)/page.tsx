"use client";
import Link from "next/link";
import { useState } from "react";
import { TextGenerateEffect } from "../components/ui/text-generate-effect";
import Image from "next/image";
import { FeaturesSectionDemo } from "../components/ui/features";
import { motion } from "framer-motion";

const words = `MOWA-LANG`;
const translations = {
  en: {
    heading: (
      <>
        What is <span className="underline underline-offset-8">MOWA-LANG</span>?
      </>
    ),
    description:
      " is your ticket to coding in Telugu style! Built with GO, this playful language mixes our native slang with programming fun. Perfect for beginners to kickstart coding with a smile or for pros to explore a fresh, entertaining twist. Get ready to code, laugh, and learn – all in mana basha!",
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
      " antey mana Telugu style lo coding ticket mowa! GO tho banchesina ee language, mana slang ni programming tho mix chesi chaala fun isthundi. Koddisepu nerchukunte chaalu – beginners ki easy start, pros ki kotha kick. Mana basha lo code chey, navvu, nerchukov – full josh!",
    whyMowa: (
      <>
        Endhuku <span className="underline underline-offset-8">MOWA-LANG</span>?
      </>
    ),
  },
};

export default function Home() {
  const [language, setLanguage] = useState<"en" | "mowa">("mowa");

  const buttonVariants = {
    hover: {
      scale: 1.05,
      boxShadow: "0px 5px 15px rgba(249, 203, 67, 0.4)",
      transition: { duration: 0.3 },
    },
    tap: { scale: 0.95 },
  };

  const tagline = language === "en"
    ? "Code in Telugu, Fun in Every Line! 🧑🏻‍💻"
    : "Ra Mowa Coding Chedhaam! 🧑🏻‍💻";

  return (
    <div className="flex flex-col min-h-screen">
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="fixed top-5 left-5 z-50"
      >
        <motion.button
          variants={buttonVariants}
          whileHover="hover"
          whileTap="tap"
          className="px-4 py-2 bg-[#F9CB43] cursor-pointer font-outfit text-black font-bold rounded-lg shadow-lg"
          onClick={() => setLanguage(language === "en" ? "mowa" : "en")}
        >
          {language === "en" ? "Switch to Mowa" : "Switch to English"}
        </motion.button>
      </motion.div>

      <section className="relative w-full h-screen flex items-center justify-center text-white overflow-hidden">
        <motion.div
          initial={{ scale: 1.2, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          transition={{ duration: 1.5, ease: "easeOut" }}
          className="absolute inset-0"
        >
          <Image
            src="/hero-bg.svg"
            alt="Hero Background"
            height={100}
            width={100}
            className="w-full h-full object-cover select-none"
            draggable={false}
            priority
          />
          <motion.div
            className="absolute inset-0 "
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.5, duration: 1 }}
          />
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 50 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 1, delay: 0.3 }}
          className="relative z-10 text-center"
        >
          <TextGenerateEffect
            className="font-gw text-6xl md:text-9xl select-none text-[#F9CB43]"
            duration={1.5}
            words={words}
          />
          <motion.p
            className="font-outfit select-none text-xl md:text-2xl text-gray-200 mt-4"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 1, duration: 0.8 }}
          >
            {tagline}
          </motion.p>

          <motion.div
            className="mt-8 flex justify-center gap-6 select-none"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 1.2, duration: 0.8 }}
          >
            <motion.div variants={buttonVariants} whileHover="hover" whileTap="tap">
              <Link
                className="p-4 bg-[#F9CB43] w-40 text-black font-medium font-outfit text-xl rounded-lg transition-all duration-300"
                href={"/play"}
              >
                Try It
              </Link>
            </motion.div>
            <motion.div variants={buttonVariants} whileHover="hover" whileTap="tap">
              <Link
                className="p-4 bg-[#F9CB43] w-40 text-black font-medium font-outfit text-xl rounded-lg transition-all duration-300"
                href={"/docs"}
              >
                Docs
              </Link>
            </motion.div>
          </motion.div>
        </motion.div>
      </section>

      <motion.hr
        className="border-4 border-[#F9CB43]"
        initial={{ width: "0%" }}
        whileInView={{ width: "100%" }}
        transition={{ duration: 1, ease: "easeInOut" }}
        viewport={{ once: true }}
      />

      <section className="min-h-screen flex flex-col-reverse lg:flex-row items-center justify-center px-5 md:gap-10 lg:mx-20 py-20">
        <motion.div
          initial={{ opacity: 0, x: -100 }}
          whileInView={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.8, delay: 0.2 }}
          viewport={{ once: true }}
          className="w-full lg:w-2/4 flex flex-col gap-6 lg:gap-12"
        >
          <h1 className="text-3xl md:text-5xl lg:text-6xl font-gw select-none sm:text-center lg:text-left text-[#F9CB43]">
            {translations[language].heading}
          </h1>
          <motion.p
            className="text-lg lg:text-xl text-gray-200 font-outfit select-none"
            initial={{ opacity: 0 }}
            whileInView={{ opacity: 1 }}
            transition={{ duration: 0.8, delay: 0.4 }}
            viewport={{ once: true }}
          >
            <span className="text-[#F9CB43] font-semibold">MOWA-LANG</span>{" "}
            {translations[language].description}
          </motion.p>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, scale: 0.8 }}
          whileInView={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.8, type: "spring" }}
          viewport={{ once: true }}
          className="w-full lg:w-2/4 flex justify-center mb-10 lg:mb-0"
          whileHover={{ scale: 1.05 }}
        >
          <Image
            width={500}
            height={500}
            src="/what.svg"
            draggable={false}
            alt="what is mowa-lang"
            className="max-w-full h-auto select-none rounded-lg"
            priority
          />
        </motion.div>
      </section>

      <section className="min-h-screen flex flex-col items-center justify-center px-5 md:px-10 lg:px-20 py-20 bg-[#001220]/95">
        <motion.h1
          initial={{ opacity: 0, y: -50 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
          viewport={{ once: true }}
          className="text-3xl md:text-5xl lg:text-6xl font-gw select-none text-center text-[#F9CB43] mb-12"
        >
          {translations[language].whyMowa}
        </motion.h1>
        <FeaturesSectionDemo />
      </section>
    </div>
  );
}
