import { Loader } from "lucide-react";
import cn from "../../utils/cn"; // Adjust path if necessary: from ui/ to utils/ is ../../

const LoadingSpinner = ({ size = 24, className }) => (
  <Loader
    className={cn("animate-spin text-indigo-500", className)}
    size={size}
  />
);

export default LoadingSpinner;
