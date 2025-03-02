"use client";
import { useEffect } from "react";
import { motion, stagger, useAnimate } from "framer-motion";
import { cn } from "@/lib/utils";

export const TextGenerateEffect = ({
  words,
  className,
  duration = 1,
}: {
  words: string;
  className?: string;
  duration?: number;
}) => {
  const [scope, animate] = useAnimate();
  const lettersArray = words.split("");

  useEffect(() => {
    animate(
      "span",
      { opacity: 1, x: 0 },
      { duration: duration, delay: stagger(duration / 50) }
    );
  }, [scope.current]);

  const renderLetters = () => {
    return (
      <motion.div ref={scope} className={cn("inline-block", className)}>
        {lettersArray.map((letter, idx) => (
          <motion.span
            key={idx}
            className={cn("opacity-0 inline-block", className)}
            initial={{ opacity: 0, x: -5 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{
              duration: duration,
              delay: idx * duration,
            }}
          >
            {letter}
          </motion.span>
        ))}
      </motion.div>
    );
  };

  return <div className={cn("font-light", className)}>{renderLetters()}</div>;
};
