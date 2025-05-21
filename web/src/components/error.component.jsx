import { AlertCircle } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "./ui/alert";

export default function ErrorComponent({ error }) {
    return (
        <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertTitle>{error?.name || "Error"}</AlertTitle>
            <AlertDescription>
                {error?.message}
            </AlertDescription>
            <AlertDescription>
                {`${error?.cause}`}
            </AlertDescription>
        </Alert>
    )
}
