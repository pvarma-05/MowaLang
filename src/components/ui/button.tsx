import * as React from "react";
import { cn } from "@/lib/utils"; // Utility for merging class names

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {}

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, ...props }, ref) => {
    return (
      <button
        ref={ref}
        className={cn(
          "px-4 py-2 bg-yellow-500 text-black font-medium rounded-lg hover:bg-yellow-600 transition duration-200",
          className
        )}
        {...props}
      />
    );
  }
);

Button.displayName = "Button";
