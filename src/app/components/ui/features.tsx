import { cn } from "@/lib/utils";
import {
  IconLanguage,
  IconDeviceDesktopCode,
  IconWorld,
  IconMessageDots,
} from "@tabler/icons-react";
import { motion } from "framer-motion";

export function FeaturesSectionDemo() {
  const features = [
    {
      title: "Fun Telugu Syntax",
      description:
        "Code like you speak! MOWA-LANG uses Telugu slang to make programming more enjoyable and natural.",
      icon: <IconLanguage className="w-12 h-12 text-[#F9CB43]" />,
    },
    {
      title: "CLI & Web Execution",
      description:
        "Run your code seamlessly on both the command line and an interactive web playground.",
      icon: <IconDeviceDesktopCode className="w-12 h-12 text-[#F9CB43]" />,
    },
    {
      title: "Cross-Platform Support",
      description:
        "MOWA-LANG runs on Windows, macOS, and Linux, so you can code anywhere, anytime.",
      icon: <IconWorld className="w-12 h-12 text-[#F9CB43]" />,
    },
    {
      title: "Movie Dialogues for Errors",
      description:
        "Get error and success messages in iconic Telugu movie dialogues, customizable by your favorite hero!",
      icon: <IconMessageDots className="w-12 h-12 text-[#F9CB43]" />,
    },
  ];

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.2,
      },
    },
  };

  return (
    <motion.div
      className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 relative z-10 py-12 max-w-7xl mx-auto font-outfit"
      variants={containerVariants}
      initial="hidden"
      whileInView="visible"
      viewport={{ once: true, margin: "-100px" }}
    >
      {features.map((feature, index) => (
        <Feature key={feature.title} {...feature} index={index} />
      ))}
    </motion.div>
  );
}

const Feature = ({
  title,
  description,
  icon,
  index,
}: {
  title: string;
  description: string;
  icon: React.ReactNode;
  index: number;
}) => {
  const cardVariants = {
    hidden: { opacity: 0, y: 50, scale: 0.95 },
    visible: {
      opacity: 1,
      y: 0,
      scale: 1,
      transition: {
        duration: 0.6,
        ease: "easeOut",
        type: "spring",
        bounce: 0.3
      }
    },
    hover: {
      scale: 1.05,
      y: -10,
      boxShadow: "0 15px 30px rgba(249, 203, 67, 0.2)",
      transition: {
        duration: 0.3,
        type: "spring",
        stiffness: 300
      }
    }
  };

  const iconVariants = {
    hidden: { scale: 0, rotate: -180 },
    visible: {
      scale: 1,
      rotate: 0,
      transition: {
        duration: 0.5,
        delay: 0.2
      }
    },
    hover: {
      rotate: 360,
      transition: {
        duration: 0.8
      }
    }
  };

  return (
    <motion.div
      variants={cardVariants}
      initial="hidden"
      whileInView="visible"
      whileHover="hover"
      viewport={{ once: true }}
      className={cn(
        "flex flex-col p-6 bg-[#002b45]/30 rounded-xl border border-[#F9CB43]/20",
        "backdrop-blur-md relative overflow-hidden group/feature",
        "transition-all duration-300 select-none"
      )}
    >
      {/* Background glow effect */}
      <motion.div
        className="absolute inset-0 bg-gradient-to-br from-[#F9CB43]/10 to-transparent opacity-0 group-hover/feature:opacity-100"
        transition={{ duration: 0.3 }}
      />

      {/* Icon container */}
      <motion.div
        variants={iconVariants}
        className="mb-6 flex justify-center z-10"
      >
        {icon}
      </motion.div>

      {/* Title */}
      <div className="text-xl font-bold mb-3 relative z-10 text-center">
        <motion.span
          className="inline-block text-neutral-100 group-hover/feature:text-[#F9CB43]"
          initial={{ x: -20, opacity: 0 }}
          whileInView={{ x: 0, opacity: 1 }}
          transition={{ duration: 0.5, delay: 0.3 }}
          viewport={{ once: true }}
        >
          {title}
        </motion.span>
        <motion.div
          className="h-1 w-16 bg-[#F9CB43] rounded-full mx-auto mt-2"
          initial={{ scaleX: 0 }}
          whileInView={{ scaleX: 1 }}
          transition={{ duration: 0.4, delay: 0.4 }}
          viewport={{ once: true }}
        />
      </div>

      {/* Description */}
      <motion.p
        className="text-sm text-neutral-300 text-center z-10"
        initial={{ opacity: 0 }}
        whileInView={{ opacity: 1 }}
        transition={{ duration: 0.5, delay: 0.5 }}
        viewport={{ once: true }}
      >
        {description}
      </motion.p>
    </motion.div>
  );
};
