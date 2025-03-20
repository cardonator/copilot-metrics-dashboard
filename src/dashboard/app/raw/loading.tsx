import { Skeleton } from "@/components/ui/skeleton";
import { PageHeader, PageTitle } from "@/features/page-header/page-header";

export default function Loading() {
  return (
    <div className="flex flex-col h-[calc(100vh-64px)]">
      <PageHeader>
        <PageTitle>Raw Data</PageTitle>
      </PageHeader>

      <div className="mx-auto w-full max-w-6xl container flex-1 overflow-hidden p-4">
        <div className="h-full flex flex-col border rounded-lg">
          <div className="p-6 flex-none">
            <Skeleton className="h-10 w-80" />
          </div>
          <div className="flex-1 p-4">
            <Skeleton className="h-full w-full" />
          </div>
        </div>
      </div>
    </div>
  );
}
