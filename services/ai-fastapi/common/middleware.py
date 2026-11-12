import time
import logging
from fastapi import Request
from starlette.middleware.base import BaseHTTPMiddleware

logger = logging.getLogger("ai_service.audit")


class AuditMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next):
        start_time = time.time()
        response = await call_next(request)
        process_time = (time.time() - start_time) * 1000
        logger.info(
            f"[AI AUDIT] {request.method} {request.url.path} - Status: {response.status_code} ({process_time:.2f}ms)"
        )
        response.headers["X-Process-Time"] = f"{process_time:.2f}ms"
        return response
