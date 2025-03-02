"use client";
import React from "react";
import { TextHoverEffect } from "@/app/components/ui/text-hover-effect";
import { motion } from "framer-motion";
import Link from "next/link";
import Image from "next/image";

export default function NotFound() {
  return (
    <div className="relative h-screen w-full flex flex-col items-center justify-center overflow-hidden">
      {/* Background Image with Lower z-index */}
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 1 }}
        className="absolute inset-0 z-0"
      >
        <Image
          src="/hero-bg.svg"
          alt="Hero Background"
          className="w-full h-full object-cover select-none"
          draggable={false}
        />
      </motion.div>

      {/* Enlarged 404 Text Effect */}
      <div className="relative z-10 scale-[1.4] sm:scale-[1.6]">
        <TextHoverEffect text="404" />
      </div>

      {/* Message and Button */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 1 }}
        className="relative z-10 text-center mt-4"
      >
        <p className="font-outfit select-none text-2xl text-gray-200">
          MOWA!! Thappu Page ki Vachesaav. Back ki Vellipo.. 🤦‍♂️
        </p>

        {/* Button */}
        <motion.div
          className="mt-6 flex justify-center gap-4 select-none"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.5, duration: 0.8 }}
        >
          <Link
            className="p-3 bg-[#F9CB43] w-40 text-black font-medium font-outfit text-lg rounded-lg transition hover:bg-[#E0B537]"
            href={"/"}
          >
            Home Page
          </Link>
        </motion.div>
      </motion.div>
    </div>
  );
}
