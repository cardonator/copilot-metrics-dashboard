import * as React from "react";
import { Column } from "@tanstack/react-table";
import { Check, PlusCircle } from "lucide-react";
import { cn } from "@/lib/utils";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList, CommandSeparator } from "@/components/ui/command";

interface DataTableLanguageFilterProps<TData, TValue> {
    column?: Column<TData, TValue>;
    title?: string;
}

export function DataTableLanguageFilter<TData, TValue>({ column, title }: DataTableLanguageFilterProps<TData, TValue>) {
    const facets = column?.getFacetedUniqueValues();
    const selectedValues = new Set(column?.getFilterValue() as string[]);
    const [searchValue, setSearchValue] = React.useState("");

    // Extract individual languages from comma-separated lists
    const languageOptions = React.useMemo(() => {
        const uniqueLanguages = new Set<string>();
        if (facets) {
            Array.from(facets.keys()).forEach(value => {
                if (typeof value === 'string' && value.toLowerCase() !== 'n/a') {
                    // Split comma-separated language lists and trim whitespace
                    const languages = value.split(',').map(lang => lang.trim());
                    languages.forEach(lang => {
                        if (lang) uniqueLanguages.add(lang);
                    });
                }
            });
        }
        return Array.from(uniqueLanguages).sort();
    }, [facets]);

    // Filter language options based on search term
    const filteredLanguageOptions = React.useMemo(() => {
        if (!searchValue) return languageOptions;
        return languageOptions.filter(language => 
            language.toLowerCase().includes(searchValue.toLowerCase())
        );
    }, [languageOptions, searchValue]);

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
                    <CommandInput 
                        placeholder={title} 
                        value={searchValue}
                        onValueChange={setSearchValue}
                    />
                    <CommandList>
                        <CommandEmpty>No languages found.</CommandEmpty>
                        <CommandGroup>
                            {filteredLanguageOptions.map((language) => {
                                const isSelected = selectedValues.has(language);
                                return (
                                    <CommandItem
                                        key={language}
                                        onSelect={() => {
                                            if (isSelected) {
                                                selectedValues.delete(language);
                                            } else {
                                                selectedValues.add(language);
                                            }
                                            const filterValues = Array.from(selectedValues);
                                            column?.setFilterValue(filterValues.length ? filterValues : undefined);
                                        }}>
                                        <div className={cn("mr-2 flex h-4 w-4 items-center justify-center rounded-sm border border-primary", isSelected ? "bg-primary text-primary-foreground" : "opacity-50 [&_svg]:invisible")}>
                                            <Check />
                                        </div>
                                        <span>{language}</span>
                                    </CommandItem>
                                );
                            })}
                        </CommandGroup>
                        {selectedValues.size > 0 && (
                            <React.Fragment>
                                <CommandSeparator />
                                <CommandGroup>
                                    <CommandItem
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