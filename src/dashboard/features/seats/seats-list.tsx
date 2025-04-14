"use client";
import { useDashboard } from "./seats-state";
import { ChartHeader } from "@/features/common/chart-header";
import { Card, CardContent } from "@/components/ui/card";
import { stringIsNullOrEmpty } from "@/utils/helpers";
import { DataTable } from "@/components/ui/data-table";
import { ColumnDef } from "@tanstack/react-table";
import { DataTableColumnHeader } from "@/components/ui/data-table-column-header";
import { useState, useEffect } from "react";

interface SeatData {
    user: string;
    organization: string | null;
    team: string | null;
    createdAt: string;
    updatedAt: string;
    lastActivityAt: string;
    lastActivityEditor: string;
    planType: string;
    pendingCancellationDate: string;
    // New fields for usage data
    acceptanceRate: string;
    totalSuggestions: string;
    activeDays: string;
    timeSaved: string;
    mostUsedLanguages: string;
}

function formatEditorName(editor: string): string {
    if (stringIsNullOrEmpty(editor)) {
        return editor;
    }
    const editorInfo = editor.split('/');
    const editorName = `${editorInfo[0]} (${editorInfo[1]})`;

    return editorName;
}

const arrayIncludes = (row: any, id: string, value: any[]) => {
    return value.includes(row.getValue(id));
};

const stringIncludes = (row: any, id: string, value: string) => {
    return row.getValue(id).includes(value);
};

// Custom filter function for comma-separated language lists
const languageFilter = (row: any, id: string, filterValue: any) => {
    const languagesString = String(row.getValue(id)).toLowerCase();
    
    // Skip empty or N/A values
    if (languagesString === 'n/a' || !languagesString) {
        return false;
    }

    // Split the row's languages into an array of individual languages
    const rowLanguages = languagesString.split(',').map(lang => lang.trim());
    
    // For array filter values (selected from dropdown)
    if (Array.isArray(filterValue)) {
        // Check if ANY of the selected languages match 
        return filterValue.some(selectedLanguage => 
            // Make sure we're doing an exact match for each language
            rowLanguages.some(rowLanguage => 
                rowLanguage === selectedLanguage.toLowerCase()
            )
        );
    }
    
    // For string filter values (used in search inputs)
    if (filterValue && typeof filterValue === 'string') {
        const filterString = filterValue.toLowerCase();
        // Check if any language contains the search string
        return rowLanguages.some(language => language.includes(filterString));
    }
    
    // Default to true if the filter value is invalid
    return true;
};

const columns: ColumnDef<SeatData>[] = [
    // Row number column with custom render logic for accurate indexing
    {
        id: "rowNumber",
        header: () => <div className="text-center">#</div>,
        enableSorting: false,
        enableHiding: false,
        size: 50,
        cell: info => {
            // Use flatRows which accounts for filtering
            const rowIndex = info.table.getSortedRowModel().flatRows.findIndex(row => row.id === info.row.id);
            return <div className="text-center">{rowIndex + 1}</div>;
        },
    },
    // Rest of the columns with proper headers
    {
        accessorKey: "user",
        header: ({ column }) => <DataTableColumnHeader column={column} title="User" />,
        cell: ({ row }) => <div className="ml-2">{row.getValue("user")}</div>,
        filterFn: stringIncludes,
        size: 130,
    },
    {
        accessorKey: "acceptanceRate",
        header: ({ column }) => <DataTableColumnHeader column={column} title="Accept Rate" />,
        cell: ({ row }) => <div className="ml-2 text-center">{row.getValue("acceptanceRate")}</div>,
        size: 100,
    },
    {
        accessorKey: "totalSuggestions",
        header: ({ column }) => <DataTableColumnHeader column={column} title="Suggestions" />,
        cell: ({ row }) => <div className="ml-2 text-center">{row.getValue("totalSuggestions")}</div>,
        size: 110,
    },
    {
        accessorKey: "activeDays",
        header: ({ column }) => <DataTableColumnHeader column={column} title="Active Days" />,
        cell: ({ row }) => <div className="ml-2 text-center">{row.getValue("activeDays")}</div>,
        size: 100,
    },
    {
        accessorKey: "timeSaved",
        header: ({ column }) => <DataTableColumnHeader column={column} title="Time Saved" />,
        cell: ({ row }) => <div className="ml-2 text-center">{row.getValue("timeSaved")}</div>,
        size: 100,
    },
    {
        accessorKey: "mostUsedLanguages",
        header: ({ column }) => <DataTableColumnHeader column={column} title="Languages" />,
        cell: ({ row }) => {
            const languages = String(row.getValue("mostUsedLanguages"));
            
            // If no languages or N/A, just return as is
            if (languages === 'N/A' || !languages) {
                return <div className="ml-2">{languages}</div>;
            }
            
            // Split by comma and display each language on a new line
            return (
                <div className="ml-2">
                    {languages.split(',').map((lang, index) => (
                        <div key={index} className="whitespace-nowrap">{lang.trim()}</div>
                    ))}
                </div>
            );
        },
        filterFn: languageFilter,
        size: 120,
    },
    // Moved lastActivityAt and lastActivityEditor to here (after mostUsedLanguages)
    {
        accessorKey: "lastActivityAt",
        header: ({ column }) => <DataTableColumnHeader column={column} title="Last Activity" />,
        cell: ({ row }) => <div className="ml-2">{row.getValue("lastActivityAt")}</div>,
        size: 110,
    },
    {
        accessorKey: "lastActivityEditor",
        header: ({ column }) => <DataTableColumnHeader column={column} title="Last Editor" />,
        cell: ({ row }) => <div className="ml-2">{row.getValue("lastActivityEditor")}</div>,
        size: 140,
    },
    {
        accessorKey: "organization",
        header: ({ column }) => <DataTableColumnHeader column={column} title="Organization" />,
        cell: ({ row }) => <div className="ml-2">{row.getValue("organization")}</div>,
        filterFn: arrayIncludes,
    },
    {
        accessorKey: "team",
        header: ({ column }) => <DataTableColumnHeader column={column} title="Team" />,
        cell: ({ row }) => <div className="ml-2">{row.getValue("team")}</div>,
        filterFn: arrayIncludes,
    },
    {
        accessorKey: "createdAt",
        header: ({ column }) => <DataTableColumnHeader column={column} title="Create Date" />,
        cell: ({ row }) => <div className="ml-2">{row.getValue("createdAt")}</div>,
    },
    {
        accessorKey: "updatedAt",
        header: ({ column }) => <DataTableColumnHeader column={column} title="Update Date" />,
        cell: ({ row }) => <div className="ml-2">{row.getValue("updatedAt")}</div>,
    },
    {
        accessorKey: "planType",
        header: ({ column }) => <DataTableColumnHeader column={column} title="Plan" />,
        cell: ({ row }) => <div className="ml-2">{row.getValue("planType")}</div>,
        filterFn: arrayIncludes,
    },
    {
        accessorKey: "pendingCancellationDate",
        header: ({ column }) => <DataTableColumnHeader column={column} title="Pending Cancellation" />,
        cell: ({ row }) => <div className="ml-2">{row.getValue("pendingCancellationDate")}</div>,
    },
];

