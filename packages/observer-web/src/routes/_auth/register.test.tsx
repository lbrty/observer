import { afterEach, describe, expect, it, mock } from "bun:test";
import { cleanup } from "@testing-library/react";

let CapturedComponent: any;

mock.module("react-i18next", () => ({
  useTranslation: () => ({ t: (k: string) => k }),
}));

mock.module("@tanstack/react-router", () => ({
  createFileRoute: () => (opts: any) => {
    CapturedComponent = opts.component;
    return opts;
  },
  Link: ({ children, ...props }: any) => <a {...props}>{children}</a>,
  useNavigate: () => () => {},
  useRouter: () => ({}),
  useRouterState: () => ({}),
  useMatch: () => ({}),
  useSearch: () => ({}),
  useParams: () => ({}),
  useLoaderData: () => ({}),
  useRouteContext: () => ({}),
}));

mock.module("@/stores/auth", () => ({
  useAuth: () => ({
    register: () => Promise.resolve(),
  }),
}));

mock.module("@base-ui/react/field", () => {
  const Root = ({ children, ...props }: any) => <div {...props}>{children}</div>;
  const Label = ({ children, ...props }: any) => <label {...props}>{children}</label>;
  const Control = (props: any) => <input {...props} />;
  return { Field: { Root, Label, Control } };
});

mock.module("@/components/button", () => ({
  Button: ({ children, ...props }: any) => <button {...props}>{children}</button>,
}));

mock.module("@/lib/form-error", () => ({
  handleApiError: () => Promise.resolve("Error occurred"),
}));

const kyInstance = {
  get: () => ({ json: () => Promise.resolve({}) }),
  post: () => ({ json: () => Promise.resolve({}) }),
  create: () => kyInstance,
};

mock.module("ky", () => ({
  default: kyInstance,
  HTTPError: class extends Error {
    response: any;
    constructor(r: any) { super("HTTP Error"); this.response = r; }
  },
}));

mock.module("@/lib/api", () => ({
  api: kyInstance,
  HTTPError: class extends Error {
    response: any;
    constructor(r: any) { super("HTTP Error"); this.response = r; }
  },
}));

import { render, screen } from "@testing-library/react";

await import("@/routes/_auth/register");

afterEach(cleanup);

describe("RegisterPage", () => {
  it("renders register form with email, password, and confirm password fields", () => {
    render(<CapturedComponent />);

    expect(screen.getByText("common.email")).toBeDefined();
    expect(screen.getByText("auth.password")).toBeDefined();
    expect(screen.getByText("auth.confirmPassword")).toBeDefined();
  });

  it("renders register button", () => {
    render(<CapturedComponent />);

    expect(screen.getByRole("button", { name: "auth.register" })).toBeDefined();
  });

  it("renders link to login page", () => {
    render(<CapturedComponent />);

    const link = screen.getByText("auth.login");
    expect(link).toBeDefined();
    expect(link.closest("a")?.getAttribute("to")).toBe("/login");
  });
});
