# Component Extraction Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Split every frontend component file over 150 lines into a folder with focused sibling files; extract sub-components from large route files into `components/`.

**Architecture:** Approach B — sibling files. Each oversized component becomes a folder; `index.tsx` orchestrates and re-exports so existing import sites are unchanged. Route files stay intact; extracted pieces land in `components/` (or `components/reports/`). Drawer indexes get logic hooks to shed lines.

**Tech Stack:** React 19, TypeScript, Tailwind, `@base-ui/react`, TanStack Router/Query, D3, ky

---

## Phase 1 — Shared chart primitives

### Task 1: Shared `ChartTooltip` and `useChartTooltip`

**Files:**
- Create: `packages/observer-web/src/components/charts/chart-tooltip.tsx`
- Create: `packages/observer-web/src/components/charts/use-chart-tooltip.ts`

- [ ] **Step 1: Create `chart-tooltip.tsx`**

```tsx
// packages/observer-web/src/components/charts/chart-tooltip.tsx
interface ChartTooltipProps {
  visible: boolean;
  x: number;
  y: number;
  children: React.ReactNode;
}

export function ChartTooltip({ visible, x, y, children }: ChartTooltipProps) {
  if (!visible) return null;
  return (
    <div
      style={{
        position: "absolute",
        left: x,
        top: y,
        background: "var(--bg-secondary)",
        border: "1px solid var(--border-secondary)",
        borderRadius: 8,
        padding: "6px 10px",
        fontSize: 12,
        color: "var(--fg)",
        boxShadow: "var(--shadow-elevated)",
        pointerEvents: "none",
        whiteSpace: "nowrap",
      }}
    >
      {children}
    </div>
  );
}
```

- [ ] **Step 2: Create `use-chart-tooltip.ts`**

```ts
// packages/observer-web/src/components/charts/use-chart-tooltip.ts
import { type RefObject, useState } from "react";

interface TooltipState {
  visible: boolean;
  x: number;
  y: number;
}

export function useChartTooltip<T extends string>(containerRef: RefObject<HTMLDivElement | null>) {
  const [tooltip, setTooltip] = useState<TooltipState & { text: string }>({
    visible: false,
    x: 0,
    y: 0,
    text: "",
  });

  function show(event: MouseEvent, text: string) {
    const bounds = containerRef.current?.getBoundingClientRect();
    if (!bounds) return;
    setTooltip({ visible: true, x: event.clientX - bounds.left + 12, y: event.clientY - bounds.top - 10, text });
  }

  function hide() {
    setTooltip((prev) => ({ ...prev, visible: false }));
  }

  return { tooltip, show, hide };
}
```

- [ ] **Step 3: Commit**

```bash
git add packages/observer-web/src/components/charts/chart-tooltip.tsx \
        packages/observer-web/src/components/charts/use-chart-tooltip.ts
git commit -m "Add shared ChartTooltip component and useChartTooltip hook"
```

---

## Phase 2 — Component folder conversions

### Task 2: `search-palette.tsx` → `search-palette/`

**Files:**
- Create: `src/components/search-palette/project-group-section.tsx`
- Create: `src/components/search-palette/index.tsx`  
- Delete: `src/components/search-palette.tsx`

- [ ] **Step 1: Create `project-group-section.tsx`**

Extract the `ProjectGroupSection` function (currently lines 135–204 of `search-palette.tsx`):

```tsx
// packages/observer-web/src/components/search-palette/project-group-section.tsx
import { Command } from "cmdk";

import { FolderSimpleIcon, PawPrintIcon, UserFocusIcon } from "@/components/icons";
import type { ProjectGroup } from "@/types/search";

interface ProjectGroupSectionProps {
  group: ProjectGroup;
  onPerson: (projectId: string, personId: string) => void;
  onPet: (projectId: string, petId: string) => void;
  onProject: (projectId: string) => void;
  t: (key: string) => string;
}

export function ProjectGroupSection({ group, onPerson, onPet, onProject, t }: ProjectGroupSectionProps) {
  return (
    <Command.Group
      heading={group.project_name}
      className="[&_[cmdk-group-heading]]:px-2 [&_[cmdk-group-heading]]:py-1 [&_[cmdk-group-heading]]:text-[11px] [&_[cmdk-group-heading]]:font-semibold [&_[cmdk-group-heading]]:uppercase [&_[cmdk-group-heading]]:tracking-wide [&_[cmdk-group-heading]]:text-fg-tertiary"
    >
      {group.people.length > 0 && (
        <Command.Group
          heading={t("search.people")}
          className="[&_[cmdk-group-heading]]:pl-4 [&_[cmdk-group-heading]]:text-[10px] [&_[cmdk-group-heading]]:text-fg-tertiary"
        >
          {group.people.map((p) => (
            <Command.Item
              key={p.id}
              value={`person-${p.id}-${p.first_name} ${p.last_name}`}
              onSelect={() => onPerson(group.project_id, p.id)}
              className="flex cursor-pointer items-center gap-2 rounded-lg px-3 py-1.5 text-sm text-fg outline-none aria-selected:bg-bg-tertiary"
            >
              <UserFocusIcon size={14} className="shrink-0 text-fg-tertiary" />
              {p.first_name} {p.last_name}
            </Command.Item>
          ))}
        </Command.Group>
      )}

      {group.pets.length > 0 && (
        <Command.Group
          heading={t("search.pets")}
          className="[&_[cmdk-group-heading]]:pl-4 [&_[cmdk-group-heading]]:text-[10px] [&_[cmdk-group-heading]]:text-fg-tertiary"
        >
          {group.pets.map((pet) => (
            <Command.Item
              key={pet.id}
              value={`pet-${pet.id}-${pet.name}`}
              onSelect={() => onPet(group.project_id, pet.id)}
              className="flex cursor-pointer items-center gap-2 rounded-lg px-3 py-1.5 text-sm text-fg outline-none aria-selected:bg-bg-tertiary"
            >
              <PawPrintIcon size={14} className="shrink-0 text-fg-tertiary" />
              {pet.name}
            </Command.Item>
          ))}
        </Command.Group>
      )}

      {group.projects.length > 0 && (
        <Command.Group
          heading={t("search.projects")}
          className="[&_[cmdk-group-heading]]:pl-4 [&_[cmdk-group-heading]]:text-[10px] [&_[cmdk-group-heading]]:text-fg-tertiary"
        >
          {group.projects.map((proj) => (
            <Command.Item
              key={proj.id}
              value={`project-${proj.id}-${proj.name}`}
              onSelect={() => onProject(proj.id)}
              className="flex cursor-pointer items-center gap-2 rounded-lg px-3 py-1.5 text-sm text-fg outline-none aria-selected:bg-bg-tertiary"
            >
              <FolderSimpleIcon size={14} className="shrink-0 text-fg-tertiary" />
              {proj.name}
            </Command.Item>
          ))}
        </Command.Group>
      )}
    </Command.Group>
  );
}
```

- [ ] **Step 2: Create `search-palette/index.tsx`** (SearchPalette without the ProjectGroupSection definition)

Copy `search-palette.tsx` to `search-palette/index.tsx`, then:
- Remove the `ProjectGroupSection` function and its interface (lines 135–204)
- Add import: `import { ProjectGroupSection } from "./project-group-section";`

- [ ] **Step 3: Delete the old file and verify**

```bash
rm packages/observer-web/src/components/search-palette.tsx
cd packages/observer-web && bun run build 2>&1 | tail -5
```

Expected: no errors about `search-palette`.

- [ ] **Step 4: Commit**

```bash
git add packages/observer-web/src/components/search-palette/
git rm packages/observer-web/src/components/search-palette.tsx
git commit -m "Split search-palette into folder with ProjectGroupSection sibling"
```

---

### Task 3: `charts/bar-chart.tsx` → `charts/bar-chart/`

The 308-line file has two large render branches (horizontal + vertical) inside a single `useEffect`. Extract each branch into its own module.

**Files:**
- Create: `src/components/charts/bar-chart/render-horizontal.ts`
- Create: `src/components/charts/bar-chart/render-vertical.ts`
- Create: `src/components/charts/bar-chart/chart-legend.tsx`
- Create: `src/components/charts/bar-chart/index.tsx`
- Delete: `src/components/charts/bar-chart.tsx`

- [ ] **Step 1: Create the shared render options type and `render-horizontal.ts`**

