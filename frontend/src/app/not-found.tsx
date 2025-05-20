"use client";
import React from "react";
import { TextHoverEffect } from "@/components/ui/text-hover-effect";
import { motion } from "framer-motion";
import Link from "next/link";

export default function NotFound() {
  return (
    <div className="relative h-screen w-full flex flex-col items-center justify-center overflow-hidden">

      {/* Enlarged 404 Text Effect */}
      <div className="relative z-10 scale-[1.4] sm:scale-[1.6]">
        <TextHoverEffect text="404" />
      </div>

      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 1 }}
        className="relative z-10 text-center mt-4"
      >
        <p className="font-outfit select-none text-2xl text-[#A93E39]">
          MOWA!! Thappu Page ki Vachesaav. Back ki Vellipo.. 🤦‍♂️
        </p>

        <motion.div
          className="mt-6 flex justify-center gap-4 select-none"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.5, duration: 0.8 }}
        >
          <Link
            className="p-3 bg-[#A93E39] w-40 text-white font-medium font-outfit text-lg rounded-lg transition hover:bg-[#E0B537]"
            href={"/"}
          >
            Home Page
          </Link>
        </motion.div>
      </motion.div>
    </div>
  );
}
