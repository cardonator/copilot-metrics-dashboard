import * as React from "react";
import { Column } from "@tanstack/react-table";
import { Check, PlusCircle } from "lucide-react";
import { cn } from "@/lib/utils";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList, CommandSeparator } from "@/components/ui/command";

interface DataTableFacetedFilterProps<TData, TValue> {
    column?: Column<TData, TValue>;
    title?: string;
}

export function DataTableFacetedFilter<TData, TValue>({ column, title }: DataTableFacetedFilterProps<TData, TValue>) {
    const facets = column?.getFacetedUniqueValues();
    const selectedValues = new Set(column?.getFilterValue() as string[]);

    return (
        <Popover>
            <PopoverTrigger asChild>
                <Button
                    variant="outline"
                    size="sm"
                    className="h-8 border-dashed gap-2">
                    <PlusCircle className="h-4 w-4" />
                    {title}
                    {selectedValues?.size > 0 && (
                        <>
                            <Separator
                                orientation="vertical"
                                className="mx-2 h-4"
                            />
                            <Badge
                                variant="secondary"
                                className="rounded-sm px-1 font-normal lg:hidden">
                                {selectedValues.size}
                            </Badge>
                            <div className="hidden space-x-1 lg:flex">
                                {selectedValues.size > 2 ? (
                                    <Badge
                                        variant="secondary"
                                        className="rounded-sm px-1 font-normal">
                                        {selectedValues.size} selected
                                    </Badge>
                                ) : (
                                    Array.from(selectedValues).map((value) => (
                                        <Badge
                                            key={value}
                                            variant="secondary"
                                            className="rounded-sm px-1 font-normal">
                                            {value}
                                        </Badge>
                                    ))
                                )}
                            </div>
                        </>
                    )}
                </Button>
            </PopoverTrigger>
            <PopoverContent
                className="w-[200px] p-0"
                align="start">
                <Command>
                    <CommandInput placeholder={title} key="command-input" />
                    <CommandList key="command-list">
                        <CommandEmpty key="command-empty">No results found.</CommandEmpty>
                        <CommandGroup key="facets-group">
                            {Array.from(facets ?? []).map(([value, count]) => {
                                const isSelected = selectedValues.has(value);
                                const keyValue = String(value || 'empty');
                                return (
                                    <CommandItem
                                        key={keyValue}
                                        onSelect={() => {
                                            if (isSelected) {
                                                selectedValues.delete(value);
                                            } else {
                                                selectedValues.add(value);
                                            }
                                            const filterValues = Array.from(selectedValues);
                                            column?.setFilterValue(filterValues.length ? filterValues : undefined);
                                        }}>
                                        <div className={cn("mr-2 flex h-4 w-4 items-center justify-center rounded-sm border border-primary", isSelected ? "bg-primary text-primary-foreground" : "opacity-50 [&_svg]:invisible")}>
                                            <Check />
                                        </div>
                                        <span>{value}</span>
                                        {count && <span className="ml-auto flex h-4 w-4 items-center justify-center font-mono text-xs">{count}</span>}
                                    </CommandItem>
                                );
                            })}
                        </CommandGroup>
                        {selectedValues.size > 0 && (
                            <React.Fragment key="clear-filters-section">
                                <CommandSeparator key="command-separator" />
                                <CommandGroup key="clear-filters-group">
                                    <CommandItem
                                        key="clear-filters-item"
                                        onSelect={() => column?.setFilterValue(undefined)}
                                        className="justify-center text-center">
                                        Clear filters
                                    </CommandItem>
                                </CommandGroup>
                            </React.Fragment>
                        )}
                    </CommandList>
                </Command>
            </PopoverContent>
        </Popover>
    );
}
