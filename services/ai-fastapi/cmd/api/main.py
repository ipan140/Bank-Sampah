import logging
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from config.settings import settings
from common.errors import global_exception_handler
from common.middleware import AuditMiddleware

from app.modules.classify.handler import router as classify_router
from app.modules.estimate.handler import router as estimate_router
from app.modules.routing.handler import router as routing_router

# Setup Logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
)
logger = logging.getLogger("ai_service")

app = FastAPI(
    title="Bank Sampah AI & Optimization Microservice",
    description="Waste visual classification (ONNX), price & carbon estimation, and Google OR-Tools VRP routing",
    version="1.0.0",
)

# Global CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
app.add_middleware(AuditMiddleware)
app.add_exception_handler(Exception, global_exception_handler)


@app.get("/health", tags=["Health"])
async def health():
    return {
        "status": "healthy",
        "service": "ai-fastapi",
        "environment": settings.APP_ENV,
    }


# Include Routers
app.include_router(classify_router, prefix="/api/v1")
app.include_router(estimate_router, prefix="/api/v1")
app.include_router(routing_router, prefix="/api/v1")

if __name__ == "__main__":
    import uvicorn
    uvicorn.run("cmd.api.main:app", host="0.0.0.0", port=settings.APP_PORT, reload=True)
