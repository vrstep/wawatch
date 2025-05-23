import React, { useState } from "react";
import { useAuth } from "../contexts/AuthContext"; // Add this line
import ErrorMessage from "../components/ui/ErrorMessage";
import LoadingSpinner from "../components/ui/LoadingSpinner";
import Label from "../components/ui/Label";
import Input from "../components/ui/Input";
import Button from "../components/ui/Button";
import apiService from "../api/api";
import SectionTitle from "../components/layout/SectionTitle";
import { Settings } from "lucide-react";
import { useEffect } from "react";

const SettingsPage = () => {
  const { user, loading: authLoading } = useAuth();
  // States for forms
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmNewPassword, setConfirmNewPassword] = useState("");

  const [newUsername, setNewUsername] = useState("");
  const [usernamePassword, setUsernamePassword] = useState("");

  const [newEmail, setNewEmail] = useState("");
  const [emailPassword, setEmailPassword] = useState("");

  const [error, setError] = useState({});
  const [success, setSuccess] = useState({});
  const [formLoading, setFormLoading] = useState({
    password: false,
    username: false,
    email: false,
  });

  if (authLoading)
    return (
      <div className="flex justify-center items-center h-screen">
        <LoadingSpinner size={48} />
      </div>
    );
  if (!user) {
    useEffect(() => {
      window.location.hash = "#/login";
    }, []);
    return null;
  }

  const handlePasswordChange = async (e) => {
    e.preventDefault();
    if (newPassword !== confirmNewPassword) {
      setError((prev) => ({
        ...prev,
        password: "New passwords do not match.",
      }));
      return;
    }
    setError((prev) => ({ ...prev, password: null }));
    setSuccess((prev) => ({ ...prev, password: null }));
    setFormLoading((prev) => ({ ...prev, password: true }));
    try {
      await apiService.changePassword({
        current_password: currentPassword,
        new_password: newPassword,
      });
      setSuccess((prev) => ({
        ...prev,
        password: "Password changed successfully!",
      }));
      setCurrentPassword("");
      setNewPassword("");
      setConfirmNewPassword("");
    } catch (err) {
      setError((prev) => ({
        ...prev,
        password: err.message || "Failed to change password.",
      }));
    } finally {
      setFormLoading((prev) => ({ ...prev, password: false }));
    }
  };
  // Similar handlers for changeUsername and changeEmail

  return (
    <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-8 max-w-2xl">
      <SectionTitle title="Account Settings" icon={Settings} />
      <div className="space-y-8">
        {/* Change Password Form */}
        <form
          onSubmit={handlePasswordChange}
          className="p-6 bg-slate-100 dark:bg-slate-800 rounded-lg shadow"
        >
          <h3 className="text-lg font-medium mb-4 text-slate-800 dark:text-slate-100">
            Change Password
          </h3>
          {error.password && <ErrorMessage message={error.password} />}
          {success.password && (
            <div className="p-3 mb-3 text-sm text-green-700 bg-green-100 rounded-lg dark:bg-green-200 dark:text-green-800">
              {success.password}
            </div>
          )}
          <div className="space-y-4">
            <div>
              <Label htmlFor="currentPassword">Current Password</Label>
              <Input
                type="password"
                id="currentPassword"
                value={currentPassword}
                onChange={(e) => setCurrentPassword(e.target.value)}
                required
              />
            </div>
            <div>
              <Label htmlFor="newPassword">New Password</Label>
              <Input
                type="password"
                id="newPassword"
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                required
              />
            </div>
            <div>
              <Label htmlFor="confirmNewPassword">Confirm New Password</Label>
              <Input
                type="password"
                id="confirmNewPassword"
                value={confirmNewPassword}
                onChange={(e) => setConfirmNewPassword(e.target.value)}
                required
              />
            </div>
          </div>
          <Button
            type="submit"
            className="mt-6"
            disabled={formLoading.password}
          >
            {formLoading.password ? (
              <LoadingSpinner size={20} />
            ) : (
              "Change Password"
            )}
          </Button>
        </form>
        {/* TODO: Change Username Form */}
        {/* TODO: Change Email Form */}
      </div>
    </div>
  );
};
export default SettingsPage;
