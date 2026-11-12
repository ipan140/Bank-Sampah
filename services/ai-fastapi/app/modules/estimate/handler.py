from fastapi import APIRouter
from app.modules.estimate.schema import EstimateRequest, EstimateResponse
from app.modules.estimate.service import estimate_service

router = APIRouter(prefix="/estimate", tags=["Price & Weight Estimation"])


@router.post("", response_model=EstimateResponse)
async def estimate_waste(req: EstimateRequest):
    return estimate_service.calculate_estimate(req)
