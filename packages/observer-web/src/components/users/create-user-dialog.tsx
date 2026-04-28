import { type SyntheticEvent, useState } from "react";

import { Field } from "@base-ui/react/field";
import { useTranslation } from "react-i18next";

import { FormDialog } from "@/components/dialogs/form-dialog";
import { FormField } from "@/components/forms/form-field";
import { UISelect } from "@/components/ui/ui-select";
import { UISwitch } from "@/components/ui/ui-switch";
import { useOffices } from "@/hooks/reference/use-offices";
import { useCreateUser } from "@/hooks/users/use-users";
import { handleApiError } from "@/lib/form-error";
import { toSelectOptions } from "@/lib/options";

interface CreateUserDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onCreated?: () => void;
}

const defaultForm = {
  first_name: "",
  last_name: "",
  email: "",
  phone: "",
  password: "",
  role: "consultant",
  office_id: "",
  is_active: true,
  is_verified: true,
};

export function CreateUserDialog({ open, onOpenChange, onCreated }: CreateUserDialogProps) {
  const { t } = useTranslation();
  const createUser = useCreateUser();
  const { data: officesData } = useOffices();

  const [form, setForm] = useState(defaultForm);
  const [error, setError] = useState("");

  const roleOptions = [
    { label: t("admin.users.roleAdmin"), value: "admin" },
    { label: t("admin.users.roleStaff"), value: "staff" },
    { label: t("admin.users.roleConsultant"), value: "consultant" },
    { label: t("admin.users.roleGuest"), value: "guest" },
  ];

  const officeOptions = [{ label: "—", value: "" }, ...toSelectOptions(officesData)];

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    setError("");
    try {
      await createUser.mutateAsync({
        first_name: form.first_name,
        last_name: form.last_name || undefined,
        email: form.email,
        phone: form.phone || undefined,
        password: form.password,
        role: form.role,
        office_id: form.office_id || null,
        is_active: form.is_active,
        is_verified: form.is_verified,
      });
      onOpenChange(false);
      setForm(defaultForm);
      onCreated?.();
    } catch (err) {
      setError(await handleApiError(err, t));
    }
  }

  return (
    <FormDialog
      open={open}
      onOpenChange={onOpenChange}
      title={t("admin.users.addTitle")}
      loading={createUser.isPending}
      onSubmit={handleSubmit}
      maxWidth="md"
      error={error}
    >
      <div className="grid grid-cols-2 gap-3">
        <FormField
          label={t("admin.users.firstName")}
          required
          value={form.first_name}
          onChange={(v) => setForm((f) => ({ ...f, first_name: v }))}
        />
        <FormField
          label={t("admin.users.lastName")}
          value={form.last_name}
          onChange={(v) => setForm((f) => ({ ...f, last_name: v }))}
        />
      </div>

      <FormField
        label={t("admin.users.email")}
        type="email"
        required
        value={form.email}
        onChange={(v) => setForm((f) => ({ ...f, email: v }))}
      />

      <FormField
        label={t("admin.users.phone")}
        type="tel"
        value={form.phone}
        onChange={(v) => setForm((f) => ({ ...f, phone: v }))}
      />

      <FormField
        label={t("admin.users.password")}
        type="password"
        required
        value={form.password}
        onChange={(v) => setForm((f) => ({ ...f, password: v }))}
      />

      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("admin.users.role")}
        </Field.Label>
        <UISelect
          value={form.role}
          onValueChange={(v) => setForm((f) => ({ ...f, role: v }))}
          options={roleOptions}
          fullWidth
        />
      </Field.Root>

      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("admin.users.office")}
        </Field.Label>
        <UISelect
          value={form.office_id}
          onValueChange={(v) => setForm((f) => ({ ...f, office_id: v }))}
          options={officeOptions}
          placeholder="—"
          fullWidth
        />
      </Field.Root>

      <div className="flex gap-6">
        <UISwitch
          checked={form.is_active}
          onCheckedChange={(v) => setForm((f) => ({ ...f, is_active: v }))}
          label={t("admin.users.active")}
        />
        <UISwitch
          checked={form.is_verified}
          onCheckedChange={(v) => setForm((f) => ({ ...f, is_verified: v }))}
          label={t("admin.users.verified")}
        />
      </div>
    </FormDialog>
  );
}
