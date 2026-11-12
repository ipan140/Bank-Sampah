from fastapi import APIRouter, File, UploadFile, HTTPException
from app.modules.classify.schema import ClassificationResponse
from app.modules.classify.service import classifier_service

router = APIRouter(prefix="/classify", tags=["Waste Classification"])


@router.post("", response_model=ClassificationResponse)
async def classify_waste(file: UploadFile = File(...)):
    if not file.content_type.startswith("image/"):
        raise HTTPException(status_code=400, detail="Uploaded file must be an image (JPEG, PNG, WEBP)")

    content = await file.read()
    if len(content) == 0:
        raise HTTPException(status_code=400, detail="Empty image uploaded")

    return await classifier_service.classify_image(content)
