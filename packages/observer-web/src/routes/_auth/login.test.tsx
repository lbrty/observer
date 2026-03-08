import { afterEach, describe, expect, it, mock } from "bun:test";
import { cleanup } from "@testing-library/react";

let CapturedComponent: any;

class MockHTTPError extends Error {
  response: any;
  constructor(r: any) {
    super("HTTP Error");
    this.response = r;
  }
}

const kyInstance = {
  get: () => ({ json: () => Promise.resolve({}) }),
  post: () => ({ json: () => Promise.resolve({}) }),
  patch: () => ({ json: () => Promise.resolve({}) }),
  put: () => ({ json: () => Promise.resolve({}) }),
  delete: () => Promise.resolve(),
  create: () => kyInstance,
};

mock.module("ky", () => ({
  default: kyInstance,
  HTTPError: MockHTTPError,
}));

mock.module("@/lib/api", () => ({
  api: kyInstance,
  HTTPError: MockHTTPError,
}));

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
    login: () => Promise.resolve({ requires_mfa: false }),
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

import { render, screen } from "@testing-library/react";

await import("@/routes/_auth/login");

afterEach(cleanup);

describe("LoginPage", () => {
  it("renders login form with email and password fields", () => {
    render(<CapturedComponent />);

    expect(screen.getByText("common.email")).toBeDefined();
    expect(screen.getByText("auth.password")).toBeDefined();
  });

  it("renders login button", () => {
    render(<CapturedComponent />);

    expect(screen.getByRole("button", { name: "auth.login" })).toBeDefined();
  });

  it("renders link to register page", () => {
    render(<CapturedComponent />);

    const link = screen.getByText("auth.register");
    expect(link).toBeDefined();
    expect(link.getAttribute("to")).toBe("/register");
  });
});
