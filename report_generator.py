"""
Cross-Language BI Pipeline: Python Consumer
Queries the PostgreSQL metrics store and auto-generates executive Excel reports.
"""
import pandas as pd
import psycopg2
from datetime import datetime
import os

DB_CONFIG = {
    "dbname": "postgres",
    "user": "postgres",
    "password": "postgres",
    "host": "localhost",
    "port": "5432"
}

def generate_report():
    print("Initiating cross-language data handoff...")
    try:
        # Connect to the PostgreSQL database populated by the Go engine
        conn = psycopg2.connect(**DB_CONFIG)
        
        # Extract the latest aggregated metrics
        query = "SELECT * FROM ingestion_metrics ORDER BY timestamp DESC LIMIT 100"
        df = pd.read_sql_query(query, conn)
        
        if df.empty:
            print("--> [WARN] No data found in PostgreSQL. Run the Go processor first.")
            return

        # Generate Executive Excel Report
        os.makedirs("reports", exist_ok=True)
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        filename = f"reports/Executive_Metrics_Report_{timestamp}.xlsx"
        
        # Calculate derived metrics for the BI handoff
        df['failure_rate_pct'] = (df['failure_count'] / df['total_logs']) * 100
        df['failure_rate_pct'] = df['failure_rate_pct'].round(2)
        
        # Export
        df.to_excel(filename, index=False, sheet_name="Ingestion_Metrics")
        print(f"--> [SUCCESS] Enterprise BI Report generated: {filename}")
        
    except psycopg2.OperationalError:
        print("--> [WARN] PostgreSQL is not running locally. Run a local Postgres instance to generate reports.")
    except Exception as e:
        print(f"--> [ERROR] Pipeline failure: {e}")
    finally:
        if 'conn' in locals() and conn is not None:
            conn.close()

if __name__ == "__main__":
    generate_report()