import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { PageHeader } from "@/components/layout/page-header";
import { AppearanceSettings } from "@/components/profile/appearance-settings";
import { ChangePasswordForm } from "@/components/profile/change-password-form";
import { MFASettings } from "@/components/profile/mfa-settings";
import { ProfileForm } from "@/components/profile/profile-form";
import { useAuth } from "@/stores/auth";

export const Route = createFileRoute("/_app/profile")({
  component: ProfilePage,
});

function ProfilePage() {
  const { t } = useTranslation();
  const { user, setUser } = useAuth();

  return (
    <div className="mx-auto w-full max-w-xl px-5 py-6">
      <PageHeader title={t("profile.title")} />
      <div className="space-y-6">
        <ProfileForm user={user} setUser={setUser} />
        <div className="h-px bg-border-secondary" />
        <MFASettings />
        <div className="h-px bg-border-secondary" />
        <AppearanceSettings />
        <div className="h-px bg-border-secondary" />
        <ChangePasswordForm />
      </div>
    </div>
  );
}
