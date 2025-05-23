import React, { useState } from "react";
import cn from "../../utils/cn"; // Adjust path if necessary: from ui/ to utils/ is ../../

const Label = React.forwardRef(({ className, ...props }, ref) => (
  <label
    ref={ref}
    className={cn(
      "text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 dark:text-slate-200",
      className
    )}
    {...props}
  />
));

Label.displayName = "Label";
export default Label;
