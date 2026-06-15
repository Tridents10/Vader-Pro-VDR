import sqlite3

DB_FILE = "vdr_telemetry.db"

def init_db():
    """Creates the SQLite database and enables WAL mode."""
    conn = sqlite3.connect(DB_FILE)
    cursor = conn.cursor()
    cursor.execute('PRAGMA journal_mode=WAL;')
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS telemetry (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            agent_uuid TEXT NOT NULL,
            received_at TEXT NOT NULL,
            inventory_data JSON NOT NULL
        )
    ''')
    conn.commit()
    conn.close()

def get_db_connection():
    """Helper function to open and return a database connection."""
    return sqlite3.connect(DB_FILE)