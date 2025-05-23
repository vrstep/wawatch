import cn from "../../utils/cn"; // Adjust path if necessary: from ui/ to utils/ is ../../
import { AlertTriangle } from "lucide-react";

const ErrorMessage = ({ message, className }) => (
  <div
    className={cn(
      "p-4 mb-4 text-sm text-red-700 bg-red-100 rounded-lg dark:bg-red-200 dark:text-red-800",
      className
    )}
    role="alert"
  >
    <AlertTriangle className="inline w-5 h-5 mr-2" />
    <span className="font-medium">Error:</span> {message}
  </div>
);
export default ErrorMessage;
