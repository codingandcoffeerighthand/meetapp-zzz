import { Loader2 } from "lucide-react"
import { cn } from "@/lib/utils"


export default function LoadingScreen({
    className,
    loadingText = "Loading...",
    spinnerSize = 40,
    variant = "primary",
}) {
    // Determine spinner color based on variant
    const spinnerColor =
        variant === "primary" ? "text-primary" : variant === "secondary" ? "text-secondary" : "text-foreground"

    return (
        <div
            className={cn(
                "fixed inset-0 flex flex-col items-center justify-center bg-background/80 backdrop-blur-sm z-50",
                className,
            )}
            role="status"
            aria-live="polite"
        >
            <Loader2 className={cn("animate-spin", spinnerColor)} size={spinnerSize} />
            {loadingText && <p className="mt-4 text-sm font-medium text-muted-foreground">{loadingText}</p>}
        </div>
    )
}
