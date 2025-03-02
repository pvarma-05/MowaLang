import { cn } from "@/lib/utils";
import {
  IconLanguage,
  IconDeviceDesktopCode,
  IconWorld,
  IconMessageDots,
} from "@tabler/icons-react";

export function FeaturesSectionDemo() {
  const features = [
    {
      title: "Fun Telugu Syntax",
      description:
        "Code like you speak! MOWA-LANG uses Telugu slang to make programming more enjoyable and natural.",
      icon: <IconLanguage className="w-10 h-10 text-[#F9CB43]" />,
    },
    {
      title: "CLI & Web Execution",
      description:
        "Run your code seamlessly on both the command line and an interactive web playground.",
      icon: <IconDeviceDesktopCode className="w-10 h-10 text-[#F9CB43]" />,
    },
    {
      title: "Cross-Platform Support",
      description: "MOWA-LANG runs on Windows, macOS, and Linux, so you can code anywhere, anytime.",
      icon: <IconWorld className="w-10 h-10 text-[#F9CB43]" />,
    },
    {
      title: "Movie Dialogues for Errors",
      description:
        "Get error and success messages in iconic Telugu movie dialogues, customizable by your favorite hero!",
      icon: <IconMessageDots className="w-10 h-10 text-[#F9CB43]" />,
    },
  ];
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4  relative z-10 py-10 max-w-7xl mx-auto font-outfit">
      {features.map((feature, index) => (
        <Feature key={feature.title} {...feature} index={index} />
      ))}
    </div>
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
  return (
    <div
      className={cn(
        "flex flex-col lg:border-r  py-10 relative group/feature dark:border-neutral-800",
        (index === 0 || index === 4) && "lg:border-l dark:border-neutral-800",
        index < 4 && "lg:border-b dark:border-neutral-800"
      )}
    >
      {index < 4 && (
        <div className="select-none opacity-0 group-hover/feature:opacity-100 transition duration-200 absolute inset-0 h-full w-full bg-gradient-to-t from-neutral-100 dark:from-neutral-800 to-transparent pointer-events-none" />
      )}
      {index >= 4 && (
        <div className="select-none opacity-0 group-hover/feature:opacity-100 transition duration-200 absolute inset-0 h-full w-full bg-gradient-to-b from-neutral-100 dark:from-neutral-800 to-transparent pointer-events-none" />
      )}
      <div className="select-none mb-4 relative z-10 px-10 text-neutral-600 dark:text-neutral-400">
        {icon}
      </div>
      <div className="select-none text-lg font-bold mb-2 relative z-10 px-10">
        <div className="absolute left-0 inset-y-0 h-6 group-hover/feature:h-8 w-1 rounded-tr-full rounded-br-full bg-neutral-300 dark:bg-neutral-700 group-hover/feature:bg-[#F9CB43] transition-all duration-200 origin-center" />
        <span className="group-hover/feature:translate-x-2 transition duration-200 inline-block text-neutral-800 dark:text-neutral-100">
          {title}
        </span>
      </div>
      <p className="select-none text-sm text-neutral-600 dark:text-neutral-300 max-w-xs relative z-10 px-10">
        {description}
      </p>
    </div>
  );
};
