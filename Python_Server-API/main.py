from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

# Import our custom modules
from database import init_db
from routers import ingestion, dashboard

app = FastAPI(title="VDR Central Management Server (Modular)")

# 1. Enable CORS for the web dashboard
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"], 
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# 2. Setup the Database on startup
init_db()

# 3. Plug in the Routers
app.include_router(ingestion.router)
app.include_router(dashboard.router)

# 4. Global Health Check
@app.get("/api/v1/health")
async def health_check():
    """A simple endpoint to test if the modular server is online."""
    return {"status": "online", "message": "Modular VDR Server is running"}