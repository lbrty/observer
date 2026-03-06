# Bold Page Backgrounds Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add two-layer decorative SVG backgrounds to every page — a base topo layer on the app shell and per-zone thematic accent illustrations.

**Architecture:** Pure CSS using `mask-image` with inline SVG data URIs, applied via `::before` pseudo-elements. Base layer on `_app.tsx` content wrapper, accent classes on individual page root divs. All patterns use `var(--fg)` so they adapt to all 4 themes.

**Tech Stack:** CSS (Tailwind v4 + custom classes in `main.css`), inline SVG data URIs


### Task 2: Add accent pattern CSS classes

**Files:**
- Modify: `packages/observer-web/src/main.css` (add 7 `.page-bg-*` accent classes after `.page-bg-base`)

**Step 1: Add all 7 accent classes to `main.css`**

Add after the `.page-bg-base` block. Each class uses `::after` (base uses `::before`) so they layer correctly. All share the same positioning logic: top area, offset ~15% from right, ~6-8% opacity.

```css
/* Page accent — shared structure */
.page-bg-dashboard,
.page-bg-people,
.page-bg-reports,
.page-bg-admin,
.page-bg-support,
.page-bg-profile,
.page-bg-reference {
  position: relative;
  overflow: hidden;
}
.page-bg-dashboard::after,
.page-bg-people::after,
.page-bg-reports::after,
.page-bg-admin::after,
.page-bg-support::after,
.page-bg-profile::after,
.page-bg-reference::after {
  content: "";
  position: absolute;
  top: -20px;
  right: 10%;
  width: 320px;
  height: 320px;
  background-color: var(--fg);
  pointer-events: none;
  -webkit-mask-repeat: no-repeat;
  mask-repeat: no-repeat;
  -webkit-mask-position: center;
  mask-position: center;
  -webkit-mask-size: contain;
  mask-size: contain;
}

/* Dashboard — compass rose */
.page-bg-dashboard::after {
  opacity: 0.06;
  -webkit-mask-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Ccircle cx='100' cy='100' r='90' fill='none' stroke='black' stroke-width='1'/%3E%3Ccircle cx='100' cy='100' r='60' fill='none' stroke='black' stroke-width='0.8'/%3E%3Ccircle cx='100' cy='100' r='30' fill='none' stroke='black' stroke-width='0.6'/%3E%3Cline x1='100' y1='5' x2='100' y2='195' stroke='black' stroke-width='1'/%3E%3Cline x1='5' y1='100' x2='195' y2='100' stroke='black' stroke-width='1'/%3E%3Cline x1='30' y1='30' x2='170' y2='170' stroke='black' stroke-width='0.6'/%3E%3Cline x1='170' y1='30' x2='30' y2='170' stroke='black' stroke-width='0.6'/%3E%3Cpolygon points='100,8 106,40 100,30 94,40' fill='black'/%3E%3Cpolygon points='192,100 160,106 170,100 160,94' fill='black'/%3E%3Cpolygon points='100,192 94,160 100,170 106,160' fill='black'/%3E%3Cpolygon points='8,100 40,94 30,100 40,106' fill='black'/%3E%3Ctext x='100' y='20' text-anchor='middle' font-size='10' font-weight='bold' fill='black'%3EN%3C/text%3E%3Ctext x='100' y='190' text-anchor='middle' font-size='8' fill='black'%3ES%3C/text%3E%3Ctext x='188' y='104' text-anchor='middle' font-size='8' fill='black'%3EE%3C/text%3E%3Ctext x='14' y='104' text-anchor='middle' font-size='8' fill='black'%3EW%3C/text%3E%3C/svg%3E");
  mask-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Ccircle cx='100' cy='100' r='90' fill='none' stroke='black' stroke-width='1'/%3E%3Ccircle cx='100' cy='100' r='60' fill='none' stroke='black' stroke-width='0.8'/%3E%3Ccircle cx='100' cy='100' r='30' fill='none' stroke='black' stroke-width='0.6'/%3E%3Cline x1='100' y1='5' x2='100' y2='195' stroke='black' stroke-width='1'/%3E%3Cline x1='5' y1='100' x2='195' y2='100' stroke='black' stroke-width='1'/%3E%3Cline x1='30' y1='30' x2='170' y2='170' stroke='black' stroke-width='0.6'/%3E%3Cline x1='170' y1='30' x2='30' y2='170' stroke='black' stroke-width='0.6'/%3E%3Cpolygon points='100,8 106,40 100,30 94,40' fill='black'/%3E%3Cpolygon points='192,100 160,106 170,100 160,94' fill='black'/%3E%3Cpolygon points='100,192 94,160 100,170 106,160' fill='black'/%3E%3Cpolygon points='8,100 40,94 30,100 40,106' fill='black'/%3E%3Ctext x='100' y='20' text-anchor='middle' font-size='10' font-weight='bold' fill='black'%3EN%3C/text%3E%3Ctext x='100' y='190' text-anchor='middle' font-size='8' fill='black'%3ES%3C/text%3E%3Ctext x='188' y='104' text-anchor='middle' font-size='8' fill='black'%3EE%3C/text%3E%3Ctext x='14' y='104' text-anchor='middle' font-size='8' fill='black'%3EW%3C/text%3E%3C/svg%3E");
}

/* People — overlapping circles (community) */
.page-bg-people::after {
  opacity: 0.07;
  -webkit-mask-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Ccircle cx='70' cy='80' r='45' fill='none' stroke='black' stroke-width='1.5'/%3E%3Ccircle cx='130' cy='80' r='45' fill='none' stroke='black' stroke-width='1.5'/%3E%3Ccircle cx='100' cy='130' r='45' fill='none' stroke='black' stroke-width='1.5'/%3E%3Ccircle cx='100' cy='95' r='15' fill='none' stroke='black' stroke-width='1'/%3E%3Ccircle cx='70' cy='80' r='8' fill='black' fill-opacity='0.3'/%3E%3Ccircle cx='130' cy='80' r='8' fill='black' fill-opacity='0.3'/%3E%3Ccircle cx='100' cy='130' r='8' fill='black' fill-opacity='0.3'/%3E%3C/svg%3E");
  mask-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Ccircle cx='70' cy='80' r='45' fill='none' stroke='black' stroke-width='1.5'/%3E%3Ccircle cx='130' cy='80' r='45' fill='none' stroke='black' stroke-width='1.5'/%3E%3Ccircle cx='100' cy='130' r='45' fill='none' stroke='black' stroke-width='1.5'/%3E%3Ccircle cx='100' cy='95' r='15' fill='none' stroke='black' stroke-width='1'/%3E%3Ccircle cx='70' cy='80' r='8' fill='black' fill-opacity='0.3'/%3E%3Ccircle cx='130' cy='80' r='8' fill='black' fill-opacity='0.3'/%3E%3Ccircle cx='100' cy='130' r='8' fill='black' fill-opacity='0.3'/%3E%3C/svg%3E");
}

/* Reports — rising bar chart silhouette */
.page-bg-reports::after {
  opacity: 0.07;
  -webkit-mask-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Crect x='15' y='120' width='22' height='65' rx='3' fill='black'/%3E%3Crect x='47' y='90' width='22' height='95' rx='3' fill='black'/%3E%3Crect x='79' y='60' width='22' height='125' rx='3' fill='black'/%3E%3Crect x='111' y='40' width='22' height='145' rx='3' fill='black'/%3E%3Crect x='143' y='70' width='22' height='115' rx='3' fill='black'/%3E%3Crect x='175' y='25' width='22' height='160' rx='3' fill='black'/%3E%3Cpath d='M26,115 C50,82 82,52 122,35 C150,24 178,20 186,20' fill='none' stroke='black' stroke-width='1.5' stroke-dasharray='4,3'/%3E%3Ccircle cx='26' cy='115' r='3' fill='black'/%3E%3Ccircle cx='58' cy='85' r='3' fill='black'/%3E%3Ccircle cx='90' cy='55' r='3' fill='black'/%3E%3Ccircle cx='122' cy='35' r='3' fill='black'/%3E%3Ccircle cx='154' cy='65' r='3' fill='black'/%3E%3Ccircle cx='186' cy='20' r='3' fill='black'/%3E%3C/svg%3E");
  mask-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Crect x='15' y='120' width='22' height='65' rx='3' fill='black'/%3E%3Crect x='47' y='90' width='22' height='95' rx='3' fill='black'/%3E%3Crect x='79' y='60' width='22' height='125' rx='3' fill='black'/%3E%3Crect x='111' y='40' width='22' height='145' rx='3' fill='black'/%3E%3Crect x='143' y='70' width='22' height='115' rx='3' fill='black'/%3E%3Crect x='175' y='25' width='22' height='160' rx='3' fill='black'/%3E%3Cpath d='M26,115 C50,82 82,52 122,35 C150,24 178,20 186,20' fill='none' stroke='black' stroke-width='1.5' stroke-dasharray='4,3'/%3E%3Ccircle cx='26' cy='115' r='3' fill='black'/%3E%3Ccircle cx='58' cy='85' r='3' fill='black'/%3E%3Ccircle cx='90' cy='55' r='3' fill='black'/%3E%3Ccircle cx='122' cy='35' r='3' fill='black'/%3E%3Ccircle cx='154' cy='65' r='3' fill='black'/%3E%3Ccircle cx='186' cy='20' r='3' fill='black'/%3E%3C/svg%3E");
}

/* Admin — hexagonal grid */
.page-bg-admin::after {
  opacity: 0.06;
  -webkit-mask-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cpolygon points='60,20 90,20 105,46 90,72 60,72 45,46' fill='none' stroke='black' stroke-width='1.2'/%3E%3Cpolygon points='110,20 140,20 155,46 140,72 110,72 95,46' fill='none' stroke='black' stroke-width='1.2'/%3E%3Cpolygon points='35,62 65,62 80,88 65,114 35,114 20,88' fill='none' stroke='black' stroke-width='1.2'/%3E%3Cpolygon points='85,62 115,62 130,88 115,114 85,114 70,88' fill='none' stroke='black' stroke-width='1.2'/%3E%3Cpolygon points='135,62 165,62 180,88 165,114 135,114 120,88' fill='none' stroke='black' stroke-width='1.2'/%3E%3Cpolygon points='60,104 90,104 105,130 90,156 60,156 45,130' fill='none' stroke='black' stroke-width='1.2'/%3E%3Cpolygon points='110,104 140,104 155,130 140,156 110,156 95,130' fill='none' stroke='black' stroke-width='1.2'/%3E%3Cpolygon points='35,146 65,146 80,172 65,198 35,198 20,172' fill='none' stroke='black' stroke-width='1.2'/%3E%3Cpolygon points='85,146 115,146 130,172 115,198 85,198 70,172' fill='none' stroke='black' stroke-width='1.2'/%3E%3Cpolygon points='135,146 165,146 180,172 165,198 135,198 120,172' fill='none' stroke='black' stroke-width='1.2'/%3E%3Ccircle cx='75' cy='46' r='4' fill='black' fill-opacity='0.4'/%3E%3Ccircle cx='100' cy='88' r='4' fill='black' fill-opacity='0.4'/%3E%3Ccircle cx='75' cy='130' r='4' fill='black' fill-opacity='0.4'/%3E%3C/svg%3E");
  mask-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cpolygon points='60,20 90,20 105,46 90,72 60,72 45,46' fill='none' stroke='black' stroke-width='1.2'/%3E%3Cpolygon points='110,20 140,20 155,46 140,72 110,72 95,46' fill='none' stroke='black' stroke-width='1.2'/%3E%3Cpolygon points='35,62 65,62 80,88 65,114 35,114 20,88' fill='none' stroke='black' stroke-width='1.2'/%3E%3Cpolygon points='85,62 115,62 130,88 115,114 85,114 70,88' fill='none' stroke='black' stroke-width='1.2'/%3E%3Cpolygon points='135,62 165,62 180,88 165,114 135,114 120,88' fill='none' stroke='black' stroke-width='1.2'/%3E%3Cpolygon points='60,104 90,104 105,130 90,156 60,156 45,130' fill='none' stroke='black' stroke-width='1.2'/%3E%3Cpolygon points='110,104 140,104 155,130 140,156 110,156 95,130' fill='none' stroke='black' stroke-width='1.2'/%3E%3Cpolygon points='35,146 65,146 80,172 65,198 35,198 20,172' fill='none' stroke='black' stroke-width='1.2'/%3E%3Cpolygon points='85,146 115,146 130,172 115,198 85,198 70,172' fill='none' stroke='black' stroke-width='1.2'/%3E%3Cpolygon points='135,146 165,146 180,172 165,198 135,198 120,172' fill='none' stroke='black' stroke-width='1.2'/%3E%3Ccircle cx='75' cy='46' r='4' fill='black' fill-opacity='0.4'/%3E%3Ccircle cx='100' cy='88' r='4' fill='black' fill-opacity='0.4'/%3E%3Ccircle cx='75' cy='130' r='4' fill='black' fill-opacity='0.4'/%3E%3C/svg%3E");
}

/* Support records — interlocking chain links */
.page-bg-support::after {
  opacity: 0.065;
  -webkit-mask-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cellipse cx='65' cy='60' rx='40' ry='25' fill='none' stroke='black' stroke-width='1.5' transform='rotate(-20,65,60)'/%3E%3Cellipse cx='105' cy='75' rx='40' ry='25' fill='none' stroke='black' stroke-width='1.5' transform='rotate(-20,105,75)'/%3E%3Cellipse cx='80' cy='120' rx='40' ry='25' fill='none' stroke='black' stroke-width='1.5' transform='rotate(15,80,120)'/%3E%3Cellipse cx='120' cy='135' rx='40' ry='25' fill='none' stroke='black' stroke-width='1.5' transform='rotate(15,120,135)'/%3E%3Cellipse cx='95' cy='175' rx='40' ry='25' fill='none' stroke='black' stroke-width='1.5' transform='rotate(-10,95,175)'/%3E%3Ccircle cx='85' cy='68' r='3' fill='black'/%3E%3Ccircle cx='100' cy='128' r='3' fill='black'/%3E%3C/svg%3E");
  mask-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cellipse cx='65' cy='60' rx='40' ry='25' fill='none' stroke='black' stroke-width='1.5' transform='rotate(-20,65,60)'/%3E%3Cellipse cx='105' cy='75' rx='40' ry='25' fill='none' stroke='black' stroke-width='1.5' transform='rotate(-20,105,75)'/%3E%3Cellipse cx='80' cy='120' rx='40' ry='25' fill='none' stroke='black' stroke-width='1.5' transform='rotate(15,80,120)'/%3E%3Cellipse cx='120' cy='135' rx='40' ry='25' fill='none' stroke='black' stroke-width='1.5' transform='rotate(15,120,135)'/%3E%3Cellipse cx='95' cy='175' rx='40' ry='25' fill='none' stroke='black' stroke-width='1.5' transform='rotate(-10,95,175)'/%3E%3Ccircle cx='85' cy='68' r='3' fill='black'/%3E%3Ccircle cx='100' cy='128' r='3' fill='black'/%3E%3C/svg%3E");
}

/* Profile — circle with radiating lines */
.page-bg-profile::after {
  opacity: 0.06;
  width: 240px;
  height: 240px;
  -webkit-mask-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Ccircle cx='100' cy='100' r='35' fill='none' stroke='black' stroke-width='1.5'/%3E%3Ccircle cx='100' cy='100' r='55' fill='none' stroke='black' stroke-width='0.8'/%3E%3Ccircle cx='100' cy='100' r='12' fill='black' fill-opacity='0.3'/%3E%3Cline x1='100' y1='30' x2='100' y2='10' stroke='black' stroke-width='1'/%3E%3Cline x1='100' y1='170' x2='100' y2='190' stroke='black' stroke-width='1'/%3E%3Cline x1='30' y1='100' x2='10' y2='100' stroke='black' stroke-width='1'/%3E%3Cline x1='170' y1='100' x2='190' y2='100' stroke='black' stroke-width='1'/%3E%3Cline x1='50' y1='50' x2='36' y2='36' stroke='black' stroke-width='0.8'/%3E%3Cline x1='150' y1='50' x2='164' y2='36' stroke='black' stroke-width='0.8'/%3E%3Cline x1='50' y1='150' x2='36' y2='164' stroke='black' stroke-width='0.8'/%3E%3Cline x1='150' y1='150' x2='164' y2='164' stroke='black' stroke-width='0.8'/%3E%3C/svg%3E");
  mask-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Ccircle cx='100' cy='100' r='35' fill='none' stroke='black' stroke-width='1.5'/%3E%3Ccircle cx='100' cy='100' r='55' fill='none' stroke='black' stroke-width='0.8'/%3E%3Ccircle cx='100' cy='100' r='12' fill='black' fill-opacity='0.3'/%3E%3Cline x1='100' y1='30' x2='100' y2='10' stroke='black' stroke-width='1'/%3E%3Cline x1='100' y1='170' x2='100' y2='190' stroke='black' stroke-width='1'/%3E%3Cline x1='30' y1='100' x2='10' y2='100' stroke='black' stroke-width='1'/%3E%3Cline x1='170' y1='100' x2='190' y2='100' stroke='black' stroke-width='1'/%3E%3Cline x1='50' y1='50' x2='36' y2='36' stroke='black' stroke-width='0.8'/%3E%3Cline x1='150' y1='50' x2='164' y2='36' stroke='black' stroke-width='0.8'/%3E%3Cline x1='50' y1='150' x2='36' y2='164' stroke='black' stroke-width='0.8'/%3E%3Cline x1='150' y1='150' x2='164' y2='164' stroke='black' stroke-width='0.8'/%3E%3C/svg%3E");
}

/* Reference data — globe with latitude lines */
.page-bg-reference::after {
  opacity: 0.065;
  -webkit-mask-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Ccircle cx='100' cy='100' r='80' fill='none' stroke='black' stroke-width='1.5'/%3E%3Cellipse cx='100' cy='100' rx='50' ry='80' fill='none' stroke='black' stroke-width='1'/%3E%3Cellipse cx='100' cy='100' rx='20' ry='80' fill='none' stroke='black' stroke-width='0.8'/%3E%3Cpath d='M22,70 C50,65 75,62 100,62 C125,62 150,65 178,70' fill='none' stroke='black' stroke-width='0.8'/%3E%3Cpath d='M22,130 C50,135 75,138 100,138 C125,138 150,135 178,130' fill='none' stroke='black' stroke-width='0.8'/%3E%3Cline x1='20' y1='100' x2='180' y2='100' stroke='black' stroke-width='1'/%3E%3Ccircle cx='100' cy='100' r='3' fill='black'/%3E%3Ccircle cx='60' cy='68' r='2' fill='black'/%3E%3Ccircle cx='140' cy='132' r='2' fill='black'/%3E%3C/svg%3E");
  mask-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Ccircle cx='100' cy='100' r='80' fill='none' stroke='black' stroke-width='1.5'/%3E%3Cellipse cx='100' cy='100' rx='50' ry='80' fill='none' stroke='black' stroke-width='1'/%3E%3Cellipse cx='100' cy='100' rx='20' ry='80' fill='none' stroke='black' stroke-width='0.8'/%3E%3Cpath d='M22,70 C50,65 75,62 100,62 C125,62 150,65 178,70' fill='none' stroke='black' stroke-width='0.8'/%3E%3Cpath d='M22,130 C50,135 75,138 100,138 C125,138 150,135 178,130' fill='none' stroke='black' stroke-width='0.8'/%3E%3Cline x1='20' y1='100' x2='180' y2='100' stroke='black' stroke-width='1'/%3E%3Ccircle cx='100' cy='100' r='3' fill='black'/%3E%3Ccircle cx='60' cy='68' r='2' fill='black'/%3E%3Ccircle cx='140' cy='132' r='2' fill='black'/%3E%3C/svg%3E");
}
```

