from pydantic import BaseModel, Field
from typing import List


class Coordinate(BaseModel):
    latitude: float
    longitude: float


class PickupNode(BaseModel):
    id: str
    latitude: float
    longitude: float
    weight_kg: float = 5.0
    address: str = ""


class RoutingRequest(BaseModel):
    depot: Coordinate = Field(..., description="Starting hub / waste bank depot location")
    vehicle_capacity_kg: float = Field(default=100.0, description="Max load capacity of courier vehicle")
    pickups: List[PickupNode] = Field(..., description="List of pickup points to visit")


class Waypoint(BaseModel):
    pickup_id: str
    latitude: float
    longitude: float
    address: str
    load_after_pickup_kg: float


class RouteSolution(BaseModel):
    total_distance_km: float
    total_weight_kg: float
    vehicle_utilization_percent: float
    route: List[Waypoint]
