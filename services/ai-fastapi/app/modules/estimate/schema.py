from pydantic import BaseModel, Field
from typing import Optional


class EstimateRequest(BaseModel):
    category: str = Field(..., description="Category: PLASTIC, PAPER, METAL, GLASS, ORGANIC")
    volume_liters: Optional[float] = Field(None, description="Volume in Liters, e.g. trash bag size (30L, 50L)")
    count_items: Optional[int] = Field(None, description="Approximate count of items, e.g. 20 bottles")
    is_compacted: bool = Field(default=False, description="Whether waste is crushed/flattened")


class EstimateResponse(BaseModel):
    category: str
    estimated_weight_kg: float
    unit_price_per_kg: float
    estimated_total_value_idr: float
    carbon_offset_kg: float
