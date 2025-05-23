import { clsx } from "clsx";
import { twMerge } from "tailwind-merge";

const cn = (...classes) => classes.filter(Boolean).join(" ");

export default cn;
