from fastapi import APIRouter, Request, Header, HTTPException
from Crypto.Cipher import AES
import json
import datetime
from database import get_db_connection

# Create a router specifically for the ingestion endpoints
router = APIRouter()

SECRET_KEY = b"12345678901234567890123456789012"
API_TOKEN = "your-secure-pre-shared-token"

def decrypt_payload(encrypted_payload: bytes) -> dict:
    """Unlocks the AES-256-GCM payload."""
    try:
        # Standard GCM layout: 12-byte Nonce + Ciphertext + 16-byte Tag
        nonce = encrypted_payload[:12]
        tag = encrypted_payload[-16:]
        ciphertext = encrypted_payload[12:-16]
        
        # Initialize cipher and decrypt/verify simultaneously
        cipher = AES.new(SECRET_KEY, AES.MODE_GCM, nonce=nonce)
        decrypted_bytes = cipher.decrypt_and_verify(ciphertext, tag)
        
        return json.loads(decrypted_bytes.decode('utf-8'))
    except Exception as e:
        print(f"[!] Decryption Error: {e}")
        raise ValueError("Decryption failed!")

@router.post("/api/v1/submit")
async def receive_telemetry(
    request: Request,
    x_agent_uuid: str = Header(None),
    authorization: str = Header(None)
):
    """Endpoint for edge agents to submit encrypted telemetry."""
    expected_auth = f"Bearer {API_TOKEN}"
    if authorization != expected_auth:
        raise HTTPException(status_code=401, detail="Unauthorized")
    
    if not x_agent_uuid:
        raise HTTPException(status_code=400, detail="Missing X-Agent-UUID")

    encrypted_data = await request.body()
    if not encrypted_data:
        raise HTTPException(status_code=400, detail="Empty payload")

    try:
        inventory_data = decrypt_payload(encrypted_data)
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))

    try:
        # Ask database.py for a connection
        conn = get_db_connection()
        cursor = conn.cursor()
        cursor.execute(
            "INSERT INTO telemetry (agent_uuid, received_at, inventory_data) VALUES (?, ?, ?)",
            (x_agent_uuid, datetime.datetime.utcnow().isoformat(), json.dumps(inventory_data))
        )
        conn.commit()
        conn.close()
    except Exception as e:
        print(f"[!] DB Error: {e}")
        raise HTTPException(status_code=500, detail="Database Error")

    print(f"[+] Decrypted telemetry from Agent: {x_agent_uuid}")
    return {"status": "success", "message": "Telemetry securely stored"}