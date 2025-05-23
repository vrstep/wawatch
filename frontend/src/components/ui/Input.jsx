import React from "react";
import cn from "../../utils/cn"; // Adjust path if necessary: from ui/ to utils/ is ../../

const Input = React.forwardRef(({ className, type, ...props }, ref) => {
  return (
    <input
      type={type}
      className={cn(
        "flex h-10 w-full rounded-md border border-slate-300 bg-transparent px-3 py-2 text-sm placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 dark:border-slate-700 dark:text-slate-50 dark:placeholder:text-slate-500 dark:focus:ring-indigo-500 dark:focus:border-indigo-500",
        className
      )}
      ref={ref}
      {...props}
    />
  );
});
Input.displayName = "Input";
export default Input;