**Step 2: Commit**

```bash
git add packages/observer-web/src/main.css
git commit -m "add 7 page-level accent background pattern classes"
```


### Task 4: Apply accent classes to Project pages

**Files:**
- Modify: `packages/observer-web/src/routes/_app/projects/$projectId/people/index.tsx` (people list)
- Modify: `packages/observer-web/src/routes/_app/projects/$projectId/people/$personId.tsx` (person detail layout)
- Modify: `packages/observer-web/src/routes/_app/projects/$projectId/support-records/index.tsx`
- Modify: `packages/observer-web/src/routes/_app/projects/$projectId/reports/index.tsx`

**Step 1: People list — add `page-bg-people`**

Read `packages/observer-web/src/routes/_app/projects/$projectId/people/index.tsx` and add `page-bg-people` to the root div of the page component.

**Step 2: Person detail — add `page-bg-people`**

Read `packages/observer-web/src/routes/_app/projects/$projectId/people/$personId.tsx` and add `page-bg-people` to the root div.

**Step 3: Support records — add `page-bg-support`**

Read `packages/observer-web/src/routes/_app/projects/$projectId/support-records/index.tsx` and add `page-bg-support` to the root div.

**Step 4: Reports — add `page-bg-reports`**

Read `packages/observer-web/src/routes/_app/projects/$projectId/reports/index.tsx` and add `page-bg-reports` to the root div.

