import Button from "../components/ui/Button";

const NotFoundPage = ({ message = "Page Not Found" }) => (
  <div className="container mx-auto px-4 py-16 text-center">
    <h1 className="text-4xl font-bold text-slate-800 dark:text-slate-100 mb-4">
      404
    </h1>
    <p className="text-xl text-slate-600 dark:text-slate-400">{message}</p>
    <Button onClick={() => (window.location.hash = "#/")} className="mt-8">
      Go Home
    </Button>
  </div>
);
export default NotFoundPage;
