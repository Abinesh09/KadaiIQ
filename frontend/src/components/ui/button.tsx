import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"
import { cn } from "@/lib/utils"

const buttonVariants = cva(
  "inline-flex items-center justify-center whitespace-nowrap rounded-lg text-xs font-medium transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#8B5CF6] disabled:pointer-events-none disabled:opacity-50 select-none cursor-pointer",
  {
    variants: {
      variant: {
        default: "bg-[#8B5CF6] text-white shadow-2xs hover:bg-[#7C3AED] active:scale-[0.99]",
        primary: "bg-[#8B5CF6] text-white shadow-2xs hover:bg-[#7C3AED] active:scale-[0.99]",
        accent: "bg-[#8B5CF6] text-white shadow-2xs hover:bg-[#7C3AED] active:scale-[0.99]",
        success: "bg-[#16A36A] text-white shadow-2xs hover:bg-[#138958] active:scale-[0.99]",
        destructive: "bg-[#DC2626] text-white shadow-2xs hover:bg-[#B91C1C] active:scale-[0.99]",
        outline: "border border-[#EDE9FE] bg-white text-[#5B45B5] shadow-2xs hover:bg-[#F5F3FF] active:scale-[0.99]",
        secondary: "bg-[#EDE9FE] text-[#5B45B5] hover:bg-[#F5F3FF] active:scale-[0.99]",
        ghost: "hover:bg-[#F5F3FF] text-[#5B45B5] hover:text-[#7C3AED]",
        link: "text-[#8B5CF6] underline-offset-4 hover:underline",
      },
      size: {
        default: "h-9 px-3.5 py-2",
        sm: "h-8 px-2.5 text-xs rounded-md",
        lg: "h-10 px-5 text-sm rounded-lg",
        icon: "h-9 w-9 p-0",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  }
)

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, ...props }, ref) => {
    return (
      <button
        className={cn(buttonVariants({ variant, size, className }))}
        ref={ref}
        {...props}
      />
    )
  }
)
Button.displayName = "Button"

export { Button, buttonVariants }