**Step 5: Verify visually**

Run: `cd packages/observer-web && bun dev`

Check: People pages show overlapping circles, support records show chain links, reports show bar chart silhouette.

**Step 6: Commit**

```bash
git add 'packages/observer-web/src/routes/_app/projects/$projectId/people/index.tsx' \
  'packages/observer-web/src/routes/_app/projects/$projectId/people/$personId.tsx' \
  'packages/observer-web/src/routes/_app/projects/$projectId/support-records/index.tsx' \
  'packages/observer-web/src/routes/_app/projects/$projectId/reports/index.tsx'
git commit -m "apply accent backgrounds to people, support, reports pages"
```


### Task 6: Print stylesheet exclusion

**Files:**
- Modify: `packages/observer-web/src/main.css` (add rule to `@media print` block)

**Step 1: Hide decorative backgrounds in print**

In the existing `@media print` block in `main.css` (around line 318), add:

```css
  .page-bg-base::before,
  .page-bg-dashboard::after,
  .page-bg-people::after,
  .page-bg-reports::after,
  .page-bg-admin::after,
  .page-bg-support::after,
  .page-bg-profile::after,
  .page-bg-reference::after {
    display: none !important;
  }
```

**Step 2: Commit**

```bash
git add packages/observer-web/src/main.css
git commit -m "hide page backgrounds in print stylesheet"
```