```ts
// packages/observer-web/src/components/charts/bar-chart/render-horizontal.ts
import * as d3 from "d3";

import type { CountResult } from "@/types/report";

import { getColor } from "../colors";

export interface RenderOptions {
  width: number;
  height: number;
  resolvedHeight: number;
  yAxisLabel: string;
  selectedLabel: string | null;
  colorMap?: Record<string, string>;
  clipId: string;
  onMouseOver: (event: MouseEvent, d: CountResult) => void;
  onMouseOut: () => void;
  onClick: (event: MouseEvent, d: CountResult) => void;
}

export function renderHorizontalBars(
  svg: d3.Selection<SVGSVGElement, unknown, null, undefined>,
  data: CountResult[],
  opts: RenderOptions,
) {
  const { width, resolvedHeight, selectedLabel, colorMap, clipId, onMouseOver, onMouseOut, onClick } = opts;
  const margin = { top: 10, right: 40, bottom: 20, left: 120 };
  const w = width - margin.left - margin.right;
  const h = resolvedHeight - margin.top - margin.bottom;
  const axisColor = "var(--fg-tertiary, #6b7280)";

  const g = svg
    .attr("viewBox", `0 0 ${width} ${resolvedHeight}`)
    .append("g")
    .attr("transform", `translate(${margin.left},${margin.top})`);

  const yBand = d3.scaleBand<string>().domain(data.map((d) => d.label)).range([0, h]).padding(0.3);
  const xLinear = d3.scaleLinear().domain([0, d3.max(data, (d) => d.count) ?? 0]).nice().range([0, w]);

  const xAxis = g.append("g").attr("transform", `translate(0,${h})`).call(d3.axisBottom(xLinear).ticks(5));
  xAxis.selectAll("text").style("font-size", "9px").style("fill", axisColor);
  xAxis.selectAll("line").style("stroke", axisColor);
  xAxis.select(".domain").style("stroke", axisColor);

  const yAxis = g.append("g").call(d3.axisLeft(yBand));
  yAxis.selectAll("text").style("font-size", "9px").style("fill", axisColor);
  yAxis.selectAll("line").style("stroke", axisColor);
  yAxis.select(".domain").style("stroke", axisColor);

  g.append("clipPath").attr("id", clipId).append("rect").attr("width", w).attr("height", h);

  const opacityFn = (d: CountResult) => (selectedLabel === null || selectedLabel === d.label ? 1 : 0.3);
  const rH = 3;

  g.append("g")
    .attr("clip-path", `url(#${clipId})`)
    .selectAll(".bar")
    .data(data)
    .join("path")
    .attr("class", "bar")
    .attr("d", (d) => {
      const bx = 0;
      const by = yBand(d.label) ?? 0;
      const bw = xLinear(d.count);
      const bh = yBand.bandwidth();
      return `M${bx},${by} H${bx + bw - rH} a${rH},${rH} 0 0 1 ${rH},${rH} V${by + bh - rH} a${rH},${rH} 0 0 1 ${-rH},${rH} H${bx} Z`;
    })
    .attr("fill", (d, i) => getColor(d.label, colorMap, i))
    .attr("opacity", opacityFn)
    .style("cursor", "pointer")
    .on("mouseover", onMouseOver)
    .on("mousemove", onMouseOver)
    .on("mouseout", onMouseOut)
    .on("click", onClick);

  g.selectAll(".label")
    .data(data)
    .join("text")
    .attr("class", "label")
    .attr("x", (d) => xLinear(d.count) + 6)
    .attr("y", (d) => (yBand(d.label) ?? 0) + yBand.bandwidth() / 2)
    .attr("dominant-baseline", "central")
    .attr("text-anchor", "start")
    .style("font-size", "9px")
    .style("fill", "var(--fg-secondary, #6b7280)")
    .attr("opacity", opacityFn)
    .text((d) => d.count);
}
```

- [ ] **Step 2: Create `render-vertical.ts`**

```ts
// packages/observer-web/src/components/charts/bar-chart/render-vertical.ts
import * as d3 from "d3";

import type { CountResult } from "@/types/report";

import { getColor } from "../colors";
import type { RenderOptions } from "./render-horizontal";

export function renderVerticalBars(
  svg: d3.Selection<SVGSVGElement, unknown, null, undefined>,
  data: CountResult[],
  opts: RenderOptions,
) {
  const { width, height, selectedLabel, colorMap, clipId, yAxisLabel, onMouseOver, onMouseOut, onClick } = opts;
  const margin = { top: 20, right: 20, bottom: 60, left: 50 };
  const w = width - margin.left - margin.right;
  const h = height - margin.top - margin.bottom;
  const axisColor = "var(--fg-tertiary, #6b7280)";

  const g = svg
    .attr("viewBox", `0 0 ${width} ${height}`)
    .append("g")
    .attr("transform", `translate(${margin.left},${margin.top})`);

  const x = d3.scaleBand<string>().domain(data.map((d) => d.label)).range([0, w]).padding(0.3);
  const y = d3.scaleLinear().domain([0, d3.max(data, (d) => d.count) ?? 0]).nice().range([h, 0]);

  const xAxis = g.append("g").attr("transform", `translate(0,${h})`).call(d3.axisBottom(x));
  xAxis.selectAll("text").attr("transform", "rotate(-35)").style("text-anchor", "end").style("font-size", "9px").style("fill", axisColor);
  xAxis.selectAll("line").style("stroke", axisColor);
  xAxis.select(".domain").style("stroke", axisColor);

  const yAxis = g.append("g").call(d3.axisLeft(y).ticks(5));
  yAxis.selectAll("text").style("font-size", "9px").style("fill", axisColor);
  yAxis.selectAll("line").style("stroke", axisColor);
  yAxis.select(".domain").style("stroke", axisColor);

  g.append("text")
    .attr("transform", "rotate(-90)")
    .attr("x", -h / 2)
    .attr("y", -margin.left + 14)
    .attr("text-anchor", "middle")
    .style("font-size", "10px")
    .style("fill", "var(--fg-tertiary, #6b7280)")
    .text(yAxisLabel);

  g.append("clipPath").attr("id", clipId).append("rect").attr("width", w).attr("height", h);

  const opacityFn = (d: CountResult) => (selectedLabel === null || selectedLabel === d.label ? 1 : 0.3);
  const rV = 3;

  g.append("g")
    .attr("clip-path", `url(#${clipId})`)
    .selectAll(".bar")
    .data(data)
    .join("path")
    .attr("class", "bar")
    .attr("d", (d) => {
      const bx = x(d.label) ?? 0;
      const by = y(d.count);
      const bw = x.bandwidth();
      const bh = h - y(d.count);
      return `M${bx},${by + rV} a${rV},${rV} 0 0 1 ${rV},${-rV} H${bx + bw - rV} a${rV},${rV} 0 0 1 ${rV},${rV} V${by + bh} H${bx} Z`;
    })
    .attr("fill", (d, i) => getColor(d.label, colorMap, i))
    .attr("opacity", opacityFn)
    .style("cursor", "pointer")
    .on("mouseover", onMouseOver)
    .on("mousemove", onMouseOver)
    .on("mouseout", onMouseOut)
    .on("click", onClick);

  g.selectAll(".label")
    .data(data)
    .join("text")
    .attr("class", "label")
    .attr("x", (d) => (x(d.label) ?? 0) + x.bandwidth() / 2)
    .attr("y", (d) => y(d.count) - 4)
    .attr("text-anchor", "middle")
    .style("font-size", "9px")
    .style("fill", "var(--fg-secondary, #6b7280)")
    .attr("opacity", opacityFn)
    .text((d) => d.count);
}
```

- [ ] **Step 3: Create `chart-legend.tsx`**

```tsx
// packages/observer-web/src/components/charts/bar-chart/chart-legend.tsx
export interface BarLegendItem {
  short: string;
  full: string;
}

interface ChartLegendProps {
  items: BarLegendItem[];
}

export function ChartLegend({ items }: ChartLegendProps) {
  if (items.length === 0) return null;
  return (
    <div className="mt-3 flex flex-wrap gap-x-4 gap-y-1 border-t border-border-secondary pt-3">
      {items.map((item) => (
        <span key={item.short} className="text-[11px] text-fg-tertiary">
          <span className="font-medium text-fg-secondary">{item.short}</span> — {item.full}
        </span>
      ))}
    </div>
  );
}
```

- [ ] **Step 4: Create `bar-chart/index.tsx`**

```tsx
// packages/observer-web/src/components/charts/bar-chart/index.tsx
import { useEffect, useId, useRef, useState } from "react";

import * as d3 from "d3";

import type { CountResult } from "@/types/report";

import { ChartTooltip } from "../chart-tooltip";
import { ChartLegend, type BarLegendItem } from "./chart-legend";
import { renderHorizontalBars } from "./render-horizontal";
import { renderVerticalBars } from "./render-vertical";

export type { BarLegendItem };

