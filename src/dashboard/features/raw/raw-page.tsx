"use client";

import { Card, CardContent } from "@/components/ui/card";
import { PageHeader, PageTitle } from "@/features/page-header/page-header";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { Check, Copy } from "lucide-react";
import { useState } from "react";
import { toast } from "@/components/ui/use-toast";

interface RawDataPageProps {
  metricsData: any;
  seatsData: any | null;
}

export function RawDataPage({ metricsData, seatsData }: RawDataPageProps) {
  const [activeTab, setActiveTab] = useState<string>("metrics");
  const [copied, setCopied] = useState(false);
  
  const handleTabChange = (value: string) => {
    setActiveTab(value);
    setCopied(false); // Reset copied state when changing tabs
  };

  // Get the raw API responses
  const metricsRawResponse = metricsData.rawApiResponse || JSON.stringify(metricsData, null, 2);
  const seatsRawResponse = seatsData?.rawApiResponse || (seatsData ? JSON.stringify(seatsData, null, 2) : null);

  // Get the current content based on active tab
  const currentContent = activeTab === "metrics" ? metricsRawResponse : seatsRawResponse;

  const copyToClipboard = () => {
    if (!currentContent) return;
    
    // Create a text area element to copy from
    const textArea = document.createElement('textarea');
    textArea.value = currentContent;
    
    // Make it invisible and add to the DOM
    textArea.style.position = 'fixed';
    textArea.style.opacity = '0';
    document.body.appendChild(textArea);
    
    // Select and copy the text
    textArea.select();
    
    try {
      const successful = document.execCommand('copy');
      if (successful) {
        setCopied(true);
        toast({
          title: "Copied!",
          description: "Raw data copied to clipboard",
          duration: 2000,
        });
        
        // Reset the copied status after 2 seconds
        setTimeout(() => {
          setCopied(false);
        }, 2000);
      }
    } catch (err) {
      console.error('Failed to copy text: ', err);
      toast({
        title: "Failed to copy",
        description: "There was an error copying to clipboard",
        variant: "destructive",
        duration: 2000,
      });
    }
    
    // Remove the text area
    document.body.removeChild(textArea);
  };

  return (
    <div className="flex flex-col h-[calc(100vh-64px)]">
      <PageHeader>
        <PageTitle>Raw API Responses</PageTitle>
      </PageHeader>

      <div className="mx-auto w-full max-w-6xl container flex-1 overflow-hidden p-4">
        <Card className="h-full flex flex-col">
          <Tabs defaultValue="metrics" onValueChange={handleTabChange} className="flex flex-col h-full">
            <div className="px-6 pt-6 pb-2 flex-none">
              <TabsList className="grid w-full grid-cols-2 max-w-md">
                <TabsTrigger value="metrics">Metrics API Response</TabsTrigger>
                <TabsTrigger value="seats" disabled={!seatsData}>Seats API Response</TabsTrigger>
              </TabsList>
            </div>
            
            <CardContent className="flex-1 overflow-hidden p-4">
              <div className="relative h-full w-full rounded-md border bg-muted/5">
                <div className="absolute right-2 top-2 z-10">
                  <Button
                    size="sm"
                    variant="secondary"
                    className="h-8 w-8 p-0 rounded-full"
                    onClick={copyToClipboard}
                    title="Copy to clipboard"
                  >
                    {copied ? (
                      <Check className="h-4 w-4 text-green-500" />
                    ) : (
                      <Copy className="h-4 w-4" />
                    )}
                    <span className="sr-only">Copy to clipboard</span>
                  </Button>
                </div>
                <ScrollArea className="h-full w-full p-4">
                  <pre className="whitespace-pre-wrap font-mono text-xs pt-8">
                    {currentContent}
                  </pre>
                </ScrollArea>
              </div>
            </CardContent>
          </Tabs>
        </Card>
      </div>
    </div>
  );
}
