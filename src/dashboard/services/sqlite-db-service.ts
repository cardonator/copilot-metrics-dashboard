import sqlite3 from 'sqlite3';
import { open, Database } from 'sqlite';
import { homedir } from 'os';
import path from 'path';
import fs from 'fs/promises';
import { ServerActionResponse } from '@/features/common/server-action-response';

let dbInstance: Database | null = null;

export const DEFAULT_DB_PATH = path.join(homedir(), '.copilot-metrics', 'copilot-metrics.db');

export const initializeDb = async (): Promise<Database> => {
  if (dbInstance) {
    return dbInstance;
  }

  const dbPath = process.env.SQLITE_DB_PATH || DEFAULT_DB_PATH;
  
  // Ensure directory exists
  const dbDir = path.dirname(dbPath);
  try {
    await fs.mkdir(dbDir, { recursive: true });
  } catch (error) {
    console.error('Error creating directory for SQLite database:', error);
  }

  dbInstance = await open({
    filename: dbPath,
    driver: sqlite3.Database
  });

  return dbInstance;
};

export const sqliteConfiguration = async (): Promise<boolean> => {
  // Check if storage type is set to sqlite
  if (process.env.STORAGE_TYPE === 'sqlite') {
    return true;
  }
  
  // Check if the default or specified file exists
  const dbPath = process.env.SQLITE_DB_PATH || DEFAULT_DB_PATH;
  try {
    return await fs.access(dbPath).then(() => true).catch(() => false);
  } catch {
    return false;
  }
};

export const closeDb = async (): Promise<void> => {
  if (dbInstance) {
    await dbInstance.close();
    dbInstance = null;
  }
};

export const queryDb = async <T>(
  query: string, 
  params: any[] = []
): Promise<ServerActionResponse<T[]>> => {
  try {
    const db = await initializeDb();
    
    // Handle empty params array properly to avoid syntax errors
    const results = params.length > 0 ? 
      await db.all(query, ...params) : 
      await db.all(query);
      
    return {
      status: "OK",
      response: results as T[]
    };
  } catch (error) {
    console.error('SQLite query error:', error);
    // Include the query in error message for better debugging
    return {
      status: "ERROR",
      errors: [{
        message: `Database query error: ${error instanceof Error ? error.message : 'Unknown error'} in query: ${query}`
      }]
    };
  }
};

// Add a new helper function to safely build WHERE clauses
export const buildWhereClause = (conditions: Record<string, any>): { 
  whereClause: string, 
  params: any[] 
} => {
  const whereConditions: string[] = [];
  const params: any[] = [];
  
  Object.entries(conditions).forEach(([key, value]) => {
    if (value !== undefined && value !== null) {
      whereConditions.push(`${key} = ?`);
      params.push(value);
    }
  });
  
  const whereClause = whereConditions.length > 0 ? 
    `WHERE ${whereConditions.join(' AND ')}` : '';
    
  return { whereClause, params };
};