export const SeatsList = () => {
    const { filteredData } = useDashboard();
    const currentData = filteredData;
    const hasOrganization = currentData?.seats.some((seat) => seat.organization);
    const hasTeam = currentData?.seats.some((seat) => seat.assigning_team);
    
    // State to hold the user metrics
    const [userMetrics, setUserMetrics] = useState<Record<string, any>>({});
    const [isLoading, setIsLoading] = useState(true);
    
    // Fetch user metrics when the component mounts or filtered data changes
    useEffect(() => {
        const fetchUserMetrics = async () => {
            if (!currentData || !currentData.organization) return;
            
            setIsLoading(true);
            try {
                // Call server action to get user metrics
                const response = await fetch(`/api/user-metrics?organization=${currentData.organization}`);
                if (response.ok) {
                    const data = await response.json();
                    setUserMetrics(data);
                }
            } catch (error) {
                console.error("Error fetching user metrics:", error);
            } finally {
                setIsLoading(false);
            }
        };
        
        fetchUserMetrics();
    }, [currentData]);
    
    // Map seats to table data, including user metrics when available
    const tableData = (currentData?.seats ?? []).map((seat) => {
        const username = seat.assignee.login;
        const metrics = userMetrics[username] || {
            acceptanceRate: 'N/A',
            totalSuggestions: 'N/A',
            activeDays: 'N/A',
            timeSaved: 'N/A',
            mostUsedLanguages: 'N/A'
        };
        
        // Format metrics values for display
        const formatMetricValue = (value: any) => {
            if (value === null || value === undefined) return 'N/A';
            return value;
        };
        
        // Ensure mostUsedLanguages is a valid string
        if (typeof metrics.mostUsedLanguages !== 'string') {
            metrics.mostUsedLanguages = 'N/A';
        }

        // Sanitize values to ensure they can be safely used as keys
        const sanitizeForKey = (value: any): string => {
            if (value === null || value === undefined) return 'empty';
            // Convert to string and ensure it's not empty
            return String(value).trim() || 'empty';
        };
        
        return {
            // No rowNumber field - we'll let the table component handle this dynamically
            user: sanitizeForKey(username),
            organization: seat.organization?.login,
            team: seat.assigning_team?.name,
            createdAt: new Date(seat.created_at).toLocaleDateString(),
            updatedAt: new Date(seat.updated_at).toLocaleDateString(),
            lastActivityAt: seat.last_activity_at ? new Date(seat.last_activity_at).toLocaleDateString() : "-",
            lastActivityEditor: formatEditorName(seat.last_activity_editor),
            planType: sanitizeForKey(seat.plan_type),
            pendingCancellationDate: seat.pending_cancellation_date ? new Date(seat.pending_cancellation_date).toLocaleDateString() : "N/A",
            // Include user metrics with sanitization
            acceptanceRate: formatMetricValue(metrics.acceptanceRate),
            totalSuggestions: formatMetricValue(metrics.totalSuggestions),
            activeDays: formatMetricValue(metrics.activeDays),
            timeSaved: formatMetricValue(metrics.timeSaved),
            mostUsedLanguages: formatMetricValue(metrics.mostUsedLanguages)
        };
    });
    
    return (
        <Card className="col-span-4">
            <ChartHeader
                title="Assigned Seats with Usage Metrics"
                description={isLoading ? "Loading usage metrics..." : ""}
            />
            <CardContent>
                <DataTable
                    columns={columns.filter((col) => col.id !== "organization" || hasOrganization)}
                    data={tableData}
                    initialVisibleColumns={{
                        updatedAt: false,
                        pendingCancellationDate: false,
                        organization: false,
                        planType: false,
                        createdAt: false
                    }}
                    search={{
                        column: "user",
                        placeholder: "Filter seats...",
                    }}
                    filters={[
                        ...(hasOrganization ? [{ column: "organization", label: "Organizations" }] : []), 
                        ...(hasTeam ? [{ column: "team", label: "Team" }] : []),
                        { column: "planType", label: "Plan Type" },
                        { column: "mostUsedLanguages", label: "Languages" }
                    ]}
                    enableExport
                />
            </CardContent>
        </Card>
    );
};
