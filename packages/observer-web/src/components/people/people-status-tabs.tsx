import { Tabs } from "@base-ui/react/tabs";

interface PeopleStatusTabsProps {
  tabs: readonly string[];
  value: string;
  onValueChange: (value: string) => void;
  getLabel: (tab: string) => string;
}

export function PeopleStatusTabs({ tabs, value, onValueChange, getLabel }: PeopleStatusTabsProps) {
  return (
    <Tabs.Root defaultValue="" value={value} onValueChange={onValueChange} className="mb-4">
      <Tabs.List className="flex gap-0 rounded-lg border border-border-secondary bg-bg-secondary p-0.5">
        {tabs.map((tab) => (
          <Tabs.Tab
            key={tab}
            value={tab}
            className="cursor-pointer rounded-sm px-4 py-1.5 m-0.5 text-sm font-medium text-fg-tertiary transition-colors hover:text-fg data-active:bg-bg data-active:text-fg data-active:shadow-card"
          >
            {getLabel(tab)}
          </Tabs.Tab>
        ))}
      </Tabs.List>
    </Tabs.Root>
  );
}