interface BarChartProps {
  data: CountResult[];
  width?: number;
  height?: number;
  yAxisLabel?: string;
  legend?: BarLegendItem[];
  direction?: "vertical" | "horizontal" | "auto";
  colorMap?: Record<string, string>;
}

export function BarChart({
  data,
  width = 700,
  height = 260,
  yAxisLabel = "Count",
  legend,
  direction,
  colorMap,
}: BarChartProps) {
  const ref = useRef<SVGSVGElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const clipId = useId().replace(/:/g, "_");
  const [tooltip, setTooltip] = useState({ visible: false, x: 0, y: 0, label: "", count: 0 });
  const [selectedLabel, setSelectedLabel] = useState<string | null>(null);

  const resolvedDirection =
    direction === "auto" ? (data.length > 6 ? "horizontal" : "vertical") : (direction ?? "vertical");
  const resolvedHeight =
    resolvedDirection === "horizontal" ? Math.max(200, data.length * 24) : height;

  useEffect(() => {
    if (!ref.current || data.length === 0) return;
    const svg = d3.select(ref.current);
    svg.selectAll("*").remove();

    function handleMouseOver(event: MouseEvent, d: CountResult) {
      const bounds = containerRef.current?.getBoundingClientRect();
      if (!bounds) return;
      setTooltip({ visible: true, x: event.clientX - bounds.left + 12, y: event.clientY - bounds.top - 10, label: d.label, count: d.count });
    }
    function handleMouseOut() { setTooltip((prev) => ({ ...prev, visible: false })); }
    function handleClick(_: MouseEvent, d: CountResult) { setSelectedLabel((prev) => (prev === d.label ? null : d.label)); }

    const opts = { width, height, resolvedHeight, yAxisLabel, selectedLabel, colorMap, clipId, onMouseOver: handleMouseOver, onMouseOut: handleMouseOut, onClick: handleClick };

    if (resolvedDirection === "horizontal") {
      renderHorizontalBars(svg, data, opts);
    } else {
      renderVerticalBars(svg, data, opts);
    }
  }, [data, width, height, resolvedHeight, yAxisLabel, selectedLabel, resolvedDirection, colorMap, clipId]);

  if (data.length === 0) return null;

  return (
    <div ref={containerRef} className="relative w-full">
      <svg ref={ref} className="w-full" style={{ maxHeight: 360 }} />
      <ChartTooltip visible={tooltip.visible} x={tooltip.x} y={tooltip.y}>
        <strong>{tooltip.label}</strong>: {tooltip.count}
      </ChartTooltip>
      {legend && <ChartLegend items={legend} />}
    </div>
  );
}
```

- [ ] **Step 5: Delete old file and verify**

```bash
rm packages/observer-web/src/components/charts/bar-chart.tsx
cd packages/observer-web && bun run build 2>&1 | grep -E "error|Error" | head -10
```

- [ ] **Step 6: Commit**

```bash
git add packages/observer-web/src/components/charts/bar-chart/
git rm packages/observer-web/src/components/charts/bar-chart.tsx
git commit -m "Split bar-chart into folder with render-horizontal, render-vertical, chart-legend"
```

---

### Task 4: `charts/sankey-chart.tsx` → `charts/sankey-chart/`

**Files:**
- Create: `src/components/charts/sankey-chart/index.tsx`
- Delete: `src/components/charts/sankey-chart.tsx`

The sankey chart is 202 lines. Extract tooltip state via `useChartTooltip` and replace the tooltip div with `<ChartTooltip>`.

- [ ] **Step 1: Create `sankey-chart/index.tsx`**

Read `charts/sankey-chart.tsx` then rewrite as `sankey-chart/index.tsx` with these changes:
- Add imports: `import { ChartTooltip } from "../chart-tooltip"; import { useChartTooltip } from "../use-chart-tooltip";`
- Replace the `useState<Tooltip>` block with `const { tooltip, show, hide } = useChartTooltip(containerRef);`
- In the `useEffect` link mouseover handler, call `show(event, \`${tl(...)  → ... }\`)` instead of `setTooltip({...})`
- In the link mouseout handler, call `hide()` instead of `setTooltip(prev => ...)`
- In the return JSX replace the tooltip `<div style={...}>` with `<ChartTooltip visible={tooltip.visible} x={tooltip.x} y={tooltip.y}>{tooltip.text}</ChartTooltip>`

The resulting `index.tsx` will be ~170 lines (within the soft ceiling).

- [ ] **Step 2: Delete old file and verify**

```bash
rm packages/observer-web/src/components/charts/sankey-chart.tsx
cd packages/observer-web && bun run build 2>&1 | grep -E "error|Error" | head -10
```

- [ ] **Step 3: Commit**

```bash
git add packages/observer-web/src/components/charts/sankey-chart/
git rm packages/observer-web/src/components/charts/sankey-chart.tsx
git commit -m "Split sankey-chart into folder, use shared ChartTooltip"
```

---

### Task 5: `date-picker.tsx` → `date-picker/`

**Files:**
- Create: `src/components/date-picker/utils.ts`
- Create: `src/components/date-picker/date-picker.tsx`
- Create: `src/components/date-picker/date-range-picker.tsx`
- Create: `src/components/date-picker/index.tsx`
- Move: `src/components/date-picker.css` → `src/components/date-picker/date-picker.css`
- Delete: `src/components/date-picker.tsx`

- [ ] **Step 1: Create `utils.ts`**

```ts
// packages/observer-web/src/components/date-picker/utils.ts
export const triggerClass =
  "flex h-9 w-full items-center gap-2 rounded-lg border border-border-secondary bg-bg-secondary px-3 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg disabled:opacity-50";

export function formatDisplay(iso: string): string {
  if (!iso) return "";
  const [y, m, d] = iso.split("-");
  return `${d}.${m}.${y}`;
}

export function toISO(date: Date): string {
  const y = date.getFullYear();
  const m = String(date.getMonth() + 1).padStart(2, "0");
  const d = String(date.getDate()).padStart(2, "0");
  return `${y}-${m}-${d}`;
}

export function parseISO(iso: string): Date | undefined {
  if (!iso) return undefined;
  const [y, m, d] = iso.split("-").map(Number);
  return new Date(y, m - 1, d);
}
```

- [ ] **Step 2: Create `date-picker/date-picker.tsx`**

Copy the `DatePicker` export from the original file. Replace the inline utility calls with imports from `./utils`. Replace the CSS import with `./date-picker.css`.

- [ ] **Step 3: Create `date-picker/date-range-picker.tsx`**

Copy the `DateRangePicker` export. Import utilities from `./utils` and CSS from `./date-picker.css`.

- [ ] **Step 4: Create `date-picker/index.tsx`**

```tsx
// packages/observer-web/src/components/date-picker/index.tsx
export { DatePicker } from "./date-picker";
export { DateRangePicker } from "./date-range-picker";
```

- [ ] **Step 5: Move CSS, delete old file, verify**

```bash
mv packages/observer-web/src/components/date-picker.css \
   packages/observer-web/src/components/date-picker/date-picker.css
rm packages/observer-web/src/components/date-picker.tsx
cd packages/observer-web && bun run build 2>&1 | grep -E "error|Error" | head -10
```

- [ ] **Step 6: Commit**

```bash
git add packages/observer-web/src/components/date-picker/
git rm packages/observer-web/src/components/date-picker.tsx \
       packages/observer-web/src/components/date-picker.css
git commit -m "Split date-picker into folder with utils, DatePicker, DateRangePicker"
```

---

### Task 6: `profile/mfa-settings.tsx` → `profile/mfa-settings/`

**Files:**
- Create: `src/components/profile/mfa-settings/mfa-active.tsx`
- Create: `src/components/profile/mfa-settings/mfa-setup.tsx`
- Create: `src/components/profile/mfa-settings/index.tsx`
- Delete: `src/components/profile/mfa-settings.tsx`

- [ ] **Step 1: Create `mfa-active.tsx`** (the enabled state — shows status badge + optional disable form)

```tsx
// packages/observer-web/src/components/profile/mfa-settings/mfa-active.tsx
import type { SyntheticEvent } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { FormField } from "@/components/form-field";

interface MFAActiveProps {
  showDisable: boolean;
  disableCode: string;
  isPending: boolean;
  onDisableCode: (v: string) => void;
  onShowDisable: () => void;
  onCancelDisable: () => void;
  onSubmit: (e: SyntheticEvent) => void;
}

export function MFAActive({
  showDisable,
  disableCode,
  isPending,
  onDisableCode,
  onShowDisable,
  onCancelDisable,
  onSubmit,
}: MFAActiveProps) {
  const { t } = useTranslation();
  return (
    <div className="rounded-lg border border-border-secondary p-4">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm font-medium text-fg">{t("profile.mfaActive")}</p>
          <p className="text-xs text-fg-secondary mt-0.5">{t("profile.mfaActiveHint")}</p>
        </div>
        <span className="inline-flex items-center gap-1 rounded-full bg-green-100 px-2.5 py-0.5 text-xs font-medium text-green-800 dark:bg-green-900/30 dark:text-green-300">
          {t("common.enabled")}
        </span>
      </div>
      {!showDisable ? (
        <Button variant="danger" className="mt-3" onClick={onShowDisable}>
          {t("profile.disableMFA")}
        </Button>
      ) : (
        <form onSubmit={onSubmit} className="mt-3 space-y-3">
          <p className="text-sm text-fg-secondary">{t("profile.disableMFAHint")}</p>
          <FormField label={t("auth.totpCode")} value={disableCode} onChange={onDisableCode} maxLength={6} />
          <div className="flex gap-2">
            <Button type="submit" variant="danger" disabled={isPending}>
              {t("profile.confirmDisable")}
            </Button>
            <Button type="button" variant="secondary" onClick={onCancelDisable}>
              {t("common.cancel")}
            </Button>
          </div>
        </form>
      )}
    </div>
  );
}
```

- [ ] **Step 2: Create `mfa-setup.tsx`** (the disabled state — QR code + TOTP form)

```tsx
// packages/observer-web/src/components/profile/mfa-settings/mfa-setup.tsx
import type { SyntheticEvent } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { FormField } from "@/components/form-field";

interface MFASetupProps {
  showSetup: boolean;
  qrDataURL: string;
  secret?: string;
  totpCode: string;
  isLoading: boolean;
  hasData: boolean;
  isPending: boolean;
  onTotpCode: (v: string) => void;
  onShowSetup: () => void;
  onCancelSetup: () => void;
  onSubmit: (e: SyntheticEvent) => void;
}

export function MFASetup({
  showSetup,
  qrDataURL,
  secret,
  totpCode,
  isLoading,
  hasData,
  isPending,
  onTotpCode,
  onShowSetup,
  onCancelSetup,
  onSubmit,
}: MFASetupProps) {
  const { t } = useTranslation();
  return (
    <div className="rounded-lg border border-border-secondary p-4">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm font-medium text-fg">{t("profile.mfaInactive")}</p>
          <p className="text-xs text-fg-secondary mt-0.5">{t("profile.mfaInactiveHint")}</p>
        </div>
        <span className="inline-flex items-center gap-1 rounded-full bg-bg-secondary px-2.5 py-0.5 text-xs font-medium text-fg-tertiary">
          {t("common.disabled")}
        </span>
      </div>
      {!showSetup ? (
        <Button className="mt-3" onClick={onShowSetup}>
          {t("profile.setupMFA")}
        </Button>
      ) : (
        <div className="mt-4 space-y-4">
          <p className="text-sm text-fg-secondary">{t("profile.scanQR")}</p>
          {isLoading && <p className="text-sm text-fg-tertiary">{t("common.loading")}</p>}
          {qrDataURL && (
            <div className="flex flex-col items-start gap-3">
              <img src={qrDataURL} alt={t("profile.totpQrCodeAlt")} className="rounded-lg border border-border-secondary" width={180} height={180} />
              <details className="text-xs text-fg-tertiary">
                <summary className="cursor-pointer">{t("profile.cantScanQR")}</summary>
                <code className="mt-1 block break-all font-mono text-xs text-fg-secondary select-all">{secret}</code>
              </details>
            </div>
          )}
          {hasData && (
            <form onSubmit={onSubmit} className="space-y-3">
              <FormField label={t("profile.enterCode")} value={totpCode} onChange={onTotpCode} maxLength={6} required />
              <div className="flex gap-2">
                <Button type="submit" disabled={isPending || totpCode.length !== 6}>
                  {t("profile.verifyAndEnable")}
                </Button>
                <Button type="button" variant="secondary" onClick={onCancelSetup}>
                  {t("common.cancel")}
                </Button>
              </div>
            </form>
          )}
        </div>
      )}
    </div>
  );
}
```

- [ ] **Step 3: Create `mfa-settings/index.tsx`**

```tsx
// packages/observer-web/src/components/profile/mfa-settings/index.tsx
import { type SyntheticEvent, useEffect, useState } from "react";

import QRCode from "qrcode";
import { useTranslation } from "react-i18next";

import { ErrorBanner } from "@/components/alert-banner";
import { useDisableMFA, useEnableMFA, useMFASetup } from "@/hooks/use-mfa";
import { handleApiError } from "@/lib/form-error";
import { useAuth } from "@/stores/auth";
import { useToast } from "@/stores/toast";

import { MFAActive } from "./mfa-active";
import { MFASetup } from "./mfa-setup";

export function MFASettings() {
  const { t } = useTranslation();
  const { user, setUser } = useAuth();
  const toast = useToast();

  const [showSetup, setShowSetup] = useState(false);
  const [qrDataURL, setQrDataURL] = useState("");
  const [totpCode, setTotpCode] = useState("");
  const [disableCode, setDisableCode] = useState("");
  const [showDisable, setShowDisable] = useState(false);
  const [error, setError] = useState("");

  const setup = useMFASetup(showSetup);
  const enableMFA = useEnableMFA();
  const disableMFA = useDisableMFA();

  useEffect(() => {
    if (setup.data?.otpauth_url) {
      QRCode.toDataURL(setup.data.otpauth_url, { width: 200 }).then(setQrDataURL);
    }
  }, [setup.data?.otpauth_url]);

  async function handleEnable(e: SyntheticEvent) {
    e.preventDefault();
    setError("");
    if (!setup.data) return;
    try {
      await enableMFA.mutateAsync({ secret: setup.data.secret, totpCode });
      if (user) setUser({ ...user, mfa_enabled: true });
      toast.success(t("profile.mfaEnabled"));
      setShowSetup(false);
      setTotpCode("");
      setQrDataURL("");
    } catch (err) {
      setError(await handleApiError(err, t));
    }
  }

  async function handleDisable(e: SyntheticEvent) {
    e.preventDefault();
    setError("");
    try {
      await disableMFA.mutateAsync(disableCode);
      if (user) setUser({ ...user, mfa_enabled: false });
      toast.success(t("profile.mfaDisabled"));
      setShowDisable(false);
      setDisableCode("");
    } catch (err) {
      setError(await handleApiError(err, t));
    }
  }

  return (
    <div className="space-y-4">
      <h2 className="text-sm font-semibold text-fg">{t("profile.twoFactor")}</h2>
      <ErrorBanner message={error} />
      {user?.mfa_enabled ? (
        <MFAActive
          showDisable={showDisable}
          disableCode={disableCode}
          isPending={disableMFA.isPending}
          onDisableCode={setDisableCode}
          onShowDisable={() => setShowDisable(true)}
          onCancelDisable={() => setShowDisable(false)}
          onSubmit={handleDisable}
        />
      ) : (
        <MFASetup
          showSetup={showSetup}
          qrDataURL={qrDataURL}
          secret={setup.data?.secret}
          totpCode={totpCode}
          isLoading={setup.isLoading}
          hasData={!!setup.data}
          isPending={enableMFA.isPending}
          onTotpCode={setTotpCode}
          onShowSetup={() => setShowSetup(true)}
          onCancelSetup={() => { setShowSetup(false); setTotpCode(""); setQrDataURL(""); }}
          onSubmit={handleEnable}
        />
      )}
    </div>
  );
}
```

- [ ] **Step 4: Delete old file and verify**

```bash
rm packages/observer-web/src/components/profile/mfa-settings.tsx
cd packages/observer-web && bun run build 2>&1 | grep -E "error|Error" | head -10
```

- [ ] **Step 5: Commit**

```bash
git add packages/observer-web/src/components/profile/mfa-settings/
git rm packages/observer-web/src/components/profile/mfa-settings.tsx
git commit -m "Split mfa-settings into folder with MFAActive and MFASetup"
```

---

### Task 7: `permissions/assign-dialog.tsx` → `permissions/assign-dialog/`

**Files:**
- Create: `src/components/permissions/assign-dialog/permission-toggle-row.tsx`
- Create: `src/components/permissions/assign-dialog/selected-user-card.tsx`
- Create: `src/components/permissions/assign-dialog/index.tsx`
- Delete: `src/components/permissions/assign-dialog.tsx`

- [ ] **Step 1: Create `permission-toggle-row.tsx`**

```tsx
// packages/observer-web/src/components/permissions/assign-dialog/permission-toggle-row.tsx
import { UISwitch } from "@/components/ui-switch";

interface PermissionToggleRowProps {
  label: string;
  description: string;
  checked: boolean;
  onCheckedChange: (v: boolean) => void;
}

export function PermissionToggleRow({ label, description, checked, onCheckedChange }: PermissionToggleRowProps) {
  return (
    <div className="space-y-2">
      <UISwitch checked={checked} onCheckedChange={onCheckedChange} label={label} />
      <p className="ml-11.5 text-xs text-fg-tertiary">{description}</p>
    </div>
  );
}
```

- [ ] **Step 2: Create `selected-user-card.tsx`**

```tsx
// packages/observer-web/src/components/permissions/assign-dialog/selected-user-card.tsx
import { XIcon } from "@/components/icons";
import type { AdminUser } from "@/types/admin";

interface SelectedUserCardProps {
  user: AdminUser;
  onClear: () => void;
}

export function SelectedUserCard({ user, onClear }: SelectedUserCardProps) {
  return (
    <div className="flex items-center justify-between rounded-lg border border-border-secondary bg-bg px-3 py-2">
      <div>
        <p className="text-sm text-fg">{user.first_name} {user.last_name}</p>
        <p className="text-xs text-fg-tertiary">{user.email}</p>
      </div>
      <button type="button" onClick={onClear} className="cursor-pointer rounded p-1 text-fg-tertiary hover:text-fg">
        <XIcon size={14} />
      </button>
    </div>
  );
}
```

- [ ] **Step 3: Create `assign-dialog/index.tsx`**

Copy `assign-dialog.tsx` into the folder as `index.tsx`. Replace the four `UISwitch + <p>` blocks with `<PermissionToggleRow>`. Replace the selected-user block with `<SelectedUserCard>`. Add imports for the two new components. Remove `XIcon` import if no longer used directly.

The resulting file should be ~125 lines.

- [ ] **Step 4: Delete old file, verify, commit**

```bash
rm packages/observer-web/src/components/permissions/assign-dialog.tsx
cd packages/observer-web && bun run build 2>&1 | grep -E "error|Error" | head -10
git add packages/observer-web/src/components/permissions/assign-dialog/
git rm packages/observer-web/src/components/permissions/assign-dialog.tsx
git commit -m "Split assign-dialog into folder with PermissionToggleRow and SelectedUserCard"
```

---

### Task 8: Drawer index hooks — `migration-record-drawer`

At 252 lines, the index needs a custom hook to shed ~100 lines of form/submit logic.

**Files:**
- Create: `src/components/migration-record-drawer/use-migration-record-form.ts`
- Modify: `src/components/migration-record-drawer/index.tsx`

- [ ] **Step 1: Create `use-migration-record-form.ts`**

Extract from `index.tsx`:
- `emptyForm` constant
- `useDrawerForm` call + all `useCountries/useStates/usePlaces` calls
- The `useEffect` that resolves place→state→country chains (lines 93–119)
- `handleSubmit` function (lines 123–169)
- The options derivation (lines 171–175)

```ts
// packages/observer-web/src/components/migration-record-drawer/use-migration-record-form.ts
import { type SyntheticEvent, useEffect } from "react";

import { useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";

import { useCountries } from "@/hooks/use-countries";
import { useDrawerForm } from "@/hooks/use-drawer-form";
import {
  useCreateMigrationRecord,
  useMigrationRecord,
  useUpdateMigrationRecord,
} from "@/hooks/use-migration-records";
import { usePlaces } from "@/hooks/use-places";
import { useStates } from "@/hooks/use-states";
import { handleApiError } from "@/lib/form-error";
import { toSelectOptions } from "@/lib/options";
import { useToast } from "@/stores/toast";
import type { HousingAtDestination, MovementReason } from "@/types/migration-record";

const emptyForm = {
  from_country: "",
  from_state: "",
  from_place_id: "",
  dest_country: "",
  dest_state: "",
  destination_place_id: "",
  migration_date: "",
  movement_reason: "",
  housing_at_destination: "",
  notes: "",
};

interface UseMigrationRecordFormOptions {
  open: boolean;
  projectId: string;
  personId: string;
  recordId: string | null;
}

export function useMigrationRecordForm({ open, projectId, personId, recordId }: UseMigrationRecordFormOptions) {
  const { t } = useTranslation();
  const isEdit = recordId !== null;

  const { data: record } = useMigrationRecord(projectId, personId, recordId ?? "");
  const qc = useQueryClient();
  const createRecord = useCreateMigrationRecord(projectId, personId);
  const updateRecord = useUpdateMigrationRecord(projectId, personId);
  const toast = useToast();

  const { form, set, error, setError, editingId, setEditingId } = useDrawerForm({
    initial: emptyForm,
    open,
    isEdit,
    data: record,
    mapData: (d) => ({
      from_country: "",
      from_state: "",
      from_place_id: (d.from_place_id as string) ?? "",
      dest_country: "",
      dest_state: "",
      destination_place_id: (d.destination_place_id as string) ?? "",
      migration_date: (d.migration_date as string) ?? "",
      movement_reason: (d.movement_reason as string) ?? "",
      housing_at_destination: (d.housing_at_destination as string) ?? "",
      notes: (d.notes as string) ?? "",
    }),
  });

  const { data: countries } = useCountries();
  const { data: allStates } = useStates();
  const { data: allPlaces } = usePlaces();
  const { data: fromStates } = useStates(form.from_country || undefined);
  const { data: fromPlaces } = usePlaces(form.from_state || undefined);
  const { data: destStates } = useStates(form.dest_country || undefined);
  const { data: destPlaces } = usePlaces(form.dest_state || undefined);

  useEffect(() => {
    if (!isEdit || !record || !allStates?.states || !allPlaces?.places) return;
    const states = allStates.states;
    const places = allPlaces.places;
    if (record.from_place_id && !form.from_country) {
      const place = places.find((p) => p.id === record.from_place_id);
      if (place) {
        const state = states.find((s) => s.id === place.state_id);
        if (state) { set("from_country", state.country_id); set("from_state", state.id); }
      }
    }
    if (record.destination_place_id && !form.dest_country) {
      const place = places.find((p) => p.id === record.destination_place_id);
      if (place) {
        const state = states.find((s) => s.id === place.state_id);
        if (state) { set("dest_country", state.country_id); set("dest_state", state.id); }
      }
    }
  }, [record, allStates, allPlaces]);

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    setError("");
    try {
      if (isEdit && editingId) {
        await updateRecord.mutateAsync({
          id: editingId,
          data: {
            ...(form.from_place_id && { from_place_id: form.from_place_id }),
            ...(form.destination_place_id && { destination_place_id: form.destination_place_id }),
            ...(form.migration_date && { migration_date: form.migration_date }),
            ...(form.movement_reason && { movement_reason: form.movement_reason as MovementReason }),
            ...(form.housing_at_destination && { housing_at_destination: form.housing_at_destination as HousingAtDestination }),
            ...(form.notes && { notes: form.notes }),
          },
        });
        await qc.invalidateQueries({ queryKey: ["migration-records", projectId, personId] });
        toast.success(t("project.migrationRecords.saved"));
      } else {
        const created = await createRecord.mutateAsync({
          ...(form.from_place_id && { from_place_id: form.from_place_id }),
          ...(form.destination_place_id && { destination_place_id: form.destination_place_id }),
          ...(form.migration_date && { migration_date: form.migration_date }),
          ...(form.movement_reason && { movement_reason: form.movement_reason as MovementReason }),
          ...(form.housing_at_destination && { housing_at_destination: form.housing_at_destination as HousingAtDestination }),
          ...(form.notes && { notes: form.notes }),
        });
        await qc.invalidateQueries({ queryKey: ["migration-records", projectId, personId] });
        setEditingId(created.id);
        toast.success(t("project.migrationRecords.saved"));
      }
    } catch (err) {
      setError(await handleApiError(err, t));
    }
  }

  return {
    isEdit,
    form,
    set,
    error,
    isPending: createRecord.isPending || updateRecord.isPending,
    handleSubmit,
    countryOptions: toSelectOptions(countries),
    fromStateOptions: toSelectOptions(fromStates?.states),
    fromPlaceOptions: toSelectOptions(fromPlaces?.places),
    destStateOptions: toSelectOptions(destStates?.states),
    destPlaceOptions: toSelectOptions(destPlaces?.places),
  };
}
```

- [ ] **Step 2: Rewrite `index.tsx` to use the hook**

```tsx
// packages/observer-web/src/components/migration-record-drawer/index.tsx
import type { SyntheticEvent } from "react";

import { useTranslation } from "react-i18next";

import { ErrorBanner } from "@/components/alert-banner";
import { DrawerShell } from "@/components/drawer-shell";
import { FormTextarea } from "@/components/form-field";
import { SectionHeading } from "@/components/section-heading";

import { DetailsSection } from "./details-section";
import { PlaceSection } from "./place-section";
import { useMigrationRecordForm } from "./use-migration-record-form";

interface MigrationRecordDrawerProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  projectId: string;
  personId: string;
  recordId: string | null;
}

export function MigrationRecordDrawer({ open, onOpenChange, projectId, personId, recordId }: MigrationRecordDrawerProps) {
  const { t } = useTranslation();
  const {
    isEdit,
    form,
    set,
    error,
    isPending,
    handleSubmit,
    countryOptions,
    fromStateOptions,
    fromPlaceOptions,
    destStateOptions,
    destPlaceOptions,
  } = useMigrationRecordForm({ open, projectId, personId, recordId });

  return (
    <DrawerShell
      open={open}
      onOpenChange={onOpenChange}
      title={isEdit ? t("project.migrationRecords.editTitle") : t("project.migrationRecords.addTitle")}
      onSubmit={handleSubmit}
      isPending={isPending}
      submitLabel={t("project.migrationRecords.save")}
      savingLabel={t("project.migrationRecords.saving")}
    >
      {error && <ErrorBanner message={error} />}

      <PlaceSection
        title={t("project.migrationRecords.from")}
        country={form.from_country}
        state={form.from_state}
        place={form.from_place_id}
        countryOptions={countryOptions}
        stateOptions={fromStateOptions}
        placeOptions={fromPlaceOptions}
        countryPlaceholder={t("project.people.selectCountry")}
        statePlaceholder={t("project.people.selectState")}
        placePlaceholder={t("project.people.selectPlace")}
        onCountryChange={(v) => { set("from_country", v); set("from_state", ""); set("from_place_id", ""); }}
        onStateChange={(v) => { set("from_state", v); set("from_place_id", ""); }}
        onPlaceChange={(v) => set("from_place_id", v)}
      />

      <PlaceSection
        title={t("project.migrationRecords.to")}
        country={form.dest_country}
        state={form.dest_state}
        place={form.destination_place_id}
        countryOptions={countryOptions}
        stateOptions={destStateOptions}
        placeOptions={destPlaceOptions}
        countryPlaceholder={t("project.people.selectCountry")}
        statePlaceholder={t("project.people.selectState")}
        placePlaceholder={t("project.people.selectPlace")}
        onCountryChange={(v) => { set("dest_country", v); set("dest_state", ""); set("destination_place_id", ""); }}
        onStateChange={(v) => { set("dest_state", v); set("destination_place_id", ""); }}
        onPlaceChange={(v) => set("destination_place_id", v)}
      />

      <DetailsSection
        migrationDate={form.migration_date}
        movementReason={form.movement_reason}
        housingAtDestination={form.housing_at_destination}
        onDateChange={(v) => set("migration_date", v)}
        onReasonChange={(v) => set("movement_reason", v)}
        onHousingChange={(v) => set("housing_at_destination", v)}
      />

      <SectionHeading>{t("project.migrationRecords.notes")}</SectionHeading>
      <FormTextarea label="" value={form.notes} onChange={(v) => set("notes", v)} rows={4} />
    </DrawerShell>
  );
}
```

- [ ] **Step 3: Verify and commit**

```bash
cd packages/observer-web && bun run build 2>&1 | grep -E "error|Error" | head -10
git add packages/observer-web/src/components/migration-record-drawer/
git commit -m "Extract useMigrationRecordForm hook to slim migration-record-drawer index"
```

---

### Task 9: `person-drawer/index.tsx` — extract `usePersonDrawerForm`

**Files:**
- Create: `src/components/person-drawer/use-person-drawer-form.ts`
- Modify: `src/components/person-drawer/index.tsx`

- [ ] **Step 1: Create `use-person-drawer-form.ts`**

Extract from `index.tsx`:
- `emptyForm` constant
- All hook calls (`useDrawerForm`, `useCountries`, `useStates`, `usePlaces`, `useOffices`, `usePersonTags`, `useReplacePersonTags`, `useCreatePerson`, `useUpdatePerson`)
- `tagIds` state + its two `useEffect`s
- `originPlaceLabel`/`currentPlaceLabel` state
- `resolvePlaceLabel` function
- `handleSubmit` function

The hook returns `{ isEdit, form, set, error, isPending, handleSubmit, tagIds, setTagIds, officeOptions, resolvedOriginLabel, resolvedCurrentLabel, setOriginPlaceLabel, setCurrentPlaceLabel }`.

Read `person-drawer/index.tsx` to copy exact implementations. The resulting hook will be ~120 lines; `index.tsx` will be ~80 lines.

- [ ] **Step 2: Rewrite `index.tsx`**

Import the hook, pass all values to `IdentitySection`, `LocationSection`, `CaseSection`, and `TagPicker`. The index becomes a thin orchestrator with just the JSX.

- [ ] **Step 3: Verify and commit**

```bash
cd packages/observer-web && bun run build 2>&1 | grep -E "error|Error" | head -10
git add packages/observer-web/src/components/person-drawer/
git commit -m "Extract usePersonDrawerForm hook to slim person-drawer index"
```

---

### Task 10: `support-record-drawer/index.tsx` — extract `useSupportRecordForm`

**Files:**
- Create: `src/components/support-record-drawer/use-support-record-form.ts`
- Modify: `src/components/support-record-drawer/index.tsx`

- [ ] **Step 1: Create `use-support-record-form.ts`**

Extract from `index.tsx`:
- `emptyForm` constant
- `useDrawerForm`, `useCreateSupportRecord`, `useUpdateSupportRecord`, `usePerson`, `useOffices`, `useQueryClient` calls
- `personName` state and its two `useEffect`s
- `handleSubmit`
- `officeOptions` derivation

Hook returns `{ isEdit, form, set, error, isPending, handleSubmit, personName, setPersonName, officeOptions }`.

- [ ] **Step 2: Rewrite `index.tsx`** to use the hook (~80 lines).

- [ ] **Step 3: Verify and commit**

```bash
cd packages/observer-web && bun run build 2>&1 | grep -E "error|Error" | head -10
git add packages/observer-web/src/components/support-record-drawer/
git commit -m "Extract useSupportRecordForm hook to slim support-record-drawer index"
```

---

### Task 11: `household-drawer/index.tsx` — extract `useHouseholdDrawerForm`

**Files:**
- Create: `src/components/household-drawer/use-household-drawer-form.ts`
- Modify: `src/components/household-drawer/index.tsx`

- [ ] **Step 1: Create `use-household-drawer-form.ts`**

Extract from `index.tsx`:
- `emptyForm`, `emptyMemberForm` constants
- All hook calls (`useHousehold`, `useCreateHousehold`, `useUpdateHousehold`, `useAddHouseholdMember`, `useRemoveHouseholdMember`, `usePerson`, `useDrawerForm`, `useQueryClient`)
- `memberForm`, `headPersonName`, `memberPersonName` state
- `resolvedHeadName`, `headPersonLabel` derivations
- `handleSubmit`, `handleAddMember`, `handleRemoveMember`

Hook returns `{ isEdit, form, set, error, isPending, handleSubmit, handleAddMember, handleRemoveMember, memberForm, setMemberForm, headPersonLabel, memberPersonName, setMemberPersonName, addMemberPending, editingId, household }`.

- [ ] **Step 2: Rewrite `index.tsx`** using the hook (~60 lines).

- [ ] **Step 3: Verify and commit**

```bash
cd packages/observer-web && bun run build 2>&1 | grep -E "error|Error" | head -10
git add packages/observer-web/src/components/household-drawer/
git commit -m "Extract useHouseholdDrawerForm hook to slim household-drawer index"
```

---

## Phase 3 — Route sub-component extractions

> For each task: read the route file first, identify the exact JSX block to extract, create the component, import it in the route, verify build.

### Task 12: `people/$personId/index.lazy.tsx` (298 lines)

**Extract to:**
- `components/person-detail.tsx` — the `Detail` dl item
- `components/quick-support-form.tsx` — the inline quick-add support record form

- [ ] **Step 1: Create `person-detail.tsx`**

```tsx
// packages/observer-web/src/components/person-detail.tsx
interface PersonDetailProps {
  label: string;
  value?: string | null;
}

export function PersonDetail({ label, value }: PersonDetailProps) {
  return (
    <div>
      <dt className="text-xs font-medium text-fg-tertiary">{label}</dt>
      <dd className="mt-0.5 text-sm text-fg">{value}</dd>
    </div>
  );
}
```

- [ ] **Step 2: Create `quick-support-form.tsx`**

Read `people/$personId/index.lazy.tsx`. The quick support form starts at the `{formOpen && (` block. Extract into:

```tsx
// packages/observer-web/src/components/quick-support-form.tsx
// Props: projectId, personId, onSaved, onCancel
// Contains: type/sphere UISelects, DatePicker, textarea, submit/cancel buttons
// Internally calls useCreateSupportRecord, useToast, handleApiError
// typeOptions, sphereOptions, typeKeyMap, sphereKeyMap, inputClass all live here
```

The component manages its own `type`, `sphere`, `providedAt`, `notes`, `error` state and calls `onSaved()` on success.

- [ ] **Step 3: Update the route file**

In `index.lazy.tsx`:
- Remove the `Detail` function definition and the entire `{formOpen && (...)}` block
- Import `PersonDetail` and `QuickSupportForm`
- Replace all `<Detail .../>` with `<PersonDetail .../>`
- Replace `{formOpen && (...)}` with `{formOpen && <QuickSupportForm projectId={projectId} personId={personId} onSaved={() => setFormOpen(false)} onCancel={() => setFormOpen(false)} />}`
- Remove now-unused imports (`useState` for form fields, `DatePicker`, `UISelect`, `handleApiError`, `useCreateSupportRecord`, `useToast`, etc. — keep only what remains in the route)

- [ ] **Step 4: Verify and commit**

```bash
cd packages/observer-web && bun run build 2>&1 | grep -E "error|Error" | head -10
git add packages/observer-web/src/components/person-detail.tsx \
        packages/observer-web/src/components/quick-support-form.tsx \
        packages/observer-web/src/routes/_app/projects/\$projectId/people/\$personId/index.lazy.tsx
git commit -m "Extract PersonDetail and QuickSupportForm from person overview route"
```

---

### Task 13: `people/$personId/documents.lazy.tsx` (398 lines)

**Extract to:**
- `components/document-mime-icon.tsx` — `mimeIcon()` + `formatBytes()` utils
- `components/document-preview-dialog.tsx` — image/PDF preview modal
- `components/document-upload-zone.tsx` — upload button + file handling UI

- [ ] **Step 1: Read the file**

```bash
cat -n packages/observer-web/src/routes/_app/projects/\$projectId/people/\$personId/documents.lazy.tsx
```

- [ ] **Step 2: Create `document-mime-icon.tsx`**

Extract `formatBytes` and `mimeIcon` functions. Export both. Signature:

```tsx
export function formatBytes(bytes: number): string { ... }
export function mimeIcon(type: string): Icon { ... } // returns a Phosphor icon component
```

- [ ] **Step 3: Create `document-preview-dialog.tsx`**

Extract the preview dialog (the `@base-ui/react/dialog` block that shows images and PDFs). Props:

```tsx
interface DocumentPreviewDialogProps {
  document: Document | null;
  projectId: string;
  personId: string;
  onClose: () => void;
}
```

- [ ] **Step 4: Create `document-upload-zone.tsx`**

Extract the upload button and its hidden file input. Props:

```tsx
interface DocumentUploadZoneProps {
  projectId: string;
  personId: string;
  onUploadError: (msg: string) => void;
}
// Internally calls useUploadDocument
```

- [ ] **Step 5: Update route file and verify**

Replace extracted blocks with component calls. Remove now-unused icon imports.

```bash
cd packages/observer-web && bun run build 2>&1 | grep -E "error|Error" | head -10
git add packages/observer-web/src/components/document-mime-icon.tsx \
        packages/observer-web/src/components/document-preview-dialog.tsx \
        packages/observer-web/src/components/document-upload-zone.tsx \
        "packages/observer-web/src/routes/_app/projects/\$projectId/people/\$personId/documents.lazy.tsx"
git commit -m "Extract DocumentMimeIcon, DocumentPreviewDialog, DocumentUploadZone from documents route"
```

---

### Task 14: `people/index.tsx` (334 lines)

**Extract to:**
- `components/people-filter-bar.tsx` — all filter inputs (status tabs, sex, age group, date range, tags, office)
- `components/people-columns.tsx` — the `ColumnDef<Person>[]` factory function

- [ ] **Step 1: Read the file**

```bash
cat -n packages/observer-web/src/routes/_app/projects/\$projectId/people/index.tsx
```

- [ ] **Step 2: Create `people-columns.tsx`**

Extract the column definitions. Signature:

```tsx
// packages/observer-web/src/components/people-columns.tsx
import type { Column } from "@/components/data-table";
import type { Person } from "@/types/person";

interface PeopleColumnsOptions {
  t: (key: string) => string;
  onEdit: (id: string) => void;
  canWrite: boolean;
}

export function buildPeopleColumns({ t, onEdit, canWrite }: PeopleColumnsOptions): Column<Person>[] { ... }
```

- [ ] **Step 3: Create `people-filter-bar.tsx`**

Extract the filter controls. Signature:

```tsx
// packages/observer-web/src/components/people-filter-bar.tsx
interface PeopleFilterBarProps {
  projectId: string;
  params: ListPeopleParams;
  onParamChange: (key: keyof ListPeopleParams, value: unknown) => void;
}
```

- [ ] **Step 4: Update route and verify**

```bash
cd packages/observer-web && bun run build 2>&1 | grep -E "error|Error" | head -10
git add packages/observer-web/src/components/people-columns.tsx \
        packages/observer-web/src/components/people-filter-bar.tsx \
        "packages/observer-web/src/routes/_app/projects/\$projectId/people/index.tsx"
git commit -m "Extract PeopleColumns and PeopleFilterBar from people list route"
```

---

### Task 15: `admin/users/index.lazy.tsx` (312 lines)

**Extract to:**
- `components/users-columns.tsx` — user table column definitions
- `components/create-user-dialog.tsx` — the user creation form dialog

- [ ] **Step 1: Read the file**

```bash
cat -n packages/observer-web/src/routes/_app/admin/users/index.lazy.tsx
```

- [ ] **Step 2: Create `users-columns.tsx`**

```tsx
// packages/observer-web/src/components/users-columns.tsx
import type { Column } from "@/components/data-table";
import type { AdminUser } from "@/types/admin";

interface UsersColumnsOptions {
  t: (key: string) => string;
}

export function buildUsersColumns({ t }: UsersColumnsOptions): Column<AdminUser>[] { ... }
```

- [ ] **Step 3: Create `create-user-dialog.tsx`**

Extract the dialog with the user creation form. Props:

```tsx
interface CreateUserDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onCreated: () => void;
}
```

- [ ] **Step 4: Update route and verify**

```bash
cd packages/observer-web && bun run build 2>&1 | grep -E "error|Error" | head -10
git add packages/observer-web/src/components/users-columns.tsx \
        packages/observer-web/src/components/create-user-dialog.tsx \
        packages/observer-web/src/routes/_app/admin/users/index.lazy.tsx
git commit -m "Extract UsersColumns and CreateUserDialog from admin users route"
```

---

### Task 16: `support-records/-support-records-page.tsx` (309 lines)

**Extract to:**
- `components/support-record-filter-bar.tsx`
- `components/support-record-columns.tsx`

- [ ] **Step 1: Read the file**

```bash
cat -n "packages/observer-web/src/routes/_app/projects/\$projectId/support-records/-support-records-page.tsx"
```

- [ ] **Step 2: Create `support-record-columns.tsx`**

```tsx
// packages/observer-web/src/components/support-record-columns.tsx
import type { Column } from "@/components/data-table";
import type { SupportRecord } from "@/types/support-record";

interface SupportRecordColumnsOptions {
  t: (key: string) => string;
  projectId: string;
  canWrite: boolean;
  onEdit: (id: string) => void;
  onDelete: (id: string) => void;
}

export function buildSupportRecordColumns(opts: SupportRecordColumnsOptions): Column<SupportRecord>[] { ... }
```

- [ ] **Step 3: Create `support-record-filter-bar.tsx`**

Extract sphere filter + date range filter row. Props mirror the route's filter state.

- [ ] **Step 4: Update route and verify**

```bash
cd packages/observer-web && bun run build 2>&1 | grep -E "error|Error" | head -10
git add packages/observer-web/src/components/support-record-filter-bar.tsx \
        packages/observer-web/src/components/support-record-columns.tsx \
        "packages/observer-web/src/routes/_app/projects/\$projectId/support-records/-support-records-page.tsx"
git commit -m "Extract SupportRecordColumns and SupportRecordFilterBar from support-records route"
```

---

### Task 17: `my-stats/index.lazy.tsx` (307 lines)

**Extract to:**
- `components/report-date-presets.tsx` — preset date range buttons (shared with reports pages)
- `components/my-stats-kpi-cards.tsx` — KPI stat cards row

- [ ] **Step 1: Read the file**

```bash
cat -n "packages/observer-web/src/routes/_app/projects/\$projectId/my-stats/index.lazy.tsx"
```

- [ ] **Step 2: Create `report-date-presets.tsx`**

```tsx
// packages/observer-web/src/components/report-date-presets.tsx
interface DatePreset {
  label: string;
  from: string;
  to: string;
}

interface ReportDatePresetsProps {
  from: string;
  to: string;
  presets: DatePreset[];
  onChange: (from: string, to: string) => void;
}

export function ReportDatePresets({ from, to, presets, onChange }: ReportDatePresetsProps) { ... }
```

- [ ] **Step 3: Create `my-stats-kpi-cards.tsx`**

Extract the KPI cards grid (people served, consultations, active cases, households). The component receives raw counts and renders them.

- [ ] **Step 4: Update route and verify**

```bash
cd packages/observer-web && bun run build 2>&1 | grep -E "error|Error" | head -10
git add packages/observer-web/src/components/report-date-presets.tsx \
        packages/observer-web/src/components/my-stats-kpi-cards.tsx \
        "packages/observer-web/src/routes/_app/projects/\$projectId/my-stats/index.lazy.tsx"
git commit -m "Extract ReportDatePresets and MyStatsKpiCards from my-stats route"
```

---

### Task 18: `reports/people.lazy.tsx` (518 lines)

**Extract to `components/reports/`:**
- `report-filter-bar.tsx` — shared filter controls (office, category, case status, sex, age group, support type)
- `people-kpi-cards.tsx` — 6 KPI cards
- `people-chart-section.tsx` — the repeated chart+title+filter wrapper

- [ ] **Step 1: Read the file**

```bash
cat -n "packages/observer-web/src/routes/_app/projects/\$projectId/reports/people.lazy.tsx"
```

- [ ] **Step 2: Create `components/reports/` and `report-filter-bar.tsx`**

```tsx
// packages/observer-web/src/components/reports/report-filter-bar.tsx
// Contains: date range (using ReportDatePresets), office filter, category filter,
//           case status filter, sex filter, age group filter, support type filter
// Props: projectId + all filter values + onChange handlers
```

- [ ] **Step 3: Create `people-kpi-cards.tsx`**

```tsx
// packages/observer-web/src/components/reports/people-kpi-cards.tsx
interface PeopleKpiCardsProps {
  totalPeople: number;
  totalConsultations: number;
  activeCases: number;
  idpCount: number;
  households: number;
  offices: number;
  isLoading: boolean;
}

export function PeopleKpiCards(props: PeopleKpiCardsProps) { ... }
```

- [ ] **Step 4: Create `people-chart-section.tsx`**

A reusable wrapper for a single chart card with title, optional filter, and chart content:

```tsx
// packages/observer-web/src/components/reports/people-chart-section.tsx
interface PeopleChartSectionProps {
  title: string;
  filter?: React.ReactNode;
  children: React.ReactNode;
}

export function PeopleChartSection({ title, filter, children }: PeopleChartSectionProps) {
  return (
    <div className="rounded-xl border border-border-secondary bg-bg-secondary p-5">
      <div className="mb-4 flex items-center justify-between">
        <h3 className="text-sm font-semibold text-fg">{title}</h3>
        {filter}
      </div>
      {children}
    </div>
  );
}
```

- [ ] **Step 5: Update route and verify**

```bash
cd packages/observer-web && bun run build 2>&1 | grep -E "error|Error" | head -10
git add packages/observer-web/src/components/reports/
git add "packages/observer-web/src/routes/_app/projects/\$projectId/reports/people.lazy.tsx"
git commit -m "Extract ReportFilterBar, PeopleKpiCards, PeopleChartSection from people reports route"
```

---

### Task 19: `reports/pets.lazy.tsx` (465 lines)

**Extract to `components/reports/`:**
- `pets-kpi-cards.tsx`
- `pets-chart-section.tsx` (or reuse `PeopleChartSection` if identical)

- [ ] **Step 1: Read the file**

```bash
cat -n "packages/observer-web/src/routes/_app/projects/\$projectId/reports/pets.lazy.tsx"
```

- [ ] **Step 2: Create `pets-kpi-cards.tsx`**

```tsx
interface PetsKpiCardsProps {
  total: number;
  needsShelter: number;
  adopted: number;
  isLoading: boolean;
}

export function PetsKpiCards(props: PetsKpiCardsProps) { ... }
```

- [ ] **Step 3: Reuse or create chart section wrapper**

If the pets chart section markup matches `PeopleChartSection`, import and reuse it. If different, create `pets-chart-section.tsx`.

- [ ] **Step 4: Update route and verify**

```bash
cd packages/observer-web && bun run build 2>&1 | grep -E "error|Error" | head -10
git add packages/observer-web/src/components/reports/pets-kpi-cards.tsx \
        "packages/observer-web/src/routes/_app/projects/\$projectId/reports/pets.lazy.tsx"
git commit -m "Extract PetsKpiCards from pets reports route"
```

---

### Task 20: `reports/custom.lazy.tsx` (300 lines)

**Extract to `components/reports/`:**
- `custom-report-form.tsx` — metric + dimension + filter controls
- `report-result-table.tsx` — the dynamic results table

- [ ] **Step 1: Read the file**

```bash
cat -n "packages/observer-web/src/routes/_app/projects/\$projectId/reports/custom.lazy.tsx"
```

- [ ] **Step 2: Create `custom-report-form.tsx`**

```tsx
interface CustomReportFormProps {
  projectId: string;
  metric: string;
  dimensions: string[];
  from: string;
  to: string;
  supportType: string;
  onMetricChange: (v: string) => void;
  onDimensionsChange: (v: string[]) => void;
  onFromChange: (v: string) => void;
  onToChange: (v: string) => void;
  onSupportTypeChange: (v: string) => void;
}
```

- [ ] **Step 3: Create `report-result-table.tsx`**

Extract the results table (dynamic columns based on selected dimensions, CSV export button).

- [ ] **Step 4: Update route and verify**

```bash
cd packages/observer-web && bun run build 2>&1 | grep -E "error|Error" | head -10
git add packages/observer-web/src/components/reports/custom-report-form.tsx \
        packages/observer-web/src/components/reports/report-result-table.tsx \
        "packages/observer-web/src/routes/_app/projects/\$projectId/reports/custom.lazy.tsx"
git commit -m "Extract CustomReportForm and ReportResultTable from custom reports route"
```

---

### Task 21: `pets/-pets-page.tsx` (286 lines)

**Extract to:**
- `components/pet-filter-bar.tsx`
- `components/pet-columns.tsx`

- [ ] **Step 1: Read the file**

```bash
cat -n "packages/observer-web/src/routes/_app/projects/\$projectId/pets/-pets-page.tsx"
```

- [ ] **Step 2: Create `pet-columns.tsx`**

```tsx
import type { Column } from "@/components/data-table";
import type { Pet } from "@/types/pet";

interface PetColumnsOptions {
  t: (key: string) => string;
  canWrite: boolean;
  onEdit: (id: string) => void;
  onDelete: (id: string) => void;
}

export function buildPetColumns(opts: PetColumnsOptions): Column<Pet>[] { ... }
```

- [ ] **Step 3: Create `pet-filter-bar.tsx`**

Extract date range + tag filters.

- [ ] **Step 4: Update route and verify**

```bash
cd packages/observer-web && bun run build 2>&1 | grep -E "error|Error" | head -10
git add packages/observer-web/src/components/pet-columns.tsx \
        packages/observer-web/src/components/pet-filter-bar.tsx \
        "packages/observer-web/src/routes/_app/projects/\$projectId/pets/-pets-page.tsx"
git commit -m "Extract PetColumns and PetFilterBar from pets page route"
```

---

## Phase 4 — Final verification

### Task 22: Full build + line count audit

- [ ] **Step 1: Full production build**

```bash
cd packages/observer-web && bun run build 2>&1 | tail -20
```

Expected: no TypeScript errors, bundle sizes reported.

- [ ] **Step 2: Confirm no component file exceeds 170 lines**

```bash
find packages/observer-web/src/components -name "*.tsx" -o -name "*.ts" | \
  xargs wc -l | sort -rn | awk '$1 > 170 {print}' | grep -v total
```

Expected: empty output (no files over 170 lines).

- [ ] **Step 3: Confirm route files are unchanged structurally**

```bash
find packages/observer-web/src/routes -name "*.tsx" | \
  xargs wc -l | sort -rn | head -10
```

Routes may still be large — that is expected and correct.

- [ ] **Step 4: Final commit**

```bash
git add -A
git commit -m "Complete component extraction — all components under 170 lines"
```
