from pydantic_settings import BaseSettings
from typing import Optional


class Settings(BaseSettings):
    APP_ENV: str = "development"
    APP_PORT: int = 8000
    CORE_SERVICE_URL: str = "http://localhost:8080"
    RABBITMQ_URL: Optional[str] = "amqp://guest:guest@localhost:5672/"
    MODEL_PATH: str = "models/waste_yolov8.onnx"

    class Config:
        env_file = ".env"
        extra = "ignore"


settings = Settings()
