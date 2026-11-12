from pydantic import BaseModel, Field
from typing import List, Optional


class WasteDetection(BaseModel):
    label: str = Field(..., description="Detected category: PLASTIC, PAPER, METAL, GLASS, ORGANIC")
    confidence: float = Field(..., description="Detection confidence score between 0.0 and 1.0")
    box: Optional[List[float]] = Field(None, description="[ymin, xmin, ymax, xmax] normalized coordinates")


class ClassificationResponse(BaseModel):
    success: bool
    primary_category: str
    confidence: float
    detections: List[WasteDetection]
    suggested_action: str
