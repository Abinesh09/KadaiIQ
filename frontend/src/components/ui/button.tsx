import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"
import { cn } from "@/lib/utils"

const buttonVariants = cva(
  "inline-flex items-center justify-center whitespace-nowrap rounded-lg text-xs font-medium transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#16A36A] disabled:pointer-events-none disabled:opacity-50 select-none cursor-pointer",
  {
    variants: {
      variant: {
        default: "bg-[#172B4D] text-white shadow-2xs hover:bg-[#243B5A] active:scale-[0.99]",
        primary: "bg-[#172B4D] text-white shadow-2xs hover:bg-[#243B5A] active:scale-[0.99]",
        accent: "bg-[#16A36A] text-white shadow-2xs hover:bg-[#138958] active:scale-[0.99]",
        success: "bg-[#16A36A] text-white shadow-2xs hover:bg-[#138958] active:scale-[0.99]",
        destructive: "bg-[#DC2626] text-white shadow-2xs hover:bg-[#B91C1C] active:scale-[0.99]",
        outline: "border border-[#E5EAF0] bg-white text-[#172033] shadow-2xs hover:bg-[#F5F7FA] hover:text-[#172033] active:scale-[0.99]",
        secondary: "bg-[#F5F7FA] text-[#172033] hover:bg-[#E5EAF0] active:scale-[0.99]",
        ghost: "hover:bg-[#F5F7FA] text-[#64748B] hover:text-[#172033]",
        link: "text-[#16A36A] underline-offset-4 hover:underline",
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
