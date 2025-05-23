const SectionTitle = ({ title, icon: Icon }) => (
  <h2 className="text-2xl font-semibold text-slate-800 dark:text-slate-100 mb-6 flex items-center">
    {Icon && <Icon className="mr-3 h-7 w-7 text-indigo-500" />}
    {title}
  </h2>
);

export default SectionTitle;
