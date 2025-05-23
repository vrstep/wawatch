import React, { useState } from "react";
// import { useAuth } from "../contexts/AuthContext"; // Add this line
import Label from "../components/ui/Label";
import Input from "../components/ui/Input";
import Button from "../components/ui/Button";
import LoadingSpinner from "../components/ui/LoadingSpinner";
import ErrorMessage from "../components/ui/ErrorMessage";
import apiService from "../api/api";
import { CheckCircle } from "lucide-react";

const SignupPage = () => {
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);
  const [loading, setLoading] = useState(false);
  // const { login } = useAuth();

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (password !== confirmPassword) {
      setError("Passwords do not match.");
      return;
    }
    setError(null);
    setSuccess(null);
    setLoading(true);
    try {
      await apiService.signup({ username, email, password });
      setSuccess("Account created successfully! Please log in.");
      // Optionally auto-login or redirect to login
      // await login({ username, password });
      // window.location.hash = '#/';
      setTimeout(() => (window.location.hash = "#/login"), 2000);
    } catch (err) {
      setError(err.message || "Failed to create account.");
    } finally {
      setLoading(false);
    }
  };
  return (
    <div className="flex items-center justify-center min-h-[calc(100vh-4rem)] py-12 px-4 sm:px-6 lg:px-8">
      <div className="w-full max-w-md space-y-8 p-10 bg-slate-100 dark:bg-slate-800 rounded-xl shadow-2xl">
        <div>
          <h2 className="mt-6 text-center text-3xl font-bold tracking-tight text-slate-900 dark:text-slate-50">
            Create a new account
          </h2>
        </div>
        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          {error && <ErrorMessage message={error} />}
          {success && (
            <div
              className="p-4 mb-4 text-sm text-green-700 bg-green-100 rounded-lg dark:bg-green-200 dark:text-green-800"
              role="alert"
            >
              <CheckCircle className="inline w-5 h-5 mr-2" />
              {success}
            </div>
          )}
          <div className="rounded-md shadow-sm">
            <div className="mb-4">
              <Label htmlFor="username-signup">Username</Label>
              <Input
                id="username-signup"
                name="username"
                type="text"
                required
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder="Username"
              />
            </div>
            <div className="mb-4">
              <Label htmlFor="email-signup">Email address</Label>
              <Input
                id="email-signup"
                name="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="Email address (optional)"
              />
            </div>
            <div className="mb-4">
              <Label htmlFor="password-signup">Password</Label>
              <Input
                id="password-signup"
                name="password"
                type="password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Password"
              />
            </div>
            <div>
              <Label htmlFor="confirm-password-signup">Confirm Password</Label>
              <Input
                id="confirm-password-signup"
                name="confirm-password"
                type="password"
                required
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                placeholder="Confirm Password"
              />
            </div>
          </div>
          <div>
            <Button type="submit" className="w-full" disabled={loading}>
              {loading ? <LoadingSpinner size={20} className="mr-2" /> : null}{" "}
              Create Account
            </Button>
          </div>
        </form>
        <p className="mt-2 text-center text-sm text-slate-600 dark:text-slate-400">
          Already have an account?{" "}
          <a
            href="#/login"
            className="font-medium text-indigo-600 hover:text-indigo-500 dark:text-indigo-400 dark:hover:text-indigo-300"
          >
            Sign in
          </a>
        </p>
      </div>
    </div>
  );
};

export default SignupPage;
