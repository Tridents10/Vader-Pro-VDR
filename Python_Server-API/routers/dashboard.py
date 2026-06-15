from fastapi import APIRouter, HTTPException
import json
from database import get_db_connection

# Create a router specifically for dashboard endpoints
router = APIRouter()

@router.get("/api/v1/agents")
async def get_all_agents():
    """Fetches all stored agent telemetry for the web dashboard."""
    try:
        # Ask database.py for a connection
        conn = get_db_connection()
        cursor = conn.cursor()
        cursor.execute("SELECT id, agent_uuid, received_at, inventory_data FROM telemetry ORDER BY received_at DESC")
        rows = cursor.fetchall()
        conn.close()
        
        agents_list = []
        for row in rows:
            agents_list.append({
                "id": row[0],
                "agent_uuid": row[1],
                "received_at": row[2],
                "inventory": json.loads(row[3]) 
            })
            
        return {"status": "success", "count": len(agents_list), "data": agents_list}
        
    except Exception as e:
        print(f"[!] API Error: {e}")
        raise HTTPException(status_code=500, detail="Failed to retrieve agent data")