import math
import logging
from typing import List, Tuple
from app.modules.routing.schema import RoutingRequest, RouteSolution, Waypoint

logger = logging.getLogger(__name__)


def haversine_km(lat1: float, lon1: float, lat2: float, lon2: float) -> float:
    R = 6371.0  # Earth radius in kilometers
    dlat = math.radians(lat2 - lat1)
    dlon = math.radians(lon2 - lon1)
    a = (
        math.sin(dlat / 2) ** 2
        + math.cos(math.radians(lat1))
        * math.cos(math.radians(lat2))
        * math.sin(dlon / 2) ** 2
    )
    c = 2 * math.atan2(math.sqrt(a), math.sqrt(1 - a))
    return R * c


class RoutingOptimizationService:
    def solve_cvrp(self, req: RoutingRequest) -> RouteSolution:
        if not req.pickups:
            return RouteSolution(
                total_distance_km=0.0,
                total_weight_kg=0.0,
                vehicle_utilization_percent=0.0,
                route=[],
            )

        try:
            from ortools.constraint_solver import pywrapcp, routing_enums_pb2
            return self._solve_with_ortools(req)
        except ImportError:
            logger.warning("OR-Tools not available in current environment. Using heuristic greedy solver.")
            return self._solve_greedy(req)

    def _solve_with_ortools(self, req: RoutingRequest) -> RouteSolution:
        from ortools.constraint_solver import pywrapcp, routing_enums_pb2

        # Nodes: 0 is Depot, 1..N are Pickups
        all_nodes = [req.depot] + [p for p in req.pickups]
        n_nodes = len(all_nodes)

        # Build distance matrix (in meters for integer OR-Tools solver)
        distance_matrix = []
        for i in range(n_nodes):
            row = []
            for j in range(n_nodes):
                dist_km = haversine_km(
                    all_nodes[i].latitude,
                    all_nodes[i].longitude,
                    all_nodes[j].latitude,
                    all_nodes[j].longitude,
                )
                row.append(int(dist_km * 1000))  # Convert to integer meters
            distance_matrix.append(row)

        demands = [0] + [int(p.weight_kg) for p in req.pickups]
        vehicle_capacities = [int(req.vehicle_capacity_kg)]

        manager = pywrapcp.RoutingIndexManager(n_nodes, 1, 0)
        routing = pywrapcp.RoutingModel(manager)

        def distance_callback(from_index, to_index):
            from_node = manager.IndexToNode(from_index)
            to_node = manager.IndexToNode(to_index)
            return distance_matrix[from_node][to_node]

        transit_callback_index = routing.RegisterTransitCallback(distance_callback)
        routing.SetArcCostEvaluatorOfAllVehicles(transit_callback_index)

        def demand_callback(from_index):
            from_node = manager.IndexToNode(from_index)
            return demands[from_node]

        demand_callback_index = routing.RegisterUnaryTransitCallback(demand_callback)
        routing.AddDimensionWithVehicleCapacity(
            demand_callback_index,
            0,  # null capacity slack
            vehicle_capacities,  # vehicle maximum capacities
            True,  # start cumul to zero
            "Capacity",
        )

        search_parameters = pywrapcp.DefaultRoutingSearchParameters()
        search_parameters.first_solution_strategy = (
            routing_enums_pb2.FirstSolutionStrategy.PATH_CHEAPEST_ARC
        )

        solution = routing.SolveWithParameters(search_parameters)
        if not solution:
            return self._solve_greedy(req)

        index = routing.Start(0)
        ordered_waypoints: List[Waypoint] = []
        current_load = 0.0
        total_distance_m = 0

        while not routing.IsEnd(index):
            node_index = manager.IndexToNode(index)
            if node_index != 0:
                pickup = req.pickups[node_index - 1]
                current_load += pickup.weight_kg
                ordered_waypoints.append(
                    Waypoint(
                        pickup_id=pickup.id,
                        latitude=pickup.latitude,
                        longitude=pickup.longitude,
                        address=pickup.address,
                        load_after_pickup_kg=round(current_load, 2),
                    )
                )
            previous_index = index
            index = solution.Value(routing.NextVar(index))
            total_distance_m += routing.GetArcCostForVehicle(previous_index, index, 0)

        total_distance_km = round(total_distance_m / 1000.0, 2)
        total_weight_kg = round(current_load, 2)
        utilization = round((total_weight_kg / req.vehicle_capacity_kg) * 100, 1)

        return RouteSolution(
            total_distance_km=total_distance_km,
            total_weight_kg=total_weight_kg,
            vehicle_utilization_percent=min(utilization, 100.0),
            route=ordered_waypoints,
        )

    def _solve_greedy(self, req: RoutingRequest) -> RouteSolution:
        remaining = list(req.pickups)
        curr_lat = req.depot.latitude
        curr_lon = req.depot.longitude
        curr_weight = 0.0
        total_dist = 0.0
        waypoints = []

        while remaining:
            # Nearest unvisited node
            best_idx = 0
            best_dist = float("inf")
            for i, p in enumerate(remaining):
                d = haversine_km(curr_lat, curr_lon, p.latitude, p.longitude)
                if d < best_dist:
                    best_dist = d
                    best_idx = i

            chosen = remaining.pop(best_idx)
            total_dist += best_dist
            curr_lat = chosen.latitude
            curr_lon = chosen.longitude
            curr_weight += chosen.weight_kg

            waypoints.append(
                Waypoint(
                    pickup_id=chosen.id,
                    latitude=chosen.latitude,
                    longitude=chosen.longitude,
                    address=chosen.address,
                    load_after_pickup_kg=round(curr_weight, 2),
                )
            )

        # Return to depot
        total_dist += haversine_km(curr_lat, curr_lon, req.depot.latitude, req.depot.longitude)

        return RouteSolution(
            total_distance_km=round(total_dist, 2),
            total_weight_kg=round(curr_weight, 2),
            vehicle_utilization_percent=round((curr_weight / req.vehicle_capacity_kg) * 100, 1),
            route=waypoints,
        )


routing_service = RoutingOptimizationService()
