from fastapi import APIRouter
from app.modules.routing.schema import RoutingRequest, RouteSolution
from app.modules.routing.service import routing_service

router = APIRouter(prefix="/routing", tags=["Vehicle Routing Problem"])


@router.post("/optimize", response_model=RouteSolution)
async def optimize_route(req: RoutingRequest):
    return routing_service.solve_cvrp(req)
