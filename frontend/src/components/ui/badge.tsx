import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"
import { cn } from "@/lib/utils"

const badgeVariants = cva(
  "inline-flex items-center rounded-md border px-2 py-0.5 text-xs font-medium transition-colors focus:outline-none",
  {
    variants: {
      variant: {
        default:
          "border-transparent bg-[#172B4D] text-white shadow-2xs hover:bg-[#243B5A]",
        primary:
          "border-transparent bg-[#172B4D] text-white",
        secondary:
          "border-[#E5EAF0] bg-[#F5F7FA] text-[#172033]",
        destructive:
          "border-[#DC2626]/20 bg-[#DC2626]/10 text-[#DC2626] font-semibold",
        outline:
          "border-[#E5EAF0] text-[#64748B] bg-white",
        success:
          "border-[#16A36A]/20 bg-[#16A36A]/10 text-[#16A36A] font-semibold",
        warning:
          "border-[#D97706]/20 bg-[#D97706]/10 text-[#D97706] font-semibold",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  }
)

export interface BadgeProps
  extends React.HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof badgeVariants> {}

function Badge({ className, variant, ...props }: BadgeProps) {
  return (
    <div className={cn(badgeVariants({ variant }), className)} {...props} />
  )
}

export { Badge, badgeVariants }
